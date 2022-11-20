# note this is a "dynamic" include, output string will be used as source

# generates decode functions
# frommp3 and mp3 calls decode("mp3")

[ _registry as $r
| $r.groups
| to_entries[]
# TODO: nicer way to skip "all" which also would override builtin all/*
| select(.key != "all" and ($r.formats[.key].skip_decode_function | not))
| "def \(.key)($opts): decode(\(.key | tojson); $opts);"
, "def \(.key): decode(\(.key | tojson); {});"
, "def from\(.key)($opts): decode(\(.key | tojson); $opts) | if ._error then error(._error.error) end;"
, "def from\(.key): from\(.key)({});"
] | join("\n")
