#!/usr/bin/env fq -f

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

open($FILENAME) | mp4 |
first(.boxes[] | select(.type == "moov")) |
first(.boxes[] | select(.type == "mvhd")) as $mvhd |
{
    duration: $mvhd.duration,
    time_scale: $mvhd.time_scale,
    duration_s: ($mvhd.duration / $mvhd.time_scale),
    tracks: [
        .boxes[] | select(.type == "trak") |
            first(.. | select(.type == "mdhd")) as $mdhd |
            first(.. | select(.type == "hdlr")) as $hdlr |
            first(.. | select(.type == "stsd")) as $stsd |
            first(.. | select(.type == "elst")) as $elst |
            {
                component_type: $hdlr.component_subtype,
                data_format: $stsd.sample_descriptions[0].data_format,
                media_scale: $mdhd.time_scale,
                edit_list: [
                    $elst.table[] | {
                        time_scale: $mdhd.time_scale,
                        track_duration: .track_duration,
                        media_time: .media_time,
                        track_duration_s: (.track_duration / $mvhd.time_scale),
                        media_time_s: (.media_time / $mdhd.time_scale)
                    }
                ]
            }
    ]
}
