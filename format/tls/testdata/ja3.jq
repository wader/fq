# calculate ja3 client fingerprint
# https://github.com/salesforce/ja3
# TLSVersion,Ciphers,Extensions,EllipticCurves,EllipticCurvePointFormats
# ex:
# 769,47–53–5–10–49161–49162–49171–49172–50–56–19–4,0–10–11,23–24–25,0

# ja3 string
def to_ja3:
  def grease_values:
    [ 2570,6682,10794,14906,19018,23130,27242,31354
    , 35466,39578,43690,47802,51914,56026,60138,64250
    ];
  ( .records[0].message
  | [ [.version | toactual]
    , (.cipher_suits | map(toactual) - grease_values)
    , ([.extensions[]?.type | toactual] - grease_values)
    , ( [.extensions[]? | select(.type=="supported_groups").supported_groups[] | toactual]
      | . - grease_values
      )
    , [.extensions[]? | select(.type=="ec_point_formats").ex_points_formats[]]
    ]
  | map(join("-"))
  | join(",")
  );

# ja3 md5 hex digest
def to_ja3_digest: to_ja3 | to_md5 | to_hex;

# list ja3 string and digest in pcap or pcapng
def pcap_ja3:
  [ ( ( if format == "pcap" then .
        elif format == "pcapng" then .[]
        else error("not a pcap or pcapng decode value")
        end
      ).tcp_connections[]
    | . as {$client,$server}
    | .client.stream
    | select(format=="tls")
    | to_ja3? as $ja3
    | { client_ip: $client.ip
      , client_port: ($client.port | toactual)
      , server_ip: $server.ip
      , server_port: ($server.port | toactual)
      , ja3: $ja3
      , ja3_digest: ($ja3 | to_md5 | to_hex)
      }
    )
  ];
