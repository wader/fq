# same as group_by but just counts, grouping exp value will be converted to string
def count_by(exp):
	group_by(exp) | map({(.[0] | exp | tostring): length}) | add;

def protobuf_to_value:
	.fields | map({(.name|tostring): (.enum // .value)}) | add;

# hack to parse just a box
# <binary> | mp4_box
def mp4_box:
	[0,0,0,16, "ftyp", "isom", 0, 0 , 2 ,0, .] | mp4.boxes;

def _lookup(children; name; p):
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
	| null
	| path(p)
	| _normalize_path
	| reduce .[] as $n ($c;
		if $n | type == "string" then
			children | map(select(name==$n))
		else
			nth($n)
		end
	  );

# <mp4 value> | mp4_lookup(.moov.trak[1])
def mp4_lookup(p): _lookup(.boxes; .type; p);

# <matroska value> | matroska_lookup(.Segment.Tracks[0].TrackEntry[1].CodecID)
def matroska_lookup(p): _lookup(.elements; .id._symbol; p);
