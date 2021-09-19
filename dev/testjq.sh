#!/bin/sh
set -eu

FQ="$1"
shift

for f in "$@"; do
    echo "$f"
    "$FQ" -nr -L "$(dirname "$f")" -f "$f"
done
