#!/usr/bin/env fq -rnf
[formats[] | "\(.name)"] | join(",\n")
