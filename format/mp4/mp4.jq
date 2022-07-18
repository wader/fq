# <mp4 root> | mp4_path(".moov.trak[1]") -> box
# box -> | mp4_path -> ".moov.trak[1]"
# box -> | mp4_path(<mp4 root>) -> ".moov.trak[1]"
def mp4_path(p):
  _decode_value(
    ( if format != "mp4" then error("not mp4 format") end
    | _tree_path(.boxes; .type; p)
    )
  );
def mp4_path:
  ( . as $c
  | format_root
  | mp4_path($c)
  );

def _mp4__help:
  { notes: "Support `mp4_path`",
    examples: [
      {comment: "Lookup box decode value using `mp4_path`", expr: "mp4_path(\".moov.trak[1]\")"},
      {comment: "Return `mp4_path` string for a box decode value", expr: "grep_by(.type == \"trak\") | mp4_path"}
    ],
    links: [
      {title: "ISO/IEC base media file format (MPEG-4 Part 12)", url: "https://en.wikipedia.org/wiki/ISO/IEC_base_media_file_format"},
      {title: "Quicktime file format", url: "https://developer.apple.com/standards/qtff-2001.pdf"}
    ]
  };
