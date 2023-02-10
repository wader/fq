// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//nolint:all
package tlsdecrypt

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
)

const (
	tlsRandomLength      = 32 // Length of a random nonce in TLS 1.1.
	masterSecretLength   = 48 // Length of a master secret in TLS 1.1.
	finishedVerifyLength = 12 // Length of verify_data in a Finished message.
)

var masterSecretLabel = []byte("master secret")
var extendedMasterSecretLabel = []byte("extended master secret")
var keyExpansionLabel = []byte("key expansion")
var clientFinishedLabel = []byte("client finished")
var serverFinishedLabel = []byte("server finished")
var clientFinalKeyLabel = []byte("client write key")
var serverFinalKeyLabel = []byte("server write key")
var finalIVLabel = []byte("IV block")

// Split a premaster secret in two as specified in RFC 4346, Section 5.
func splitPreMasterSecret(secret []byte) (s1, s2 []byte) {
	s1 = secret[0 : (len(secret)+1)/2]
	s2 = secret[len(secret)/2:]
	return
}

// pHash implements the P_hash function, as defined in RFC 4346, Section 5.
func pHash(result, secret, seed []byte, hash func() hash.Hash) {
	h := hmac.New(hash, secret)
	h.Write(seed)
	a := h.Sum(nil)

	j := 0
	for j < len(result) {
		h.Reset()
		h.Write(a)
		h.Write(seed)
		b := h.Sum(nil)
		copy(result[j:], b)
		j += len(b)

		h.Reset()
		h.Write(a)
		a = h.Sum(nil)
	}
}

// prf10 implements the TLS 1.0 pseudo-random function, as defined in RFC 2246, section 5.
func prf10(result, secret, label, seed []byte) {
	hashSHA1 := sha1.New
	hashMD5 := md5.New

	labelAndSeed := make([]byte, len(label)+len(seed))
	copy(labelAndSeed, label)
	copy(labelAndSeed[len(label):], seed)

	s1, s2 := splitPreMasterSecret(secret)
	pHash(result, s1, labelAndSeed, hashMD5)
	result2 := make([]byte, len(result))
	pHash(result2, s2, labelAndSeed, hashSHA1)

	for i, b := range result2 {
		result[i] ^= b
	}
}

// prf12 implements the TLS 1.2 pseudo-random function, as defined in RFC 5246, section 5.
func prf12(hashFunc func() hash.Hash) func(result, secret, label, seed []byte) {
	return func(result, secret, label, seed []byte) {
		labelAndSeed := make([]byte, len(label)+len(seed))
		copy(labelAndSeed, label)
		copy(labelAndSeed[len(label):], seed)

		pHash(result, secret, labelAndSeed, hashFunc)
	}
}

// prf30 implements the SSL 3.0 pseudo-random function, as defined in
// www.mozilla.org/projects/security/pki/nss/ssl/draft302.txt section 6.
func prf30(result, secret, label, seed []byte) {
	hashSHA1 := sha1.New()
	hashMD5 := md5.New()

	done := 0
	i := 0
	// RFC5246 section 6.3 says that the largest PRF output needed is 128
	// bytes. Since no more ciphersuites will be added to SSLv3, this will
	// remain true. Each iteration gives us 16 bytes so 10 iterations will
	// be sufficient.
	var b [11]byte
	for done < len(result) {
		for j := 0; j <= i; j++ {
			b[j] = 'A' + byte(i)
		}

		hashSHA1.Reset()
		hashSHA1.Write(b[:i+1])
		hashSHA1.Write(secret)
		hashSHA1.Write(seed)
		digest := hashSHA1.Sum(nil)

		hashMD5.Reset()
		hashMD5.Write(secret)
		hashMD5.Write(digest)

		done += copy(result[done:], hashMD5.Sum(nil))
		i++
	}
}

func exportPRF30(result, secret, label, seed []byte) {
	hash := md5.New()
	hash.Write(secret)
	hash.Write(seed)
	copy(result, hash.Sum(nil))
}
func prfAndHashForVersion(version uint16, suite *cipherSuite) (func(result, secret, label, seed []byte), crypto.Hash) {
	switch version {
	case VersionTLS10, VersionTLS11:
		return prf10, crypto.Hash(0)
	case VersionTLS12:
		if suite.flags&suiteSHA384 != 0 {
			return prf12(sha512.New384), crypto.SHA384
		}
		return prf12(sha256.New), crypto.SHA256
	default:
		panic("unknown version")
	}
}

func prfForVersion(version uint16, suite *cipherSuite) func(result, secret, label, seed []byte) {
	prf, _ := prfAndHashForVersion(version, suite)
	return prf
}

// keysFromMasterSecret generates the connection keys from the master
// secret, given the lengths of the MAC key, cipher key and IV, as defined in
// RFC 2246, section 6.3.
func keysFromMasterSecret(version uint16, suite *cipherSuite, masterSecret, clientRandom, serverRandom []byte, macLen, keyLen, ivLen int) (clientMAC, serverMAC, clientKey, serverKey, clientIV, serverIV []byte) {
	if suite.flags&suiteExport > 0 {
		return exportKeysFromMasterSecret(version, suite, masterSecret, clientRandom, serverRandom, macLen, keyLen, ivLen)
	}
	var seed [tlsRandomLength * 2]byte
	copy(seed[0:len(clientRandom)], serverRandom)
	copy(seed[len(serverRandom):], clientRandom)

	n := 2*macLen + 2*keyLen + 2*ivLen
	keyMaterial := make([]byte, n)
	prfForVersion(version, suite)(keyMaterial, masterSecret, keyExpansionLabel, seed[0:])
	clientMAC = keyMaterial[:macLen]
	keyMaterial = keyMaterial[macLen:]
	serverMAC = keyMaterial[:macLen]
	keyMaterial = keyMaterial[macLen:]
	clientKey = keyMaterial[:keyLen]
	keyMaterial = keyMaterial[keyLen:]
	serverKey = keyMaterial[:keyLen]
	keyMaterial = keyMaterial[keyLen:]
	clientIV = keyMaterial[:ivLen]
	keyMaterial = keyMaterial[ivLen:]
	serverIV = keyMaterial[:ivLen]
	return
}

// The crypto wars must have been the worst
func exportKeysFromMasterSecret30(version uint16, suite *cipherSuite, masterSecret, clientRandom, serverRandom []byte, macLen, keyLen, ivLen int) (clientMAC, serverMAC, clientKey, serverKey, clientIV, serverIV []byte) {
	var seed [tlsRandomLength * 2]byte
	copy(seed[0:len(clientRandom)], serverRandom)
	copy(seed[len(serverRandom):], clientRandom)
	n := 2*macLen + 2*keyLen
	keyMaterial := make([]byte, n)
	prf30(keyMaterial, masterSecret, keyExpansionLabel, seed[0:])
	clientMAC = keyMaterial[:macLen]
	keyMaterial = keyMaterial[macLen:]
	serverMAC = keyMaterial[:macLen]
	keyMaterial = keyMaterial[macLen:]
	clientKey = keyMaterial[:keyLen]
	keyMaterial = keyMaterial[keyLen:]
	serverKey = keyMaterial[:keyLen]
	var exportSeed [tlsRandomLength * 2]byte
	copy(exportSeed[0:len(serverRandom)], clientRandom)
	copy(exportSeed[len(clientRandom):], serverRandom)
	expandedKeyLen := suite.expandedKeyLen
	finalKeyBlock := make([]byte, 2*expandedKeyLen)
	exportPRF30(finalKeyBlock[:expandedKeyLen], clientKey, clientFinalKeyLabel, exportSeed[0:])
	clientKey = finalKeyBlock[:expandedKeyLen]
	finalKeyBlock = finalKeyBlock[expandedKeyLen:]
	exportPRF30(finalKeyBlock[:expandedKeyLen], serverKey, serverFinalKeyLabel, seed[0:])
	serverKey = finalKeyBlock[:expandedKeyLen]
	ivBlock := make([]byte, 2*ivLen)
	clientIV = ivBlock[:ivLen]
	exportPRF30(clientIV, []byte{}, finalIVLabel, exportSeed[0:])
	ivBlock = ivBlock[ivLen:]
	serverIV = ivBlock[:ivLen]
	exportPRF30(serverIV, []byte{}, finalIVLabel, seed[0:])
	return
}

// If a cryptographer kills me in the night, let it be known I was sorry
func exportKeysFromMasterSecretTLS(version uint16, suite *cipherSuite, masterSecret, clientRandom, serverRandom []byte, macLen, keyLen, ivLen int) (clientMAC, serverMAC, clientKey, serverKey, clientIV, serverIV []byte) {
	var seed [tlsRandomLength * 2]byte
	copy(seed[0:len(clientRandom)], serverRandom)
	copy(seed[len(serverRandom):], clientRandom)
	n := 2*macLen + 2*keyLen
	keyMaterial := make([]byte, n)
	prf := prfForVersion(version, suite)
	prf(keyMaterial, masterSecret, keyExpansionLabel, seed[0:])
	clientMAC = keyMaterial[:macLen]
	keyMaterial = keyMaterial[macLen:]
	serverMAC = keyMaterial[:macLen]
	keyMaterial = keyMaterial[macLen:]
	clientKey = keyMaterial[:keyLen]
	keyMaterial = keyMaterial[keyLen:]
	serverKey = keyMaterial[:keyLen]
	expandedKeyLen := suite.expandedKeyLen
	finalKeyBlock := make([]byte, 2*expandedKeyLen)
	var exportSeed [tlsRandomLength * 2]byte
	copy(exportSeed[0:len(serverRandom)], clientRandom)
	copy(exportSeed[len(clientRandom):], serverRandom)
	prf(finalKeyBlock[:expandedKeyLen], clientKey, clientFinalKeyLabel, exportSeed[0:])
	clientKey = finalKeyBlock[:expandedKeyLen]
	finalKeyBlock = finalKeyBlock[expandedKeyLen:]
	prf(finalKeyBlock[:expandedKeyLen], serverKey, serverFinalKeyLabel, exportSeed[0:])
	serverKey = finalKeyBlock[:expandedKeyLen]
	ivBlock := make([]byte, 2*ivLen)
	prf(ivBlock, []byte{}, finalIVLabel, exportSeed[0:])
	clientIV = ivBlock[:ivLen]
	ivBlock = ivBlock[ivLen:]
	serverIV = ivBlock[:ivLen]
	return
}

func exportKeysFromMasterSecret(version uint16, suite *cipherSuite, masterSecret, clientRandom, serverRandom []byte, macLen, keyLen, ivLen int) (clientMAC, serverMAC, clientKey, serverKey, clientIV, serverIV []byte) {
	switch version {
	case VersionSSL30:
		return exportKeysFromMasterSecret30(version, suite, masterSecret, clientRandom, serverRandom, macLen, keyLen, ivLen)
	case VersionTLS10, VersionTLS11, VersionTLS12:
		return exportKeysFromMasterSecretTLS(version, suite, masterSecret, clientRandom, serverRandom, macLen, keyLen, ivLen)
	default:
		panic("unknown version")
	}
}

// noExportedKeyingMaterial is used as a value of
// ConnectionState.ekm when renegotiation is enabled and thus
// we wish to fail all key-material export requests.
func noExportedKeyingMaterial(label string, context []byte, length int) ([]byte, error) {
	return nil, errors.New("crypto/tls: ExportKeyingMaterial is unavailable when renegotiation is enabled")
}

// ekmFromMasterSecret generates exported keying material as defined in RFC 5705.
func ekmFromMasterSecret(version uint16, suite *cipherSuite, masterSecret, clientRandom, serverRandom []byte) func(string, []byte, int) ([]byte, error) {
	return func(label string, context []byte, length int) ([]byte, error) {
		switch label {
		case "client finished", "server finished", "master secret", "key expansion":
			// These values are reserved and may not be used.
			return nil, fmt.Errorf("crypto/tls: reserved ExportKeyingMaterial label: %s", label)
		}

		seedLen := len(serverRandom) + len(clientRandom)
		if context != nil {
			seedLen += 2 + len(context)
		}
		seed := make([]byte, 0, seedLen)

		seed = append(seed, clientRandom...)
		seed = append(seed, serverRandom...)

		if context != nil {
			if len(context) >= 1<<16 {
				return nil, fmt.Errorf("crypto/tls: ExportKeyingMaterial context too long")
			}
			seed = append(seed, byte(len(context)>>8), byte(len(context)))
			seed = append(seed, context...)
		}

		keyMaterial := make([]byte, length)
		prfForVersion(version, suite)(keyMaterial, masterSecret, []byte(label), seed)
		return keyMaterial, nil
	}
}
