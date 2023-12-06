#!/usr/bin/env bash
set -Eeuo pipefail
d="$( cd "$( dirname "$0" )" && pwd -P )"

RESET=$'\e[0m'
BOLD=$'\e[1m'
GREEN=$'\e[0;32m'
RED=$'\e[0;31m'

tmpdir=$(mktemp -d -t fq-wasm-generate-fqtest.XXXXXXXX)
trap 'tear_down' 0

result=1
tear_down() {
    : "Clean up tmpdir" && {
        [[ $tmpdir ]] && rm -rf "$tmpdir"
    }

    : "Report result" && {
        if [ "$result" -eq 0 ]; then
            echo
            echo -e "${GREEN}${BOLD}OK${RESET}"
            echo
        else
            echo
            echo -e "${RED}${BOLD}FAILED${RESET}"
            echo
        fi
        exit $result
    }
}

(
  cd "$d"
  rm -f ./*.fqtest

  find "." -name '*.wasm' -print0 | sort -z -n |
    while IFS= read -r -d '' f; do
      f=${f#./}
      of=${f/%.wasm/.fqtest}
      echo "$f -> $of"
      echo "$ fq -d wasm dv $f" > "$of"
    done

  go test -run TestFQTests/wasm "$d/../../.." -update
  if grep "error" ./*.fqtest; then
    exit 1
  fi
)

result=0
