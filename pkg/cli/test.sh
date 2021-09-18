#!/bin/sh
set -eu

if which expect >/dev/null 2>&1; then
    TEMPDIR=$(mktemp -d)
    go build -o "${TEMPDIR}/fq" main.go
    PATH="${TEMPDIR}:${PATH}" expect "$1" >"${TEMPDIR}/fq.log" && FAIL=0 || FAIL=1
    if [ $FAIL = "1" ]; then
        cat "${TEMPDIR}/fq.log"
    fi
    rm -rf "${TEMPDIR}"
    if [ $FAIL = "1" ]; then
        exit 1
    fi
else
    echo "skip as expect is not installed"
fi
