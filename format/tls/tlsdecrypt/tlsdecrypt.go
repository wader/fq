package tlsdecrypt

// TODO: SSLv3 MAC
// TODO: TLS 1.3

import (
	"fmt"
	"hash"
)

type Decryptor struct {
	IsClient     bool
	Version      int
	CipherSuite  int
	MasterSecret []byte
	ClientRandom []byte
	ServerRandom []byte

	halfConn *halfConn
}

type keys struct {
	clientCipher any
	serverCipher any
	clientHash   hash.Hash
	serverHash   hash.Hash
}

// TODO: only need one direction at a time
func establishKeys(
	vers uint16,
	suite *cipherSuite,
	masterSecret []byte,
	clientRandom []byte,
	serverRandom []byte,
) keys {
	clientMAC, serverMAC, clientKey, serverKey, clientIV, serverIV :=
		keysFromMasterSecret(vers, suite, masterSecret, clientRandom, serverRandom, suite.macLen, suite.keyLen, suite.ivLen)

	var clientCipher, serverCipher any
	var clientHash, serverHash hash.Hash

	if suite.aead == nil {
		clientCipher = suite.cipher(clientKey, clientIV, true /* for reading */)
		clientHash = suite.mac(clientMAC)
		serverCipher = suite.cipher(serverKey, serverIV, true /* for reading */)
		serverHash = suite.mac(serverMAC)
	} else {
		clientCipher = suite.aead(clientKey, clientIV)
		serverCipher = suite.aead(serverKey, serverIV)
	}

	return keys{
		clientCipher: clientCipher,
		serverCipher: serverCipher,
		clientHash:   clientHash,
		serverHash:   serverHash,
	}
}

func (d *Decryptor) Decrypt(record []byte) ([]byte, error) {
	if d.halfConn == nil {
		cipherSuite := cipherSuiteByID(uint16(d.CipherSuite))
		if cipherSuite == nil {
			return nil, fmt.Errorf("unsupported cipher suit %x", d.CipherSuite)
		}

		keys := establishKeys(
			uint16(d.Version),
			cipherSuite,
			d.MasterSecret,
			d.ClientRandom,
			d.ServerRandom,
		)

		var cipher any
		var mac hash.Hash
		if d.IsClient {
			cipher = keys.clientCipher
			mac = keys.clientHash
		} else {
			cipher = keys.serverCipher
			mac = keys.serverHash
		}

		d.halfConn = &halfConn{
			version: uint16(d.Version),
			cipher:  cipher,
			mac:     mac,
			seq:     [8]byte{}, // zero
		}
	}

	plain, _, err := d.halfConn.decrypt(record)
	return plain, err
}
