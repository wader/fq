# <matroska value> | matroska_path(".Segment.Tracks[0].TrackEntry[1].CodecID") -> element
# <matroska value> | matroska_path(<matroska value>) -> ".Segment.Tracks[0]"
def matroska_path(p): tree_path(.elements; .id; p);