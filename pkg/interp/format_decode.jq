# note this is a "dynamic" include, output string will be used as source

[ _registry.groups
| to_entries[]
# TODO: nicer way to skip "all" which also would override builtin all/*
| select(.key != "all")
| "def \(.key)($opts): decode(\(.key | tojson); $opts);"
, "def \(.key): decode(\(.key | tojson); {});"
, "def from\(.key)($opts): decode(\(.key | tojson); $opts) | if ._error then error(._error.error) end;"
, "def from\(.key): from\(.key)({});"
] | join("\n")
