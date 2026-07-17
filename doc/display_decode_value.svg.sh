#!/bin/bash

FQ="$1"

s() {
    echo "\$ $1"
    sh -c "${1/fq/$FQ}"
}

s "fq 'first(grep_by(.type == \"elst\"))' file.mp4"
