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
