#!/bin/sh

if which expect >/dev/null 2>&1; then
    TEMPDIR=$(mktemp -d)
    go build -o "${TEMPDIR}/fq" main.go
    PATH="${TEMPDIR}:${PATH}" expect dev/fqbin-test.exp >"${TEMPDIR}/fq.log"
    EXIT="$?"
    if [ $EXIT != "0" ]; then
        cat "${TEMPDIR}/fq.log"
    fi
    rm -rf "${TEMPDIR}"
    if [ $EXIT != "0" ]; then
        exit 1
    fi
else
    echo "fq-test.sh: skip as expect is not installed"
fi
