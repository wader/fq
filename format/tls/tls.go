package tls

// NOTES:
// - after first cipher change records are assume to be encrypted, currently even for
//   null encryption decode will happen in post
//
// TODO: key exchange alg, decode key exchange parameters
// TODO: renegotiation, client/server hello again etc, uses current cipher state, keep track of key change
// TODO: tls 1.3, ssl? combine or own format?
// TODO: ALPN
// TODO: pcapng keylog
// TODO: add fields for seq, calculated things? prf result and decode key/iv?
// TODO: warnings to stderr decode api support?
//
// The wireshark TLS/SSL dissector code is a great reference to look at while
// reading the TLS specification:
// https://github.com/boundary/wireshark/blob/master/epan/dissectors/packet-ssl.c
// https://github.com/boundary/wireshark/blob/master/epan/dissectors/packet-ssl-utils.c

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/tls/ciphersuites"
	"github.com/wader/fq/format/tls/keylog"
	"github.com/wader/fq/format/tls/rezlib"
	"github.com/wader/fq/format/tls/tlsdecrypt"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/ranges"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed tls.md
var tlsFS embed.FS

var asn1BerGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.TLS,
		&decode.Format{
			Description:  "Transport layer security",
			Groups:       []*decode.Group{format.TCP_Stream},
			DecodeFn:     decodeTLS,
			DefaultInArg: format.TLS_In{},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.ASN1_BER}, Out: &asn1BerGroup},
			},
		})
	interp.RegisterFS(tlsFS)
}

const (
	versionSSL     = 0x0300
	versionTLS_1_0 = 0x0301
	versionTLS_1_1 = 0x0302
	versionTLS_1_2 = 0x0303
	versionTLS_1_3 = 0x0304
)

var versionNames = scalar.UintMapSymStr{
	versionSSL:     "ssl",
	versionTLS_1_0: "tls1.0",
	versionTLS_1_1: "tls1.1",
	versionTLS_1_2: "tls1.2",
	versionTLS_1_3: "tls1.3",
}

var versionValid = []uint64{
	versionSSL,
	versionTLS_1_0,
	versionTLS_1_1,
	versionTLS_1_2,
	versionTLS_1_3,
}

const (
	recordTypeChangeCipherSpec = 20
	recordTypeAlert            = 21
	recordTypeHandshake        = 22
	recordTypeApplicationData  = 23
)

var recordTypeNames = scalar.UintMapSymStr{
	recordTypeChangeCipherSpec: "change_cipher_spec",
	recordTypeAlert:            "alert",
	recordTypeHandshake:        "handshake",
	recordTypeApplicationData:  "application_data",
}

var recordTypeValid = []uint64{
	recordTypeChangeCipherSpec,
	recordTypeAlert,
	recordTypeHandshake,
	recordTypeApplicationData,
}

const (
	handshakeMsgTypeHelloRequest       = 0
	handshakeMsgTypeClientHello        = 1
	handshakeMsgTypeServerHello        = 2
	handshakeMsgTypeNewSessionTicket   = 4
	handshakeMsgTypeCertificate        = 11
	handshakeMsgTypeServerKeyExchange  = 12
	handshakeMsgTypeCertificateRequest = 13
	handshakeMsgTypeServerHelloDone    = 14
	handshakeMsgTypeCertificateVerify  = 15
	handshakeMsgTypeClientKeyExchange  = 16
	handshakeMsgTypeFinished           = 20
)

var handshakeMsgTypeNames = scalar.UintMapSymStr{
	handshakeMsgTypeHelloRequest:       "hello_request",
	handshakeMsgTypeClientHello:        "client_hello",
	handshakeMsgTypeServerHello:        "server_hello",
	handshakeMsgTypeNewSessionTicket:   "new_session_ticket",
	handshakeMsgTypeCertificate:        "certificate",
	handshakeMsgTypeServerKeyExchange:  "server_key_exchange",
	handshakeMsgTypeCertificateRequest: "certificate_request",
	handshakeMsgTypeServerHelloDone:    "server_hello_done",
	handshakeMsgTypeCertificateVerify:  "certificate_verify",
	handshakeMsgTypeClientKeyExchange:  "client_key_exchange",
	handshakeMsgTypeFinished:           "finished",
}

const (
	compressionMethodNull    = 0
	compressionMethodDeflate = 1
)

var compressionMethodNames = scalar.UintMapSymStr{
	compressionMethodNull:    "null",
	compressionMethodDeflate: "deflate",
}

const (
	changeCipherSpecTypeChangeCipherSpec = 0
)

var changeCipherSpecTypeNames = scalar.UintMapSymStr{
	changeCipherSpecTypeChangeCipherSpec: "change_cipher_spec",
}

var cipherNames = scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
	if suit, ok := ciphersuites.Suits[int(s.Actual)]; ok {
		s.Sym = suit.Name
	}
	return s, nil
})

const (
	signatureAlgorithmAnonymous         = 0
	signatureAlgorithmRSA               = 1
	signatureAlgorithmDSA               = 2
	signatureAlgorithmECDSA             = 3
	signatureAlgorithmEd25519           = 7
	signatureAlgorithmEd448             = 8
	signatureAlgorithmGOSTR34102012_256 = 64
	signatureAlgorithmGOSTR34102012_512 = 65
)

var signatureAlgorithmNames = scalar.UintMapSymStr{
	signatureAlgorithmAnonymous:         "anonymous",
	signatureAlgorithmRSA:               "rsa",
	signatureAlgorithmDSA:               "dsa",
	signatureAlgorithmECDSA:             "ecdsa",
	signatureAlgorithmEd25519:           "ed25519",
	signatureAlgorithmEd448:             "ed448",
	signatureAlgorithmGOSTR34102012_256: "gostr34102012_256",
	signatureAlgorithmGOSTR34102012_512: "gostr34102012_512",
}

const (
	hashAlgorithmnone      = 0
	hashAlgorithmMD5       = 1
	hashAlgorithmSHA1      = 2
	hashAlgorithmSHA224    = 3
	hashAlgorithmSHA256    = 4
	hashAlgorithmSHA384    = 5
	hashAlgorithmSHA512    = 6
	hashAlgorithmIntrinsic = 8
)

var hashAlgorithmNames = scalar.UintMapSymStr{
	hashAlgorithmnone:      "none",
	hashAlgorithmMD5:       "md5",
	hashAlgorithmSHA1:      "sha1",
	hashAlgorithmSHA224:    "sha224",
	hashAlgorithmSHA256:    "sha256",
	hashAlgorithmSHA384:    "sha384",
	hashAlgorithmSHA512:    "sha512",
	hashAlgorithmIntrinsic: "intrinsic",
}

type encryptedRecord struct {
	recordType uint64
	d          *decode.D
	r          ranges.Range
}

type keyExchange struct {
	msgType uint64
	d       *decode.D
	r       ranges.Range
	dataV   *decode.Value
}

type tlsCtx struct {
	rootD *decode.D

	version uint64
	random  [32]byte

	// server only
	server struct {
		currentCipherSuit uint64
		nextCipherSuit    uint64
		compressionMethod uint64
	}

	// cipher has been decided
	isEncrypted      bool
	encryptedRecords []encryptedRecord
	keyExchange      *keyExchange

	serverCtx *tlsCtx
	clientCtx *tlsCtx
}

func decodeTLSExtension(d *decode.D) {
	typ := d.FieldU16("type", extensionNames)
	length := d.FieldU16("length")
	// server sometimes use empty extension to indicate things, ex: accept SNI
	if length == 0 {
		return
	}
	d.FramedFn(int64(length)*8, func(d *decode.D) {
		switch typ {
		case extensionServerName:
			serverNameLength := d.FieldU16("server_names_length")
			d.FieldArray("server_names", func(d *decode.D) {
				d.FramedFn(int64(serverNameLength)*8, func(d *decode.D) {
					for !d.End() {
						d.FieldStruct("server_name", func(d *decode.D) {
							d.FieldU8("type") // TODO:
							length := d.FieldU16("length")
							d.FieldUTF8("name", int(length))
						})
					}
				})
			})
		case extensionApplicationLayerProtocolNegotiation:
			protocolsLength := d.FieldU16("protocols_length")
			d.FieldArray("protocols", func(d *decode.D) {
				d.FramedFn(int64(protocolsLength)*8, func(d *decode.D) {
					for !d.End() {
						d.FieldStruct("protocol", func(d *decode.D) {
							length := d.FieldU8("length")
							d.FieldUTF8("name", int(length))
						})
					}
				})
			})
		case extensionEcPointFormats:
			protocolsLength := d.FieldU8("ex_points_formats_length")
			d.FieldArray("ex_points_formats", func(d *decode.D) {
				d.FramedFn(int64(protocolsLength)*8, func(d *decode.D) {
					for !d.End() {
						d.FieldU8("ex_points_format", scalar.UintHex) // TODO: names
					}
				})
			})
		case extensionSupportedGroups:
			protocolsLength := d.FieldU16("supported_groups_length")
			d.FieldArray("supported_groups", func(d *decode.D) {
				d.FramedFn(int64(protocolsLength)*8, func(d *decode.D) {
					for !d.End() {
						d.FieldU16("supported_group", scalar.UintHex) // TODO: names
					}
				})
			})
		case extensionSignatureAlgorithms:
			protocolsLength := d.FieldU16("signature_algorithms_length")
			d.FieldArray("signature_algorithms", func(d *decode.D) {
				d.FramedFn(int64(protocolsLength)*8, func(d *decode.D) {
					for !d.End() {
						d.FieldStruct("signature_algorithm", func(d *decode.D) {
							d.FieldU8("hash", hashAlgorithmNames)
							d.FieldU8("signature", signatureAlgorithmNames)
						})
					}
				})
			})
		default:
			d.FieldRawLen("data", int64(length)*8)
		}
	})
}

func decodeTLSHandshake(d *decode.D, tc *tlsCtx) {
	msgType := d.FieldU8("type", handshakeMsgTypeNames)
	length := d.FieldU24("length")

	d.FramedFn(int64(length)*8, func(d *decode.D) {
		switch msgType {
		case handshakeMsgTypeHelloRequest:
			// TODO: nothing?
		case handshakeMsgTypeClientHello,
			handshakeMsgTypeServerHello:
			tc.version = d.FieldU16("version", versionNames, scalar.UintHex)
			copy(tc.random[:], d.PeekBytes(32))
			d.FieldStruct("random", func(d *decode.D) {
				d.FieldU32("gmt_unix_time", scalar.UintActualUnixTimeDescription(time.Second, time.RFC3339))
				d.FieldRawLen("random_bytes", 28*8)
			})

			sessionIDLength := d.FieldU8("session_id_length")
			d.FieldRawLen("session_id", int64(sessionIDLength)*8)

			if msgType == handshakeMsgTypeServerHello {
				tc.server.nextCipherSuit = d.FieldU16("cipher_suit", cipherNames, scalar.UintHex)
			} else {
				cipherSuitLength := d.FieldU16("cipher_suits_length")
				d.FramedFn(int64(cipherSuitLength)*8, func(d *decode.D) {
					d.FieldArray("cipher_suits", func(d *decode.D) {
						for !d.End() {
							d.FieldU16("cipher_suit", cipherNames, scalar.UintHex)
						}
					})
				})
			}

			if msgType == handshakeMsgTypeServerHello {
				tc.server.compressionMethod = d.FieldU8("compression_method", compressionMethodNames, scalar.UintHex)
			} else {
				compressionMethodLength := d.FieldU8("compression_methods_length")
				d.FramedFn(int64(compressionMethodLength)*8, func(d *decode.D) {
					d.FieldArray("compression_methods", func(d *decode.D) {
						for !d.End() {
							d.FieldU8("compression_method", compressionMethodNames, scalar.UintHex)
						}
					})
				})
			}

			// SSL v3 should have no extensions but we decode if there are bytes
			if d.BitsLeft() > 0 {
				extensionsLength := d.FieldU16("extensions_length")
				d.FramedFn(int64(extensionsLength)*8, func(d *decode.D) {
					d.FieldArray("extensions", func(d *decode.D) {
						for !d.End() {
							d.FieldStruct("extension", decodeTLSExtension)
						}
					})
				})
			}
		case handshakeMsgTypeCertificate:
			certificatesLength := d.FieldU24("certificates_length")
			d.FramedFn(int64(certificatesLength)*8, func(d *decode.D) {
				d.FieldArray("certificates", func(d *decode.D) {
					for !d.End() {
						d.FieldStruct("certificate", func(d *decode.D) {
							length := d.FieldU24("length")
							d.FieldFormatLen("data", int64(length)*8, &asn1BerGroup, nil)
						})
					}
				})
			})
		case handshakeMsgTypeClientKeyExchange,
			handshakeMsgTypeServerKeyExchange:
			// is decoded later in decodeTLSPostKeyExchange
			start := d.Pos()
			d.FieldRawLen("data", d.BitsLeft())
			dataV := d.FieldGet("data")
			tc.keyExchange = &keyExchange{
				msgType: msgType,
				d:       d,
				r:       ranges.Range{Start: d.Pos(), Len: d.Pos() - start},
				dataV:   dataV,
			}
		case handshakeMsgTypeFinished:
			d.FieldRawLen("verify_data", d.BitsLeft())
		case handshakeMsgTypeNewSessionTicket:
			d.FieldU32("lifetime_hint")
			ticketLength := d.FieldU16("ticket_length")
			d.FieldRawLen("ticket", int64(ticketLength)*8)
		default:
			d.FieldRawLen("data", d.BitsLeft())
		}
	})
}

func decodeTLSRecord(d *decode.D, tc *tlsCtx, isEncrypted bool) {
	recordStart := d.Pos()

	recordType := d.FieldU8("type", recordTypeNames, d.UintAssert(recordTypeValid...))
	d.FieldU16("version", versionNames, scalar.UintHex, d.UintAssert(versionValid...))
	length := d.FieldU16("length")
	d.FramedFn(int64(length)*8, func(d *decode.D) {
		if isEncrypted {
			d.FieldRawLen("encrypted_data", d.BitsLeft())
			// is decoded later in decodeTLSPostEncryptedRecords
			tc.encryptedRecords = append(tc.encryptedRecords, encryptedRecord{
				recordType: recordType,
				d:          d,
				r:          ranges.Range{Start: recordStart, Len: d.Pos() - recordStart},
			})
			return
		}

		d.FieldStruct("message", func(d *decode.D) {
			decodeTLSRecordMessage(d, tc, recordType)
		})
	})
}

func decodeTLSRecordMessage(d *decode.D, tc *tlsCtx, recordType uint64) {
	switch recordType {
	case recordTypeHandshake:
		decodeTLSHandshake(d, tc)
	case recordTypeChangeCipherSpec:
		d.FieldU8("type", changeCipherSpecTypeNames)
		tc.server.currentCipherSuit = tc.server.nextCipherSuit
		tc.isEncrypted = true
	case recordTypeApplicationData:
		d.FieldRawLen("data", d.BitsLeft())
	case recordTypeAlert:
		d.FieldU8("level", alertLevelNames)
		d.FieldU8("description", alertNames)
	default:
		d.FieldRawLen("data", d.BitsLeft())
	}
}

func decodeTLSPostEncryptedRecords(rootD *decode.D, tc *tlsCtx, kl keylog.Map) {
	// to decrypt tls we need:
	//  - client random to look up shared master secret
	//  - client and server random to generate cipher iv/key in both directions
	masterSecret, _ := kl.Lookup(keylog.ClientRandom, tc.clientCtx.random)
	if masterSecret == nil {
		// TODO: info/warn?
		return
	}

	td := tlsdecrypt.Decryptor{
		IsClient:     tc == tc.clientCtx,
		Version:      int(tc.serverCtx.version),
		CipherSuite:  int(tc.serverCtx.server.currentCipherSuit),
		MasterSecret: masterSecret,
		ClientRandom: tc.clientCtx.random[:],
		ServerRandom: tc.serverCtx.random[:],
	}

	applicationStream := &bytes.Buffer{}
	plainBuf := &bytes.Buffer{}
	var uncompressR io.Reader
	hasApplicationStream := false

	for _, r := range tc.encryptedRecords {
		encryptedRecord := r.d.ReadAllBits(rootD.BitBufRange(r.r.Start, r.r.Len))
		plain, decryptErr := td.Decrypt(encryptedRecord)
		if decryptErr != nil {
			// TODO: warn
			// log.Printf("err: %#+v\n", decryptErr)
			continue
		}

		plainBuf.Write(plain)

		// happens after plainBuf write as uncompressor init might read input
		if uncompressR == nil {
			switch tc.serverCtx.server.compressionMethod {
			case compressionMethodNull:
				uncompressR = plainBuf
			case compressionMethodDeflate:
				var uncompressRErr error
				uncompressR, uncompressRErr = rezlib.NewReader(plainBuf)
				if uncompressRErr != nil {
					// TODO: warn?
					continue
				}
			default:
				// TODO: how to inform? option?
				// rootD.Fatalf("compression method %d not supported", tc.serverCtx.server.compressionMethod)
				continue
			}
		}

		plainUncomp, plainUncompErr := io.ReadAll(uncompressR)
		if plainUncompErr != nil {
			continue
		}

		bbr := bitio.NewBitReader(plainUncomp, -1)
		// application data handled differently to get data as .message
		if r.recordType == recordTypeApplicationData {
			applicationStream.Write(plainUncomp)
			hasApplicationStream = true
			r.d.FieldRootBitBuf("message", bbr)
		} else {
			r.d.FieldStructRootBitBufFn("message", bbr, func(d *decode.D) {
				decodeTLSRecordMessage(d, tc, r.recordType)
			})
		}
	}

	if hasApplicationStream {
		applicationBytes := applicationStream.Bytes()
		rootD.FieldRootBitBuf("stream", bitio.NewBitReader(applicationBytes, -1))
	}
}

func decodeTLSPostKeyExchange(tc *tlsCtx) {
	ke := tc.keyExchange
	if ke == nil {
		return
	}

	// TODO: better way to track version and cipher changes?

	// only for TLS 1.0-1.2 for now
	switch tc.serverCtx.version {
	case versionTLS_1_0,
		versionTLS_1_1,
		versionTLS_1_2:
	default:
		return
	}
	s, ok := ciphersuites.Suits[int(tc.serverCtx.server.currentCipherSuit)]
	if !ok {
		return
	}

	// TODO: find a better way
	ke.d.SeekRel(-ke.r.Len)

	replaceData := true

	switch ke.msgType {
	case handshakeMsgTypeClientKeyExchange:
		// struct {
		// 	select (KeyExchangeAlgorithm) {
		// 		case rsa:
		// 			EncryptedPreMasterSecret;
		// 		case dhe_dss:
		// 		case dhe_rsa:
		// 		case dh_dss:
		// 		case dh_rsa:
		// 		case dh_anon:
		// 			ClientDiffieHellmanPublic;
		//      case ec_diffie_hellman:
		//          ClientECDiffieHellmanPublic;
		// 	} exchange_keys;
		// } ClientKeyExchange;
		switch s.KeyAgreement {
		// TODO: ssl/tls difference?
		// ciphersuites.DH_anon_EXPORT:
		//ciphersuites.RSA_PSK,
		//ciphersuites.RSA_EXPORT,
		//ciphersuites.RSA_FIPS:
		case
			ciphersuites.RSA:
			ke.d.FieldStruct("encrypted_premaster", func(d *decode.D) {
				length := d.FieldU16("length")
				d.FieldRawLen("data", int64(length)*8)
			})
		case ciphersuites.DHE_DSS,
			ciphersuites.DHE_RSA,
			ciphersuites.DH_DSS,
			ciphersuites.DH_RSA,
			ciphersuites.DH_anon:
			ke.d.FieldStruct("public", func(d *decode.D) {
				length := d.FieldU16("length")
				d.FieldRawLen("data", int64(length)*8)
			})
		case ciphersuites.ECDH_ECDSA,
			ciphersuites.ECDH_RSA,
			ciphersuites.ECDH_anon,
			ciphersuites.ECDHE_ECDSA,
			ciphersuites.ECDHE_PSK,
			ciphersuites.ECDHE_RSA:
			ke.d.FieldStruct("public", func(d *decode.D) {
				length := d.FieldU8("length")
				d.FieldRawLen("data", int64(length)*8)
			})
		default:
			replaceData = false
		}
	case handshakeMsgTypeServerKeyExchange:
		// struct {
		// 	select (KeyExchangeAlgorithm) {
		// 		case dh_anon:
		// 			ServerDHParams params;
		// 		case dhe_dss:
		// 		case dhe_rsa:
		// 			ServerDHParams params;
		// 			digitally-signed struct {
		// 				opaque client_random[32];
		// 				opaque server_random[32];
		// 				ServerDHParams params;
		// 			} signed_params;
		// 		case rsa:
		// 		case dh_dss:
		// 		case dh_rsa:
		// 			struct {} ;
		// 		   /* message is omitted for rsa, dh_dss, and dh_rsa */
		// 		/* may be extended, e.g., for ECDH -- see [TLSECC] */
		// 	};
		// } ServerKeyExchange;
		switch s.KeyAgreement {
		//ciphersuites.ECDHE_PSK,
		// case ciphersuites.RSA:
		case ciphersuites.ECDH_ECDSA,
			ciphersuites.ECDH_RSA,
			ciphersuites.ECDH_anon,
			ciphersuites.ECDHE_ECDSA,
			ciphersuites.ECDHE_RSA:
			ke.d.FieldStruct("curve_params", func(d *decode.D) {
				curveType := d.FieldU8("curve_type")
				// TODO: named 3=named_curve
				switch curveType {
				case 3:
					d.FieldU16("named_curve")
				}
			})
			ke.d.FieldStruct("public", func(d *decode.D) {
				length := d.FieldU8("length")
				d.FieldRawLen("data", int64(length)*8)
			})
			ke.d.FieldStruct("signature_algorithm", func(d *decode.D) {
				d.FieldU8("hash", hashAlgorithmNames)
				d.FieldU8("signature", signatureAlgorithmNames)
				length := d.FieldU16("length")
				d.FieldRawLen("data", int64(length)*8)
			})
		default:
			replaceData = false
		}
	default:
		panic(fmt.Sprintf("unknown ke type %d", ke.msgType))
	}

	if replaceData {
		if err := ke.dataV.Remove(); err != nil {
			panic(err)
		}
	}
}

// SSL v2 compatible header
func decodeV2ClientHello(d *decode.D, tc *tlsCtx) {
	// TODO: header length/padding?
	d.FieldU16("length", scalar.UintActualFn(func(a uint64) uint64 { return a & 0x7fff }))
	d.FieldU8("type")
	d.FieldU16("tls_version")
	cipherLength := d.FieldU16("cipher_spec_length")
	sessionIDLength := d.FieldU16("session_id_length")
	challengeLength := d.FieldU16("challenge_length")

	d.FramedFn(int64(cipherLength)*8, func(d *decode.D) {
		d.FieldArray("cipher_specs", func(d *decode.D) {
			for !d.End() {
				d.FieldU24("cipher_spec", cipherNames, scalar.UintHex)
			}
		})
	})
	d.FieldRawLen("session_id", int64(sessionIDLength)*8)
	copy(tc.random[:], d.PeekBytes(int(challengeLength)))
	d.FieldRawLen("challenge", int64(challengeLength)*8)
}

func decodeTLS(d *decode.D) any {
	var ti format.TLS_In
	d.ArgAs(&ti)

	isClient := false

	var tsi format.TCP_Stream_In
	if d.ArgAs(&tsi) {
		if !tsi.HasStart {
			d.Fatalf("tls requires start of byte stream")
		}
		isClient = tsi.IsClient
	}

	tc := &tlsCtx{
		rootD: d,
	}
	tc.server.currentCipherSuit = ciphersuites.TLS_NULL_WITH_NULL_NULL
	tc.server.nextCipherSuit = ciphersuites.TLS_NULL_WITH_NULL_NULL

	firstByte := d.PeekUintBits(8)
	if firstByte&0x80 != 0 {
		d.FieldStruct("ssl_v2", func(d *decode.D) {
			decodeV2ClientHello(d, tc)
		})
	}
	recordsDecoded := 0
	d.FieldArray("records", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("record", func(d *decode.D) {
				decodeTLSRecord(d, tc, tc.isEncrypted)
				recordsDecoded++
			})
		}
	})
	// as we're in the tcp group we should at least decode one record
	if recordsDecoded == 0 {
		d.Fatalf("no records found")
	}

	// client side will do post for both
	if !isClient {
		return format.TCP_Stream_Out{InArg: tc}
	}

	return format.TCP_Stream_Out{
		PostFn: func(peerIn any) {
			// peerIn will be the other peers outArg, the server *tlsCtx
			clientTc := tc
			serverTc, serverTcOk := peerIn.(*tlsCtx)
			if !serverTcOk {
				panic(fmt.Sprintf("tls PostFn in not *tlsCtx %+#v", peerIn))
			}

			tc.clientCtx = clientTc
			tc.serverCtx = serverTc
			serverTc.clientCtx = clientTc
			serverTc.serverCtx = serverTc

			decodeTLSPostKeyExchange(clientTc)
			decodeTLSPostKeyExchange(serverTc)

			if ti.Keylog == "" {
				return
			}

			km, err := keylog.Parse(ti.Keylog)
			if err != nil {
				d.Fatalf("failed to parse keylog: %s", err)
			}

			decodeTLSPostEncryptedRecords(clientTc.rootD, clientTc, km)
			decodeTLSPostEncryptedRecords(serverTc.rootD, serverTc, km)
		},
		InArg: tc,
	}
}
