# note this is a "dynamic" include, output string will be used as source

[ _registry.groups
| to_entries[]
# TODO: nicer way to skip "all" which also would override builtin all/*
| select(.key != "all")
| "def \(.key)($opts): decode(\(.key | tojson); $opts);"
, "def \(.key): decode(\(.key | tojson); {});"
] | join("\n")
