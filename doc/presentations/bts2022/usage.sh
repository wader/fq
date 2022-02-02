#!/bin/bash

c() {
    echo -e "\x1b[97m# $1\x1b[0m"
}

s() {
    echo "\$ $1"
    sh -c "${1/fq/fq -o line_bytes=16 -o unicode=true -C}"
}

n() {
    sh -c "${1/fq/fq -o line_bytes=16 -o unicode=true -C}"
}

c "display a decode value"
s "fq . file.mp3"
echo
echo

c "expression returning a number"
s "fq '.frames | length' file.mp3"
echo
echo

c "raw bytes"
s "fq 'grep_by(format == \"png\") | tobytes' file.mp3 >file.png"
s "file file.png"
echo
echo

c "interactve REPL"
echo "\$ fq -i . file.mp3"
echo "mp3> .frames | length"
n "fq '.frames | length' file.mp3"
echo "mp3> .header[0] | repl"
echo "> .headers[0] id3v2> .frames[0].text"
n "fq '.headers[0].frames[0].text' file.mp3"
echo "> .headers[0] id3v2> .frames[0].text | tovalue"
n "fq '.headers[0].frames[0].text | tovalue' file.mp3"
echo "> .headers[0] id3v2> ^D"
echo "mp3> ^D"
echo "$"
