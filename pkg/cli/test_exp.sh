#!/bin/sh
set -eu

FQ="$1"
shift

if command -v expect >/dev/null 2>&1; then
    TEMPDIR=$(mktemp -d)
    cp "$FQ" "${TEMPDIR}/fq"
    PATH="${TEMPDIR}:${PATH}" expect "$1" >"${TEMPDIR}/fq.log" && FAIL=0 || FAIL=1
    if [ $FAIL = "1" ]; then
        cat "${TEMPDIR}/fq.log"
    fi
    rm -rf "${TEMPDIR}"
    if [ $FAIL = "1" ]; then
        exit 1
    fi
    echo "$0"
else
    echo "$0: skip as expect is not installed"
fi
