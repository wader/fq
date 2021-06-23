# same as group_by but just counts, grouping exp value will be converted to string
def count_by(exp):
	group_by(exp) | map({(.[0] | exp | tostring): length}) | add;

def protobuf_to_value:
	.fields | map({(.name|tostring): (.enum // .value)}) | add;

# hack to parse just a box
# <binary> | mp4_box
def mp4_box:
	[0,0,0,16, "ftyp", "isom", 0, 0 , 2 ,0, .] | mp4.boxes;

# TODO: a bit hacky?
def expr_to_path:
	if . | type != "string" then error("require string argument") end
	| eval("null | path(\(.))");

# TODO: generalize?
def array_tree_path(children; name; p):
	# add implicit zeros to get first value
	# ["a", "b", 1] => ["a", 0, "b", 1]
	def _normalize_path:
		. as $p
		| if $p | last | type == "string" then $p+[0] end
		| (reduce .[] as $p ([[], []];
			if $p | type == "string" then
				[(.[0]+.[1]+[$p]), [0]]
			else
				[.[0]+[$p], []]
			end
		))[0];
	. as $c
	| p
	| expr_to_path
	| _normalize_path
	| reduce .[] as $n ($c;
		if $n | type == "string" then
			children | map(select(name==$n))
		else
			.[$n]
		end
	  );

def path_array_tree(name; $v):
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

# <mp4 value> | mp4_lookup(".moov.trak[1]")
def mp4_lookup(p): array_tree_path(.boxes; .type; p);
# <mp4 value> | mp4_path(<mp4 value>)
def mp4_path(p): path_array_tree(.type; p);

# <matroska value> | matroska_lookup(".Segment.Tracks[0].TrackEntry[1].CodecID")
def matroska_lookup(p): array_tree_path(.elements; .id._symbol; p);
# <matroska value> | matroska_path(<matroska value>)
def matroska_path(p): path_array_tree(.id._symbol; p);
