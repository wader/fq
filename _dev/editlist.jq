#!/usr/bin/env fq -s

# TODO: esds, make fancy printer? shared?
# TODO: handle -1 media_time

open($FILENAME) | mp4 |
(.boxes[] | select(.type == "moov").boxes[] | select(.type == "mvhd")) as $mvhd |
{
    time_scale: $mvhd.time_scale,
    tracks: [
        .. | select(.type == "trak") |
        (.. | select(.type == "mdhd")) as $mdhd |
        (.. | select(.type == "hdlr")) as $hdlr |
        (.. | select(.type == "stsd")) as $stsd |
        (.. | select(.type == "elst")) as $elst |
        {
            component_type: $hdlr.component_subtype,
            data_format: $stsd.sample_descriptions[0].data_format,
            track: {
                media_scale: $mdhd.time_scale,
                edit_list: [
                    $elst.table[] | {
                        time_scale: $mdhd.time_scale,
                        track_duration: .track_duration,
                        media_time: .media_time,
                        movie_track_duration: (.track_duration / $mvhd.time_scale),
                        movie_media_time: (.media_time / $mdhd.time_scale)
                    }
                ]
            }
        }
    ]
}
