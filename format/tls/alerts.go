package tls

import "github.com/wader/fq/pkg/scalar"

const (
	alertLevelWarning = 1
	alertLevelError   = 2
)

var alertLevelNames = scalar.UintMapSymStr{
	alertLevelWarning: "warning",
	alertLevelError:   "error",
}

const (
	alertCloseNotify                  = 0
	alertUnexpectedMessage            = 10
	alertBadRecordMAC                 = 20
	alertDecryptionFailed             = 21
	alertRecordOverflow               = 22
	alertDecompressionFailure         = 30
	alertHandshakeFailure             = 40
	alertBadCertificate               = 42
	alertUnsupportedCertificate       = 43
	alertCertificateRevoked           = 44
	alertCertificateExpired           = 45
	alertCertificateUnknown           = 46
	alertIllegalParameter             = 47
	alertUnknownCA                    = 48
	alertAccessDenied                 = 49
	alertDecodeError                  = 50
	alertDecryptError                 = 51
	alertExportRestriction            = 60
	alertProtocolVersion              = 70
	alertInsufficientSecurity         = 71
	alertInternalError                = 80
	alertInappropriateFallback        = 86
	alertUserCanceled                 = 90
	alertNoRenegotiation              = 100
	alertMissingExtension             = 109
	alertUnsupportedExtension         = 110
	alertCertificateUnobtainable      = 111
	alertUnrecognizedName             = 112
	alertBadCertificateStatusResponse = 113
	alertBadCertificateHashValue      = 114
	alertUnknownPSKIdentity           = 115
	alertCertificateRequired          = 116
	alertNoApplicationProtocol        = 120
)

var alertNames = scalar.UintMapSymStr{
	alertCloseNotify:                  "close_notify",
	alertUnexpectedMessage:            "unexpected_message",
	alertBadRecordMAC:                 "bad_record_mac",
	alertDecryptionFailed:             "decryption_failed",
	alertRecordOverflow:               "record_overflow",
	alertDecompressionFailure:         "decompression_failure",
	alertHandshakeFailure:             "handshake_failure",
	alertBadCertificate:               "bad_certificate",
	alertUnsupportedCertificate:       "unsupported_certificate",
	alertCertificateRevoked:           "revoked_certificate",
	alertCertificateExpired:           "expired_certificate",
	alertCertificateUnknown:           "unknown_certificate",
	alertIllegalParameter:             "illegal_parameter",
	alertUnknownCA:                    "unknown_certificate_authority",
	alertAccessDenied:                 "access_denied",
	alertDecodeError:                  "error_decoding_message",
	alertDecryptError:                 "error_decrypting_message",
	alertExportRestriction:            "export_restriction",
	alertProtocolVersion:              "protocol_version_not_supported",
	alertInsufficientSecurity:         "insufficient_security_level",
	alertInternalError:                "internal_error",
	alertInappropriateFallback:        "inappropriate_fallback",
	alertUserCanceled:                 "user_canceled",
	alertNoRenegotiation:              "no_renegotiation",
	alertMissingExtension:             "missing_extension",
	alertUnsupportedExtension:         "unsupported_extension",
	alertCertificateUnobtainable:      "certificate_unobtainable",
	alertUnrecognizedName:             "unrecognized_name",
	alertBadCertificateStatusResponse: "bad_certificate_status_response",
	alertBadCertificateHashValue:      "bad_certificate_hash_value",
	alertUnknownPSKIdentity:           "unknown_PSK_identity",
	alertCertificateRequired:          "certificate_required",
	alertNoApplicationProtocol:        "no_application_protocol",
}
