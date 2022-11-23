#!/bin/sh

# 8 16 24 32 40 48 56 64 80 96 112 128 144 160 256 320
for br in 8000 128000 320000; do
    for ch in 1 2; do
        # 48000 44100 32000 22050 24000 16000 11025 12000 8000
        for hz in 48000 44100 8000; do
            f="mp3-frame-${br}br-${ch}ch-${hz}hz"
            ffmpeg -y -f lavfi -i sine -ar $hz -b:a $br -ac $ch -t 1s -id3v2_version 0 -write_xing 0 -f mp3 $f.temp
            fq -d bytes '[limit(3; match([0xff,0xfb],[0xff,0xe3]; "g").offset)] as $o | .[$o[1]:$o[2]]' $f.temp >$f
            rm $f.temp
            echo "\$ fq -d mp3_frame dv $f" >$f.fqtest
        done
    done
done
