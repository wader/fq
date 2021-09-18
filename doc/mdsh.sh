#!/bin/sh
set -eu

REPODIR=$(pwd)
TEMPDIR=$(mktemp -d)
cp -a doc/* "${TEMPDIR}"
go build -o "${TEMPDIR}/fq" main.go
for f in "$@"; do
    cd "${TEMPDIR}"
    echo "Generate $f"
    mkdir -p "$(dirname "${TEMPDIR}/$f")"
    PATH="${TEMPDIR}:${PATH}" go run "${REPODIR}/doc/mdsh.go" >"${TEMPDIR}/$f" <"${REPODIR}/$f"
    mv "${TEMPDIR}/$f" "${REPODIR}/$f"
done
rm -rf "${TEMPDIR}"
