# <matroska root value> | matroska_path(".Segment.Tracks[0]") -> element
# <matroska root value> | matroska_path -> ".Segment.Tracks[0]"
# <matroska root value> | matroska_path(<matroska root value>) -> ".Segment.Tracks[0]"
def matroska_path(p):
    _decode_value(
        ( if format != "matroska" then error("not matroska format") end
        | _tree_path(.elements; .id; p)
        )
    );
def matroska_path:
    ( . as $c
    | format_root
    | matroska_path($c)
    );

def _matroska__help:
  { examples: [
      {comment: "Lookup element decode value using `matroska_path`", expr: "matroska_path(\".Segment.Tracks[0)\""},
      {comment: "Return `matroska_path` string for a box decode value", expr: "grep_by(.id == \"Tracks\") | matroska_path"}
    ],
    links: [
      {url: "https://tools.ietf.org/html/draft-ietf-cellar-ebml-00"},
      {url: "https://matroska.org/technical/specs/index.html"},
      {url: "https://www.matroska.org/technical/basics.html"},
      {url: "https://www.matroska.org/technical/codec_specs.html"},
      {url: "https://wiki.xiph.org/MatroskaOpus"}
    ]
  };