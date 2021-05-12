#!/usr/bin/env fq -rf
# find and parse a x264 SEI payload
# looks like this:
# x264 - core 161 r3020 d198931 - H.264/MPEG-4 AVC codec - Copyleft 2003-2020 - http://www.videolan.org/x264.html - options: cabac=1 ref=3 deblock=1:0:0 analyse=0x3:0x113 me=hex subme=7 psy=1 psy_rd=1.00:0.00 mixed_ref=1 me_range=16 chroma_me=1 trellis=1 8x8dct=1 cqm=0 deadzone=21,11 fast_pskip=1 chroma_qp_offset=4 threads=6 lookahead_threads=1 sliced_threads=0 nr=0 decimate=1 interlaced=0 bluray_compat=0 constrained_intra=0 bframes=3 b_pyramid=2 b_adapt=1 b_bias=0 direct=1 weightb=1 open_gop=0 weightp=2 keyint=250 keyint_min=25 scenecut=40 intra_refresh=0 rc_lookahead=40 rc=crf mbtree=1 crf=23.0 qcomp=0.60 qpmin=0 qpmax=69 qpstep=4 ip_ratio=1.40 aq=1:1.00

( ..
  | select(._format == "avc_sei" and .uuid._symbol == "x264")
  | .data
  | string[0:-1]
  | . as $full
  | split("options: ")[1]
  | ( [
      .
      | split(" ")[]
      | split("=")
      | {key: .[0], value: .[1]}
      ] | from_entries
  ) + {
      full: $full
  }
) // "No x264 SEI found"
