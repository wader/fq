# used to build functions like mp4_path
def tree_path(children; name; $v):
	def _lookup:
		# add implicit zeros to get first value
		# ["a", "b", 1] => ["a", 0, "b", 1]
		def _normalize_path:
			. as $np
			| if $np | last | type == "string" then $np+[0] end
			| (reduce .[] as $np ([[], []];
				if $np | type == "string" then
					[(.[0]+.[1]+[$np]), [0]]
				else
					[.[0]+[$np], []]
				end
			))[0];
		. as $c
		| $v
		| expr_to_path
		| _normalize_path
		| reduce .[] as $n ($c;
			if $n | type == "string" then
				children | map(select(name==$n))
			else
				.[$n]
			end
		);
	def _path:
		[
		. as $r
		| $v._path as $p
		| foreach range(($p | length)/2) as $i (
			null;
			null;
			($r | getpath($p[0:($i+1)*2]) | name) as $name
			| [($r | getpath($p[0:($i+1)*2-1]))[] | name][0:$p[($i*2)+1]+1] as $before
			| [
				$name,
				($before | map(select(. == $name)) | length)-1
			]
		)
		| [ ".", .[0],
			(.[1] | if . == 0 then empty else "[", ., "]" end)
		]
		]
		| flatten
		| join("");
	if $v | type == "string" then _lookup
	else _path end;

# <mp4 value> | mp4_path(".moov.trak[1]") -> box
# <mp4 value> | mp4_path(<mp4 value>) -> ".moov.trak"
def mp4_path(p): tree_path(.boxes; .type; p);

# <matroska value> | matroska_path(".Segment.Tracks[0].TrackEntry[1].CodecID") -> element
# <matroska value> | matroska_path(<matroska value>) -> ".Segment.Tracks[0]"
def matroska_path(p): tree_path(.elements; .id._symbol; p);
