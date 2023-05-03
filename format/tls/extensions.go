package tls

import "github.com/wader/fq/pkg/scalar"

// https://www.iana.org/assignments/tls-extensiontype-values/tls-extensiontype-values-1.csv

const (
	// fq -d csv -r '.[1:][] | select(.[1] | test("unassigned|reserved";"i") | not) | "extension\(.[1]) = \(.[0] | .[0:4] + . [7:] | ascii_downcase)"' tls-extensiontype-values-1.csv
	extensionServerName                          = 0
	extensionMaxFragmentLength                   = 1
	extensionClientCertificateUrl                = 2
	extensionTrustedCaKeys                       = 3
	extensionTruncatedHmac                       = 4
	extensionStatusRequest                       = 5
	extensionUserMapping                         = 6
	extensionClientAuthz                         = 7
	extensionServerAuthz                         = 8
	extensionCertType                            = 9
	extensionSupportedGroups                     = 10
	extensionEcPointFormats                      = 11
	extensionSrp                                 = 12
	extensionSignatureAlgorithms                 = 13
	extensionUseSrtp                             = 14
	extensionHeartbeat                           = 15
	extensionApplicationLayerProtocolNegotiation = 16
	extensionStatusRequestV2                     = 17
	extensionSignedCertificateTimestamp          = 18
	extensionClientCertificateType               = 19
	extensionServerCertificateType               = 20
	extensionPadding                             = 21
	extensionEncryptThenMac                      = 22
	extensionExtendedMasterSecret                = 23
	extensionTokenBinding                        = 24
	extensionCachedInfo                          = 25
	extensionTlsLts                              = 26
	extensionCompressCertificate                 = 27
	extensionRecordSizeLimit                     = 28
	extensionPwdProtect                          = 29
	extensionPwdClear                            = 30
	extensionPasswordSalt                        = 31
	extensionTicketPinning                       = 32
	extensionTlsCertWithExternPsk                = 33
	extensionDelegatedCredentials                = 34
	extensionSessionTicket                       = 35
	extensionTLMSP                               = 36
	extensionTLMSPProxying                       = 37
	extensionTLMSPDelegate                       = 38
	extensionSupportedEktCiphers                 = 39
	extensionPreSharedKey                        = 41
	extensionEarlyData                           = 42
	extensionSupportedVersions                   = 43
	extensionCookie                              = 44
	extensionPskKeyExchangeModes                 = 45
	extensionCertificateAuthorities              = 47
	extensionOidFilters                          = 48
	extensionPostHandshakeAuth                   = 49
	extensionSignatureAlgorithmsCert             = 50
	extensionKeyShare                            = 51
	extensionTransparencyInfo                    = 52
	extensionConnectionId53                      = 53 // deprecated
	extensionConnectionId                        = 54
	extensionExternalIdHash                      = 55
	extensionExternalSessionId                   = 56
	extensionQuicTransportParameters             = 57
	extensionTicketRequest                       = 58
	extensionDnssecChain                         = 59
	extensionRenegotiationInfo                   = 65281
)

var extensionNames = scalar.UintMapSymStr{
	// fq -d csv -r '.[1:][] | select(.[1] | test("unassigned|reserved";"i") | not) | "extension\(.[1] | gsub(`(?<c>(?:^|_).)`; .c[-1:] | ascii_upcase)) = \(.[0] | .[0:4] + . [7:] | ascii_downcase)"' tls-extensiontype-values-1.csv
	extensionServerName:                          "server_name",
	extensionMaxFragmentLength:                   "max_fragment_length",
	extensionClientCertificateUrl:                "client_certificate_url",
	extensionTrustedCaKeys:                       "trusted_ca_keys",
	extensionTruncatedHmac:                       "truncated_hmac",
	extensionStatusRequest:                       "status_request",
	extensionUserMapping:                         "user_mapping",
	extensionClientAuthz:                         "client_authz",
	extensionServerAuthz:                         "server_authz",
	extensionCertType:                            "cert_type",
	extensionSupportedGroups:                     "supported_groups",
	extensionEcPointFormats:                      "ec_point_formats",
	extensionSrp:                                 "srp",
	extensionSignatureAlgorithms:                 "signature_algorithms",
	extensionUseSrtp:                             "use_srtp",
	extensionHeartbeat:                           "heartbeat",
	extensionApplicationLayerProtocolNegotiation: "application_layer_protocol_negotiation",
	extensionStatusRequestV2:                     "status_request_v2",
	extensionSignedCertificateTimestamp:          "signed_certificate_timestamp",
	extensionClientCertificateType:               "client_certificate_type",
	extensionServerCertificateType:               "server_certificate_type",
	extensionPadding:                             "padding",
	extensionEncryptThenMac:                      "encrypt_then_mac",
	extensionExtendedMasterSecret:                "extended_master_secret",
	extensionTokenBinding:                        "token_binding",
	extensionCachedInfo:                          "cached_info",
	extensionTlsLts:                              "tls_lts",
	extensionCompressCertificate:                 "compress_certificate",
	extensionRecordSizeLimit:                     "record_size_limit",
	extensionPwdProtect:                          "pwd_protect",
	extensionPwdClear:                            "pwd_clear",
	extensionPasswordSalt:                        "password_salt",
	extensionTicketPinning:                       "ticket_pinning",
	extensionTlsCertWithExternPsk:                "tls_cert_with_extern_psk",
	extensionDelegatedCredentials:                "delegated_credentials",
	extensionSessionTicket:                       "session_ticket",
	extensionTLMSP:                               "tlmsp",
	extensionTLMSPProxying:                       "tlmsp_proxying",
	extensionTLMSPDelegate:                       "tlmsp_delegate",
	extensionSupportedEktCiphers:                 "supported_ekt_ciphers",
	extensionPreSharedKey:                        "pre_shared_key",
	extensionEarlyData:                           "early_data",
	extensionSupportedVersions:                   "supported_versions",
	extensionCookie:                              "cookie",
	extensionPskKeyExchangeModes:                 "psk_key_exchange_modes",
	extensionCertificateAuthorities:              "certificate_authorities",
	extensionOidFilters:                          "oid_filters",
	extensionPostHandshakeAuth:                   "post_handshake_auth",
	extensionSignatureAlgorithmsCert:             "signature_algorithms_cert",
	extensionKeyShare:                            "key_share",
	extensionTransparencyInfo:                    "transparency_info",
	extensionConnectionId53:                      "connection_id53",
	extensionConnectionId:                        "connection_id",
	extensionExternalIdHash:                      "external_id_hash",
	extensionExternalSessionId:                   "external_session_id",
	extensionQuicTransportParameters:             "quic_transport_parameters",
	extensionTicketRequest:                       "ticket_request",
	extensionDnssecChain:                         "dnssec_chain",
	extensionRenegotiationInfo:                   "renegotiation_info",
}
