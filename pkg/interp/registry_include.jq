[ _registry.files[][]
| select(.name | endswith(".jq"))
| .data
] | join("\n")
