# <matroska root value> | matroska_path(".segment.tracks[0]") -> element
# <matroska root value> | matroska_path -> ".segment.tracks[0]"
# <matroska root value> | matroska_path(<matroska root value>) -> ".segment.tracks[0]"
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
