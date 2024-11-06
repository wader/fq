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
c "Show header of first ID3v2 tag inside mp3 file"
s "fq '.headers[0].header' file.mp3"
echo
c "Show encoder software used"
s "fq -r '.frames[0].tag.encoder | tovalue' file.mp3"
echo
c "Decode at two offsets as mp3_frame and show bitrate"
s "fq -d bytes '.[0xb79,0xc49:] | mp3_frame.header.bitrate' file.mp3"
echo
c "Extract PNG file"
s "fq '.headers[].frames[] | select(.id == \"APIC\")?.picture | tobits' file.mp3 >file.png"
s "file file.png"
rm -f file.png
echo
c "Grep for PNG header, extract resolution and output as YAML"
s "fq -r 'grep_by(.type == \"IHDR\") | {res: {width, height}} | to_yaml' file.mp3"
#echo
c "Add query parameter to URL"
s "echo 'http://host?a=b' | fq -Rr 'from_url | .query.b = \"a b c\" | to_url'"
echo
c "Extract JSON and base64 encoded query parameter p"
s "echo 'https://host?p=eyJhIjoiaGVsbG8ifQ%3D%3D' | fq -R 'from_url.query.p | from_base64 | fromjson'"
