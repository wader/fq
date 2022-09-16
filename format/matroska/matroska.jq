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
