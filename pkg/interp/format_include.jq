# note this is a "dynamic" include, output string will be used as source

[ _registry.formats[]
| select(.files)
| .files[]
| select(.name | endswith(".jq"))
| .data
] | join("\n")
