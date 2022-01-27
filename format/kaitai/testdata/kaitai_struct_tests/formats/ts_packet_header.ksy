# Transport stream packet
# see H.222.0 F.1.1

# |7,6,5,4,2,1,0|7,6,5,4,3,2,1,0|7,6,5,4,3,2,1,0|7,6,5,4,3,2,1,0|
# | sync_byte   | , , ,<-------pid------------->|   ,afc,<-cc-->|

meta:
  id: ts_packet_header
  endian: le
doc: >
  describes the first 4 header bytes of a TS Packet header
seq:
  - id: sync_byte
    type: u1
    #contents: [0x47]
  - id: transport_error_indicator
    type: b1
  - id: payload_unit_start_indicator
    type: b1
  - id: transport_priority
    type: b1
  - id: pid
    type: b13
  - id: transport_scrambling_control
    type: b2
  - id: adaptation_field_control
    type: b2
    enum: adaptation_field_control_enum
  - id: continuity_counter
    type: b4
  - id: ts_packet_remain
    size: 184
enums:
  adaptation_field_control_enum:
    0x0: reserved
    0x1: payload_only
    0x2: adaptation_field_only
    0x3: adaptation_field_and_payload
