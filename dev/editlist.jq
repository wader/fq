#!/usr/bin/env fq -d mp4 -f

# TODO: esds, make fancy printer? shared?
# TODO: handle -1 media_time
# TODO: fragmented mp4

# root
#   moov
#     mvhd (movie header)
#     trak (track)
#       mdia
#         mdhd (media header)
#         hdlr (handler?)
#         minf
#           stbl
#             stsd (sample description)
#       elst (edit list)

( first(.boxes[] | select(.type == "moov")?)
| first(.boxes[] | select(.type == "mvhd")?) as $mvhd
| {
    duration: $mvhd.duration,
    time_scale: $mvhd.time_scale,
    duration_s: ($mvhd.duration / $mvhd.time_scale),
    tracks:
      [ .boxes[]
      | select(.type == "trak")
      | first(.. | select(.type == "mdhd")?) as $mdhd
      | first(.. | select(.type == "hdlr")?) as $hdlr
      | first(.. | select(.type == "stsd")?) as $stsd
      | first(.. | select(.type == "elst")?) as $elst
      | first(.. | select(.type == "stts")?) as $stts
      | ([$stts.entries[] | .count * .delta] | add) as $stts_sum
      | {
          component_type: $hdlr.component_subtype,
          # the sample descriptors are handled as boxes by the mp4 decoder
          data_format: $stsd.boxes[0].type,
          media_scale: $mdhd.time_scale,
          edit_list:
            [ $elst.entries[]
            | {
                track_duration: .segment_duration,
                media_time: .media_time,
                track_duration_s: (.segment_duration / $mvhd.time_scale),
                media_time_s: (.media_time / $mdhd.time_scale)
              }
            ],
          stts: {
            sum: $stts_sum,
            sum_s: ($stts_sum / $mdhd.time_scale)
          }
        }
      ]
  }
)
