#!/bin/sh
set -eu

FQ="$1"
shift

REPODIR=$(pwd)
export REPODIR
TEMPDIR=$(mktemp -d)
cp -a doc/* "${TEMPDIR}"
cp "$FQ" "${TEMPDIR}/fq"
cd "${TEMPDIR}"
for f in "$@"; do
    echo "Generate $f"
    mkdir -p "$(dirname "${TEMPDIR}/$f")"
    PATH="${TEMPDIR}:${PATH}" go run "${REPODIR}/doc/mdsh.go" >"${TEMPDIR}/$f" <"${REPODIR}/$f"
    mv "${TEMPDIR}/$f" "${REPODIR}/$f"
done
rm -rf "${TEMPDIR}"
