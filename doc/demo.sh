#!/bin/bash

FQ="$1"

c() {
    echo -e "\x1b[97m# $1\x1b[0m"
}

s() {
    echo "\$ $1"
    sh -c "${1/fq/$FQ -o unicode=true -C}"
}

c "Overview of mp3 file"
s "fq . file.mp3"
echo
c "Show ID3v2 tag inside mp3 file"
s "fq '.headers[0]' file.mp3"
echo
c "Resolution of embedded PNG cover art"
s "fq '.headers[0].frames[] | select(.id == \"APIC\").picture.chunks[] | select(.type == \"IHDR\") | {width, height}' file.mp3"
echo
c "Extract PNG file"
s "fq '.headers[].frames[] | select(.id == \"APIC\")?.picture | tobits' file.mp3 >file.png"
s "file file.png"
rm -f file.png
