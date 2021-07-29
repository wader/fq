# <mp4 value> | mp4_path(".moov.trak[1]") -> box
# <mp4 value> | mp4_path(<mp4 value>) -> ".moov.trak"
def mp4_path(p): tree_path(.boxes; .type; p);