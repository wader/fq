# go run . -L format/tls/testdata -o keylog=@format/tls/testdata/dump.pcapng.keylog -f format/tls/testdata/split.jq  format/tls/testdata/dump.pcapng | tar -C format/tls/testdata/ciphers -x

include "to_tar";

def ipv4_tcp_tuple:
  ( . as {source_ip: $sip, destination_ip: $dip}
  | grep_by(format=="tcp_segment") as {source_port: $sport, destination_port: $dport}
  | [[$sip,$sport],[$dip,$dport]]
  | sort
  );

def connetions_tuples:
  ( [ grep_by(format=="ipv4_packet")
    | ipv4_tcp_tuple
    ]
  | unique[]
  );

def to_ipv4_pcap(packets):
  # TODO: hack
  [ ("d4c3b2a1020004000000000000000000ffff0000e4000000" | from_hex)
  , ( packets | tobytes | [band(.size;0xff),band(bsr(.size;8);0xff),0,0] as $sz
  | [0,0,0,0,0,0,0,0,$sz,$sz,.])
  ] | tobytes;

( .[0].blocks
| . as $packets
| to_tar(
    ( $packets
    | connetions_tuples as $tuple
    | to_ipv4_pcap(
        ( $packets
        | grep_by(format=="ipv4_packet")
        | select($tuple ==ipv4_tcp_tuple)
        )
      )
    | . as $pcap_bytes
    | pcap
    | .tcp_connections[0].server.stream.records[0].message as {$cipher_suit}
    | {filename: "\($cipher_suit).pcap", data: $pcap_bytes}
    )
  )
)

