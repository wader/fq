#!/usr/bin/env fq -rf

# xor input array with key array
# key will be repeated to fit length of input
def xor_array($key):
  # [1,2,3] | repeat(7) -> [1,2,3,1,2,3,1]
  def repeat($len):
    ( length as $l
    | [.[range($len) % $l]]
    );
  ( . as $input
  # [$input, $key repeated]
  | [ $input
    , ($key | repeat($input | length))
    ]
  # [[$input[0], $key[0], ...]
  | transpose
  | map(bxor(.[0]; .[1]))
  );

( first(.uncompressed.files[] | select(.name == "triangle.pcap"))
| .data.tcp_connections[0].server.stream
| fromjsonl
| map(
    ( if .encoding == "xor" then
        .data |= (fromhex | explode | xor_array([71])| implode)
      end
    | if .encoding == "long_xor" then
        ( .data |=
          ( fromhex
          | explode
          | xor_array("GravityForce" | explode)
          | implode
          )
        )
      end
    | if .msg_type == "update" then .data |= fromjson end
    )
  ) as $msgs
| { "svg": {
      # move viewbox to where the objects are
      "@viewBox": "50 120 350 350",
      "@width": 350,
      "@height": 350,
      "@xmlns": "http://www.w3.org/2000/svg",
      "rect": [
        { "#seq": -1,
          "@fill": "#101010",
          "@x": 50,
          "@y": 120,
          "@width": 350,
          "@height": 350
        }
        # gather last update for all objects and draw them
        , ( reduce ($msgs[] | select(.msg_type == "update")) as $msg (
              {};
              # use tostring as object keys can only be strings
              .[$msg.data.id | tostring] = $msg.data
            )
          | { Player: {style: "fill: #0000ff", size: 15},
              Flower: {style: "fill: #00d000", size: 10},
              Rock: {style: "fill: #a0a0a0", size: 8}
            } as $types
          | .[]
          | $types[.type] as $t
          | { "@width": $t.size,
              "@height": $t.size,
              "@style": $t.style,
              "@transform": "rotate(\(.rot) \(.x-$t.size/2) \(.y-$t.size/2))",
              "@x": (.x-$t.size/2),
              "@y": (.y-$t.size/2)
            }
          )
      ],
      "polyline": {
        "#seq": 1,
        "@fill": "none",
        "@stroke": "#5050d0",
        "@stroke-dasharray": "5 10",
        "@points":
          ( [ $msgs[]
            | select(.msg_type == "update" and .data.type == "Player")
            | .data.x, .data.y
            ]
          | join(" ")
          )
      }
    }
  }
| toxml({indent: 2})
)
