package pcap

import "github.com/wader/fq/pkg/decode"

//nolint:revive
const (
	LINKTYPE_NULL                       = 0
	LINKTYPE_ETHERNET                   = 1
	LINKTYPE_AX25                       = 3
	LINKTYPE_IEEE802_5                  = 6
	LINKTYPE_ARCNET_BSD                 = 7
	LINKTYPE_SLIP                       = 8
	LINKTYPE_PPP                        = 9
	LINKTYPE_FDDI                       = 10
	LINKTYPE_PPP_HDLC                   = 50
	LINKTYPE_PPP_ETHER                  = 51
	LINKTYPE_ATM_RFC1483                = 100
	LINKTYPE_RAW                        = 101
	LINKTYPE_C_HDLC                     = 104
	LINKTYPE_IEEE802_11                 = 105
	LINKTYPE_FRELAY                     = 107
	LINKTYPE_LOOP                       = 108
	LINKTYPE_LINUX_SLL                  = 113
	LINKTYPE_LTALK                      = 114
	LINKTYPE_PFLOG                      = 117
	LINKTYPE_IEEE802_11_PRISM           = 119
	LINKTYPE_IP_OVER_FC                 = 122
	LINKTYPE_SUNATM                     = 123
	LINKTYPE_IEEE802_11_RADIOTAP        = 127
	LINKTYPE_ARCNET_LINUX               = 129
	LINKTYPE_APPLE_IP_OVER_IEEE1394     = 138
	LINKTYPE_MTP2_WITH_PHDR             = 139
	LINKTYPE_MTP2                       = 140
	LINKTYPE_MTP3                       = 141
	LINKTYPE_SCCP                       = 142
	LINKTYPE_DOCSIS                     = 143
	LINKTYPE_LINUX_IRDA                 = 144
	LINKTYPE_USER0                      = 147
	LINKTYPE_USER1                      = 148
	LINKTYPE_USER2                      = 149
	LINKTYPE_USER3                      = 150
	LINKTYPE_USER4                      = 151
	LINKTYPE_USER5                      = 152
	LINKTYPE_USER6                      = 153
	LINKTYPE_USER7                      = 154
	LINKTYPE_USER8                      = 155
	LINKTYPE_USER9                      = 156
	LINKTYPE_USER10                     = 157
	LINKTYPE_USER11                     = 158
	LINKTYPE_USER12                     = 159
	LINKTYPE_USER13                     = 160
	LINKTYPE_USER14                     = 161
	LINKTYPE_USER15                     = 162
	LINKTYPE_IEEE802_11_AVS             = 163
	LINKTYPE_BACNET_MS_TP               = 165
	LINKTYPE_PPP_PPPD                   = 166
	LINKTYPE_GPRS_LLC                   = 169
	LINKTYPE_GPF_T                      = 170
	LINKTYPE_GPF_F                      = 171
	LINKTYPE_LINUX_LAPD                 = 177
	LINKTYPE_MFR                        = 182
	LINKTYPE_BLUETOOTH_HCI_H4           = 187
	LINKTYPE_USB_LINUX                  = 189
	LINKTYPE_PPI                        = 192
	LINKTYPE_IEEE802_15_4_WITHFCS       = 195
	LINKTYPE_SITA                       = 196
	LINKTYPE_ERF                        = 197
	LINKTYPE_BLUETOOTH_HCI_H4_WITH_PHDR = 201
	LINKTYPE_AX25_KISS                  = 202
	LINKTYPE_LAPD                       = 203
	LINKTYPE_PPP_WITH_DIR               = 204
	LINKTYPE_C_HDLC_WITH_DIR            = 205
	LINKTYPE_FRELAY_WITH_DIR            = 206
	LINKTYPE_LAPB_WITH_DIR              = 207
	LINKTYPE_IPMB_LINUX                 = 209
	LINKTYPE_FLEXRAY                    = 210
	LINKTYPE_LIN                        = 212
	LINKTYPE_IEEE802_15_4_NONASK_PHY    = 215
	LINKTYPE_USB_LINUX_MMAPPED          = 220
	LINKTYPE_FC_2                       = 224
	LINKTYPE_FC_2_WITH_FRAME_DELIMS     = 225
	LINKTYPE_IPNET                      = 226
	LINKTYPE_CAN_SOCKETCAN              = 227
	LINKTYPE_IPV4                       = 228
	LINKTYPE_IPV6                       = 229
	LINKTYPE_IEEE802_15_4_NOFCS         = 230
	LINKTYPE_DBUS                       = 231
	LINKTYPE_DVB_CI                     = 235
	LINKTYPE_MUX27010                   = 236
	LINKTYPE_STANAG_5066_D_PDU          = 237
	LINKTYPE_NFLOG                      = 239
	LINKTYPE_NETANALYZER                = 240
	LINKTYPE_NETANALYZER_TRANSPARENT    = 241
	LINKTYPE_IPOIB                      = 242
	LINKTYPE_MPEG_2_TS                  = 243
	LINKTYPE_NG40                       = 244
	LINKTYPE_NFC_LLCP                   = 245
	LINKTYPE_INFINIBAND                 = 247
	LINKTYPE_SCTP                       = 248
	LINKTYPE_USBPCAP                    = 249
	LINKTYPE_RTAC_SERIAL                = 250
	LINKTYPE_BLUETOOTH_LE_LL            = 251
	LINKTYPE_NETLINK                    = 253
	LINKTYPE_BLUETOOTH_LINUX_MONITOR    = 254
	LINKTYPE_BLUETOOTH_BREDR_BB         = 255
	LINKTYPE_BLUETOOTH_LE_LL_WITH_PHDR  = 256
	LINKTYPE_PROFIBUS_DL                = 257
	LINKTYPE_PKTAP                      = 258
	LINKTYPE_EPON                       = 259
	LINKTYPE_IPMI_HPM_2                 = 260
	LINKTYPE_ZWAVE_R1_R2                = 261
	LINKTYPE_ZWAVE_R3                   = 262
	LINKTYPE_WATTSTOPPER_DLM            = 263
	LINKTYPE_ISO_14443                  = 264
	LINKTYPE_RDS                        = 265
	LINKTYPE_USB_DARWIN                 = 266
	LINKTYPE_SDLC                       = 268
	LINKTYPE_LORATAP                    = 270
	LINKTYPE_VSOCK                      = 271
	LINKTYPE_NORDIC_BLE                 = 272
	LINKTYPE_DOCSIS31_XRA31             = 273
	LINKTYPE_ETHERNET_MPACKET           = 274
	LINKTYPE_DISPLAYPORT_AUX            = 275
	LINKTYPE_LINUX_SLL2                 = 276
	LINKTYPE_OPENVIZSLA                 = 278
	LINKTYPE_EBHSCR                     = 279
	LINKTYPE_VPP_DISPATCH               = 280
	LINKTYPE_DSA_TAG_BRCM               = 281
	LINKTYPE_DSA_TAG_BRCM_PREPEND       = 282
	LINKTYPE_IEEE802_15_4_TAP           = 283
	LINKTYPE_DSA_TAG_DSA                = 284
	LINKTYPE_DSA_TAG_EDSA               = 285
	LINKTYPE_ELEE                       = 286
	LINKTYPE_Z_WAVE_SERIAL              = 287
	LINKTYPE_USB_2_0                    = 288
	LINKTYPE_ATSC_ALP                   = 289
	LINKTYPE_ETW                        = 290
)

// from https://www.tcpdump.org/linktypes.html
// TODO cleanup
var linkTypeMap = decode.UToScalar{
	LINKTYPE_NULL:                       {Sym: "null", Description: `BSD loopback encapsulation`},
	LINKTYPE_ETHERNET:                   {Sym: "ethernet", Description: `IEEE 802.3 Ethernet`},
	LINKTYPE_AX25:                       {Sym: "ax25", Description: `AX.25 packet, with nothing preceding it.`},
	LINKTYPE_IEEE802_5:                  {Sym: "ieee802_5", Description: `IEEE 802.5 Token Ring`},
	LINKTYPE_ARCNET_BSD:                 {Sym: "arcnet_bsd", Description: `ARCNET Data Packets`},
	LINKTYPE_SLIP:                       {Sym: "slip", Description: `SLIP, encapsulated with a LINKTYPE_SLIP header.`},
	LINKTYPE_PPP:                        {Sym: "ppp", Description: `PPP`},
	LINKTYPE_FDDI:                       {Sym: "fddi", Description: `FDDI`},
	LINKTYPE_PPP_HDLC:                   {Sym: "ppp_hdlc", Description: `PPP in HDLC-like framing`},
	LINKTYPE_PPP_ETHER:                  {Sym: "ppp_ether", Description: `PPPoE`},
	LINKTYPE_ATM_RFC1483:                {Sym: "atm_rfc1483", Description: `RFC 1483 LLC/SNAP-encapsulated ATM; the packet begins with an ISO 8802-2 (formerly known as IEEE 802.2) LLC header.`},
	LINKTYPE_RAW:                        {Sym: "raw", Description: `Raw IP; the packet begins with an IPv4 or IPv6 header, with the "version" field of the header indicating whether it's an IPv4 or IPv6 header.`},
	LINKTYPE_C_HDLC:                     {Sym: "c_hdlc", Description: `Cisco PPP with HDLC framing, as per section 4.3.1 of RFC 1547.`},
	LINKTYPE_IEEE802_11:                 {Sym: "ieee802_11", Description: `IEEE 802.11 wireless LAN.`},
	LINKTYPE_FRELAY:                     {Sym: "frelay", Description: `Frame Relay LAPF frames, beginning with a ITU-T Recommendation Q.922 LAPF header starting with the address field, and without an FCS at the end of the frame.`},
	LINKTYPE_LOOP:                       {Sym: "loop", Description: `OpenBSD loopback encapsulation; the link-layer header is a 4-byte field, in network byte order, containing a value of 2 for IPv4 packets, a value of either 24, 28, or 30 for IPv6 packets, a value of 7 for OSI packets, or a value of 23 for IPX packets. All of the IPv6 values correspond to IPv6 packets; code reading files should check for all of them.`},
	LINKTYPE_LINUX_SLL:                  {Sym: "linux_sll", Description: `Linux "cooked" capture encapsulation.`},
	LINKTYPE_LTALK:                      {Sym: "ltalk", Description: `Apple LocalTalk; the packet begins with an AppleTalk LocalTalk Link Access Protocol header, as described in chapter 1 of Inside AppleTalk, Second Edition.`},
	LINKTYPE_PFLOG:                      {Sym: "pflog", Description: `OpenBSD pflog; the link-layer header contains a "struct pfloghdr" structure, as defined by the host on which the file was saved. (This differs from operating system to operating system and release to release; there is nothing in the file to indicate what the layout of that structure is.)`},
	LINKTYPE_IEEE802_11_PRISM:           {Sym: "ieee802_11_prism", Description: `Prism monitor mode information followed by an 802.11 header.`},
	LINKTYPE_IP_OVER_FC:                 {Sym: "ip_over_fc", Description: `RFC 2625 IP-over-Fibre Channel, with the link-layer header being the Network_Header as described in that RFC.`},
	LINKTYPE_SUNATM:                     {Sym: "sunatm", Description: `ATM traffic, encapsulated as per the scheme used by SunATM devices.`},
	LINKTYPE_IEEE802_11_RADIOTAP:        {Sym: "ieee802_11_radiotap", Description: `Radiotap link-layer information followed by an 802.11 header.`},
	LINKTYPE_ARCNET_LINUX:               {Sym: "arcnet_linux", Description: `ARCNET Data Packets, as described by the ARCNET Trade Association standard ATA 878.1-1999, but without the Starting Delimiter, Information Length, or Frame Check Sequence fields, with only the first ISU of the Destination Identifier, and with an extra two-ISU "offset" field following the Destination Identifier. For most packet types, ARCNET Trade Association draft standard ATA 878.2 is also used; however, no exception frames are supplied, and reassembled frames, rather than fragments, are supplied. See also RFC 1051 and RFC 1201; for RFC 1051 frames, ATA 878.2 is not used.`},
	LINKTYPE_APPLE_IP_OVER_IEEE1394:     {Sym: "apple_ip_over_ieee1394", Description: `Apple IP-over-IEEE 1394 cooked header.`},
	LINKTYPE_MTP2_WITH_PHDR:             {Sym: "mtp2_with_phdr", Description: `Signaling System 7 Message Transfer Part Level 2, as specified by ITU-T Recommendation Q.703, preceded by a pseudo-header.`},
	LINKTYPE_MTP2:                       {Sym: "mtp2", Description: `Signaling System 7 Message Transfer Part Level 2, as specified by ITU-T Recommendation Q.703.`},
	LINKTYPE_MTP3:                       {Sym: "mtp3", Description: `Signaling System 7 Message Transfer Part Level 3, as specified by ITU-T Recommendation Q.704, with no MTP2 header preceding the MTP3 packet.`},
	LINKTYPE_SCCP:                       {Sym: "sccp", Description: `Signaling System 7 Signalling Connection Control Part, as specified by ITU-T Recommendation Q.711, ITU-T Recommendation Q.712, ITU-T Recommendation Q.713, and ITU-T Recommendation Q.714, with no MTP3 or MTP2 headers preceding the SCCP packet.`},
	LINKTYPE_DOCSIS:                     {Sym: "docsis", Description: `DOCSIS MAC frames, as described by the DOCSIS 3.1 MAC and Upper Layer Protocols Interface Specification or earlier specifications for MAC frames.`},
	LINKTYPE_LINUX_IRDA:                 {Sym: "linux_irda", Description: `Linux-IrDA packets, with a LINKTYPE_LINUX_IRDA header, with the payload for IrDA frames beginning with by the IrLAP header as defined by IrDA Data Specifications, including the IrDA Link Access Protocol specification.`},
	LINKTYPE_USER0:                      {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER1:                      {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER2:                      {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER3:                      {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER4:                      {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER5:                      {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER6:                      {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER7:                      {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER8:                      {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER9:                      {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER10:                     {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER11:                     {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER12:                     {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER13:                     {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER14:                     {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_USER15:                     {Sym: "user0", Description: `Reserved for private use`},
	LINKTYPE_IEEE802_11_AVS:             {Sym: "ieee802_11_avs", Description: `AVS monitor mode information followed by an 802.11 header.`},
	LINKTYPE_BACNET_MS_TP:               {Sym: "bacnet_ms_tp", Description: `BACnet MS/TP frames, as specified by section 9.3 MS/TP Frame Format of ANSI/ASHRAE Standard 135, BACnet® - A Data Communication Protocol for Building Automation and Control Networks, including the preamble and, if present, the Data CRC.`},
	LINKTYPE_PPP_PPPD:                   {Sym: "ppp_pppd", Description: `PPP in HDLC-like encapsulation, like LINKTYPE_PPP_HDLC, but with the 0xff address byte replaced by a direction indication - 0x00 for incoming and 0x01 for outgoing.`},
	LINKTYPE_GPRS_LLC:                   {Sym: "gprs_llc", Description: `General Packet Radio Service Logical Link Control, as defined by 3GPP TS 04.64.`},
	LINKTYPE_GPF_T:                      {Sym: "gpf_t", Description: `Transparent-mapped generic framing procedure, as specified by ITU-T Recommendation G.7041/Y.1303.`},
	LINKTYPE_GPF_F:                      {Sym: "gpf_f", Description: `Frame-mapped generic framing procedure, as specified by ITU-T Recommendation G.7041/Y.1303.`},
	LINKTYPE_LINUX_LAPD:                 {Sym: "linux_lapd", Description: `Link Access Procedures on the D Channel (LAPD) frames, as specified by ITU-T Recommendation Q.920 and ITU-T Recommendation Q.921, captured via vISDN, with a LINKTYPE_LINUX_LAPD header, followed by the Q.921 frame, starting with the address field.`},
	LINKTYPE_MFR:                        {Sym: "mfr", Description: `FRF.16.1 Multi-Link Frame Relay frames, beginning with an FRF.12 Interface fragmentation format fragmentation header.`},
	LINKTYPE_BLUETOOTH_HCI_H4:           {Sym: "bluetooth_hci_h4", Description: `Bluetooth HCI UART transport layer; the frame contains an HCI packet indicator byte, as specified by the UART Transport Layer portion of the most recent Bluetooth Core specification, followed by an HCI packet of the specified packet type, as specified by the Host Controller Interface Functional Specification portion of the most recent Bluetooth Core Specification.`},
	LINKTYPE_USB_LINUX:                  {Sym: "usb_linux", Description: `USB packets, beginning with a Linux USB header, as specified by the struct usbmon_packet in the Documentation/usb/usbmon.txt file in the Linux source tree. Only the first 48 bytes of that header are present. All fields in the header are in host byte order. When performing a live capture, the host byte order is the byte order of the machine on which the packets are captured. When reading a pcap file, the byte order is the byte order for the file, as specified by the file's magic number; when reading a pcapng file, the byte order is the byte order for the section of the pcapng file, as specified by the Section Header Block.`},
	LINKTYPE_PPI:                        {Sym: "ppi", Description: `Per-Packet Information information, as specified by the Per-Packet Information Header Specification, followed by a packet with the LINKTYPE_ value specified by the pph_dlt field of that header.`},
	LINKTYPE_IEEE802_15_4_WITHFCS:       {Sym: "ieee802_15_4_withfcs", Description: `IEEE 802.15.4 Low-Rate Wireless Networks, with each packet having the FCS at the end of the frame.`},
	LINKTYPE_SITA:                       {Sym: "sita", Description: `Various link-layer types, with a pseudo-header, for SITA.`},
	LINKTYPE_ERF:                        {Sym: "erf", Description: `Various link-layer types, with a pseudo-header, for Endace DAG cards; encapsulates Endace ERF records.`},
	LINKTYPE_BLUETOOTH_HCI_H4_WITH_PHDR: {Sym: "bluetooth_hci_h4_with_phdr", Description: `Bluetooth HCI UART transport layer; the frame contains a 4-byte direction field, in network byte order (big-endian), the low-order bit of which is set if the frame was sent from the host to the controller and clear if the frame was received by the host from the controller, followed by an HCI packet indicator byte, as specified by the UART Transport Layer portion of the most recent Bluetooth Core specification, followed by an HCI packet of the specified packet type, as specified by the Host Controller Interface Functional Specification portion of the most recent Bluetooth Core Specification.`},
	LINKTYPE_AX25_KISS:                  {Sym: "ax25_kiss", Description: `AX.25 packet, with a 1-byte KISS header containing a type indicator.`},
	LINKTYPE_LAPD:                       {Sym: "lapd", Description: `Link Access Procedures on the D Channel (LAPD) frames, as specified by ITU-T Recommendation Q.920 and ITU-T Recommendation Q.921, starting with the address field, with no pseudo-header.`},
	LINKTYPE_PPP_WITH_DIR:               {Sym: "ppp_with_dir", Description: `PPP, as per RFC 1661 and RFC 1662, preceded with a one-byte pseudo-header with a zero value meaning "received by this host" and a non-zero value meaning "sent by this host"; if the first 2 bytes are 0xff and 0x03, it's PPP in HDLC-like framing, with the PPP header following those two bytes, otherwise it's PPP without framing, and the packet begins with the PPP header. The data in the frame is not octet-stuffed or bit-stuffed.`},
	LINKTYPE_C_HDLC_WITH_DIR:            {Sym: "c_hdlc_with_dir", Description: `Cisco PPP with HDLC framing, as per section 4.3.1 of RFC 1547, preceded with a one-byte pseudo-header with a zero value meaning "received by this host" and a non-zero value meaning "sent by this host".`},
	LINKTYPE_FRELAY_WITH_DIR:            {Sym: "frelay_with_dir", Description: `Frame Relay LAPF frames, beginning with a one-byte pseudo-header with a zero value meaning "received by this host" (DCE->DTE) and a non-zero value meaning "sent by this host" (DTE->DCE), followed by an ITU-T Recommendation Q.922 LAPF header starting with the address field, and without an FCS at the end of the frame.`},
	LINKTYPE_LAPB_WITH_DIR:              {Sym: "lapb_with_dir", Description: `Link Access Procedure, Balanced (LAPB), as specified by ITU-T Recommendation X.25, preceded with a one-byte pseudo-header with a zero value meaning "received by this host" (DCE->DTE) and a non-zero value meaning "sent by this host" (DTE->DCE).`},
	LINKTYPE_IPMB_LINUX:                 {Sym: "ipmb_linux", Description: `IPMB over an I2C circuit, with a Linux-specific pseudo-header.`},
	LINKTYPE_FLEXRAY:                    {Sym: "flexray", Description: `FlexRay automotive bus frames or symbols, preceded by a pseudo-header.`},
	LINKTYPE_LIN:                        {Sym: "lin", Description: `Local Interconnect Network (LIN) automotive bus, preceded by a pseudo-header.`},
	LINKTYPE_IEEE802_15_4_NONASK_PHY:    {Sym: "ieee802_15_4_nonask_phy", Description: `IEEE 802.15.4 Low-Rate Wireless Networks, with each packet having the FCS at the end of the frame, and with the PHY-level data for the O-QPSK, BPSK, GFSK, MSK, and RCC DSS BPSK PHYs (4 octets of 0 as preamble, one octet of SFD, one octet of frame length + reserved bit) preceding the MAC-layer data (starting with the frame control field).`},
	LINKTYPE_USB_LINUX_MMAPPED:          {Sym: "usb_linux_mmapped", Description: `USB packets, beginning with a Linux USB header, as specified by the struct usbmon_packet in the Documentation/usb/usbmon.txt file in the Linux source tree. All 64 bytes of the header are present. All fields in the header are in host byte order. When performing a live capture, the host byte order is the byte order of the machine on which the packets are captured. When reading a pcap file, the byte order is the byte order for the file, as specified by the file's magic number; when reading a pcapng file, the byte order is the byte order for the section of the pcapng file, as specified by the Section Header Block. For isochronous transfers, the ndesc field specifies the number of isochronous descriptors that follow.`},
	LINKTYPE_FC_2:                       {Sym: "fc_2", Description: `Fibre Channel FC-2 frames, beginning with a Frame_Header.`},
	LINKTYPE_FC_2_WITH_FRAME_DELIMS:     {Sym: "fc_2_with_frame_delims", Description: `Fibre Channel FC-2 frames, beginning an encoding of the SOF, followed by a Frame_Header, and ending with an encoding of the SOF.`},
	LINKTYPE_IPNET:                      {Sym: "ipnet", Description: `Solaris ipnet pseudo-header, followed by an IPv4 or IPv6 datagram.`},
	LINKTYPE_CAN_SOCKETCAN:              {Sym: "can_socketcan", Description: `CAN (Controller Area Network) frames, with a pseudo-header followed by the frame payload.`},
	LINKTYPE_IPV4:                       {Sym: "ipv4", Description: `Raw IPv4; the packet begins with an IPv4 header.`},
	LINKTYPE_IPV6:                       {Sym: "ipv6", Description: `Raw IPv6; the packet begins with an IPv6 header.`},
	LINKTYPE_IEEE802_15_4_NOFCS:         {Sym: "ieee802_15_4_nofcs", Description: `IEEE 802.15.4 Low-Rate Wireless Network, without the FCS at the end of the frame.`},
	LINKTYPE_DBUS:                       {Sym: "dbus", Description: `Raw D-Bus messages, starting with the endianness flag, followed by the message type, etc., but without the authentication handshake before the message sequence.`},
	LINKTYPE_DVB_CI:                     {Sym: "dvb_ci", Description: `DVB-CI (DVB Common Interface for communication between a PC Card module and a DVB receiver), with the message format specified by the PCAP format for DVB-CI specification.`},
	LINKTYPE_MUX27010:                   {Sym: "mux27010", Description: `Variant of 3GPP TS 27.010 multiplexing protocol (similar to, but not the same as, 27.010).`},
	LINKTYPE_STANAG_5066_D_PDU:          {Sym: "stanag_5066_d_pdu", Description: `D_PDUs as described by NATO standard STANAG 5066, starting with the synchronization sequence, and including both header and data CRCs. The current version of STANAG 5066 is backwards-compatible with the 1.0.2 version, although newer versions are classified.`},
	LINKTYPE_NFLOG:                      {Sym: "nflog", Description: `Linux netlink NETLINK NFLOG socket log messages.`},
	LINKTYPE_NETANALYZER:                {Sym: "netanalyzer", Description: `Pseudo-header for Hilscher Gesellschaft für Systemautomation mbH netANALYZER devices, followed by an Ethernet frame, beginning with the MAC header and ending with the FCS.`},
	LINKTYPE_NETANALYZER_TRANSPARENT:    {Sym: "netanalyzer_transparent", Description: `Pseudo-header for Hilscher Gesellschaft für Systemautomation mbH netANALYZER devices, followed by an Ethernet frame, beginning with the preamble, SFD, and MAC header, and ending with the FCS.`},
	LINKTYPE_IPOIB:                      {Sym: "ipoib", Description: `IP-over-InfiniBand, as specified by RFC 4391 section 6.`},
	LINKTYPE_MPEG_2_TS:                  {Sym: "mpeg_2_ts", Description: `MPEG-2 Transport Stream transport packets, as specified by ISO 13818-1/ITU-T Recommendation H.222.0 (see table 2-2 of section 2.4.3.2 "Transport Stream packet layer").`},
	LINKTYPE_NG40:                       {Sym: "ng40", Description: `Pseudo-header for ng4T GmbH's UMTS Iub/Iur-over-ATM and Iub/Iur-over-IP format as used by their ng40 protocol tester, followed by frames for the Frame Protocol as specified by 3GPP TS 25.427 for dedicated channels and 3GPP TS 25.435 for common/shared channels in the case of ATM AAL2 or UDP traffic, by SSCOP packets as specified by ITU-T Recommendation Q.2110 for ATM AAL5 traffic, and by NBAP packets for SCTP traffic.`},
	LINKTYPE_NFC_LLCP:                   {Sym: "nfc_llcp", Description: `Pseudo-header for NFC LLCP packet captures, followed by frame data for the LLCP Protocol as specified by NFCForum-TS-LLCP_1.1.`},
	LINKTYPE_INFINIBAND:                 {Sym: "infiniband", Description: `Raw InfiniBand frames, starting with the Local Routing Header, as specified in Chapter 5 "Data packet format" of InfiniBand™ Architectural Specification Release 1.2.1 Volume 1 - General Specifications.`},
	LINKTYPE_SCTP:                       {Sym: "sctp", Description: `SCTP packets, as defined by RFC 4960, with no lower-level protocols such as IPv4 or IPv6.`},
	LINKTYPE_USBPCAP:                    {Sym: "usbpcap", Description: `USB packets, beginning with a USBPcap header.`},
	LINKTYPE_RTAC_SERIAL:                {Sym: "rtac_serial", Description: `Serial-line packet header for the Schweitzer Engineering Laboratories "RTAC" product, followed by a payload for one of a number of industrial control protocols.`},
	LINKTYPE_BLUETOOTH_LE_LL:            {Sym: "bluetooth_le_ll", Description: `Bluetooth Low Energy air interface Link Layer packets, in the format described in section 2.1 "PACKET FORMAT" of volume 6 of the Bluetooth Specification Version 4.0 (see PDF page 2200), but without the Preamble.`},
	LINKTYPE_NETLINK:                    {Sym: "netlink", Description: `Linux Netlink capture encapsulation.`},
	LINKTYPE_BLUETOOTH_LINUX_MONITOR:    {Sym: "bluetooth_linux_monitor", Description: `Bluetooth Linux Monitor encapsulation of traffic for the BlueZ stack.`},
	LINKTYPE_BLUETOOTH_BREDR_BB:         {Sym: "bluetooth_bredr_bb", Description: `Bluetooth Basic Rate and Enhanced Data Rate baseband packets.`},
	LINKTYPE_BLUETOOTH_LE_LL_WITH_PHDR:  {Sym: "bluetooth_le_ll_with_phdr", Description: `Bluetooth Low Energy link-layer packets.`},
	LINKTYPE_PROFIBUS_DL:                {Sym: "profibus_dl", Description: `PROFIBUS data link layer packets, as specified by IEC standard 61158-4-3, beginning with the start delimiter, ending with the end delimiter, and including all octets between them.`},
	LINKTYPE_PKTAP:                      {Sym: "pktap", Description: `Apple PKTAP capture encapsulation.`},
	LINKTYPE_EPON:                       {Sym: "epon", Description: `Ethernet-over-passive-optical-network packets, starting with the last 6 octets of the modified preamble as specified by 65.1.3.2 "Transmit" in Clause 65 of Section 5 of IEEE 802.3, followed immediately by an Ethernet frame.`},
	LINKTYPE_IPMI_HPM_2:                 {Sym: "ipmi_hpm_2", Description: `IPMI trace packets, as specified by Table 3-20 "Trace Data Block Format" in the PICMG HPM.2 specification. The time stamps for packets in this format must match the time stamps in the Trace Data Blocks.`},
	LINKTYPE_ZWAVE_R1_R2:                {Sym: "zwave_r1_r2", Description: `Z-Wave RF profile R1 and R2 packets, as specified by ITU-T Recommendation G.9959, with some MAC layer fields moved.`},
	LINKTYPE_ZWAVE_R3:                   {Sym: "zwave_r3", Description: `Z-Wave RF profile R3 packets, as specified by ITU-T Recommendation G.9959, with some MAC layer fields moved.`},
	LINKTYPE_WATTSTOPPER_DLM:            {Sym: "wattstopper_dlm", Description: `Formats for WattStopper Digital Lighting Management (DLM) and Legrand Nitoo Open protocol common packet structure captures.`},
	LINKTYPE_ISO_14443:                  {Sym: "iso_14443", Description: `Messages between ISO 14443 contactless smartcards (Proximity Integrated Circuit Card, PICC) and card readers (Proximity Coupling Device, PCD), with the message format specified by the PCAP format for ISO14443 specification.`},
	LINKTYPE_RDS:                        {Sym: "rds", Description: `Radio data system (RDS) groups, as per IEC 62106, encapsulated in this form.`},
	LINKTYPE_USB_DARWIN:                 {Sym: "usb_darwin", Description: `USB packets, beginning with a Darwin (macOS, etc.) USB header.`},
	LINKTYPE_SDLC:                       {Sym: "sdlc", Description: `SDLC packets, as specified by Chapter 1, "DLC Links", section "Synchronous Data Link Control (SDLC)" of Systems Network Architecture Formats, GA27-3136-20, without the flag fields, zero-bit insertion, or Frame Check Sequence field, containing SNA path information units (PIUs) as the payload.`},
	LINKTYPE_LORATAP:                    {Sym: "loratap", Description: `LoRaTap pseudo-header, followed by the payload, which is typically the PHYPayload from the LoRaWan specification.`},
	LINKTYPE_VSOCK:                      {Sym: "vsock", Description: `Protocol for communication between host and guest machines in VMware and KVM hypervisors.`},
	LINKTYPE_NORDIC_BLE:                 {Sym: "nordic_ble", Description: `Messages to and from a Nordic Semiconductor nRF Sniffer for Bluetooth LE packets, beginning with a pseudo-header.`},
	LINKTYPE_DOCSIS31_XRA31:             {Sym: "docsis31_xra31", Description: `DOCSIS packets and bursts, preceded by a pseudo-header giving metadata about the packet.`},
	LINKTYPE_ETHERNET_MPACKET:           {Sym: "ethernet_mpacket", Description: `mPackets, as specified by IEEE 802.3br Figure 99-4, starting with the preamble and always ending with a CRC field.`},
	LINKTYPE_DISPLAYPORT_AUX:            {Sym: "displayport_aux", Description: `DisplayPort AUX channel monitoring data as specified by VESA DisplayPort(DP) Standard preceded by a pseudo-header.`},
	LINKTYPE_LINUX_SLL2:                 {Sym: "linux_sll2", Description: `Linux "cooked" capture encapsulation v2.`},
	LINKTYPE_OPENVIZSLA:                 {Sym: "openvizsla", Description: `Openvizsla FPGA-based USB sniffer.`},
	LINKTYPE_EBHSCR:                     {Sym: "ebhscr", Description: `Elektrobit High Speed Capture and Replay (EBHSCR) format.`},
	LINKTYPE_VPP_DISPATCH:               {Sym: "vpp_dispatch", Description: `Records in traces from the http://fd.io VPP graph dispatch tracer, in the the graph dispatcher trace format.`},
	LINKTYPE_DSA_TAG_BRCM:               {Sym: "dsa_tag_brcm", Description: `Ethernet frames, with a switch tag inserted between the source address field and the type/length field in the Ethernet header.`},
	LINKTYPE_DSA_TAG_BRCM_PREPEND:       {Sym: "dsa_tag_brcm_prepend", Description: `Ethernet frames, with a switch tag inserted before the destination address in the Ethernet header.`},
	LINKTYPE_IEEE802_15_4_TAP:           {Sym: "ieee802_15_4_tap", Description: `IEEE 802.15.4 Low-Rate Wireless Networks, with a pseudo-header containing TLVs with metadata preceding the 802.15.4 header.`},
	LINKTYPE_DSA_TAG_DSA:                {Sym: "dsa_tag_dsa", Description: `Ethernet frames, with a switch tag inserted between the source address field and the type/length field in the Ethernet header.`},
	LINKTYPE_DSA_TAG_EDSA:               {Sym: "dsa_tag_edsa", Description: `Ethernet frames, with a programmable Ethernet type switch tag inserted between the source address field and the type/length field in the Ethernet header.`},
	LINKTYPE_ELEE:                       {Sym: "elee", Description: `Payload of lawful intercept packets using the ELEE protocol. The packet begins with the ELEE header; it does not include any transport-layer or lower-layer headers for protcols used to transport ELEE packets.`},
	LINKTYPE_Z_WAVE_SERIAL:              {Sym: "z_wave_serial", Description: `Serial frames transmitted between a host and a Z-Wave chip over an RS-232 or USB serial connection, as described in section 5 of the Z-Wave Serial API Host Application Programming Guide.`},
	LINKTYPE_USB_2_0:                    {Sym: "usb_2_0", Description: `USB 2.0, 1.1, or 1.0 packet, beginning with a PID, as described by Chapter 8 "Protocol Layer" of the the Universal Serial Bus Specification Revision 2.0.`},
	LINKTYPE_ATSC_ALP:                   {Sym: "atsc_alp", Description: `ATSC Link-Layer Protocol frames, as described in section 5 of the A/330 Link-Layer Protocol specification, found at the ATSC 3.0 standards page, beginning with a Base Header.`},
	LINKTYPE_ETW:                        {Sym: "etw", Description: `Event Tracing for Windows messages, beginning with a pseudo-header.`},
}

var linkToFormat = map[int]*decode.Group{
	LINKTYPE_ETHERNET: &pcapngEther8023Format,
}
