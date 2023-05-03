# note this is a "dynamic" include, output string will be used as source

# generates decode functions, ex:
# mp3/0 calls decode("mp3"; {})
# mp3/1 calls decode("mp3"; $opts)
# from_mp3/* same but throws error on decode error

[ _registry as $r
| $r.groups
| to_entries[]
# skip_decode_function is used to skip bits/bytes as they are special tobits/tobytes
| select($r.formats[.key].skip_decode_function | not)
| "def \(.key)($opts): decode(\(.key | tojson); $opts);"
, "def \(.key): decode(\(.key | tojson); {});"
, "def from_\(.key)($opts): decode(\(.key | tojson); $opts) | if ._error then error(._error.error) end;"
, "def from_\(.key): from_\(.key)({});"
] | join("\n")
