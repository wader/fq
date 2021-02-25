#!/usr/bin/env fq -nf
def rpad($s;$w): . + ($s * ($w+1-length))[1:];
((formats | keys | map(length) | max)+1) as $m |
"|\("Name" | rpad(" ";$m))|Description|",
"|-|-|",
(
    formats | to_entries[] |
    "|\("`"+.key+"`"|rpad(" ";$m))|\(.value.description)|"
)
