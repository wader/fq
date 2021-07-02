# same as group_by but just counts, grouping exp value will be converted to string
def count_by(exp):
	group_by(exp) | map({(.[0] | exp | tostring): length}) | add;

def protobuf_to_value:
	.fields | map({(.name|tostring): (.enum // .value)}) | add;

# hack to parse just a box
# <binary> | mp4_box
def mp4_box:
	[0,0,0,16, "ftyp", "isom", 0, 0 , 2 ,0, .] | mp4.boxes;

def flac_dump:
	["fLaC", first(.. | select(._format=="flac_metadatablock")), (.. | select(._format=="flac_frame"))] | bits;


# def recurse_foreach(init; update; extract):
#     def _recurse_foreach($state; $c):
#         (. as $c | ["_recurse_foreach $state", $state] | debug | $c) |
#         (. as $b | ["_recurse_foreach $c", $c] | debug | $b) |
#         foreach $c[]? as $e (
#             $state
#             ;
#             (. as $c | ["update $state", $state] | debug | $c) |
#             (. as $b | ["update $e", $e] | debug | $b) |
#             null | _recurse_foreach({state: $state, e: $e} | update; $e)
#             ;
#             (. as $c | ["extract $state", $state] | debug | $c) |
#             (. as $b | ["extract $e", $e] | debug | $b) |
#             {state: $state, e: .} |
#             extract
#         )
#         // $c;
#     _recurse_foreach(init; .);

# # TODO: split? can't really switch on type
# def grep(f):
# 	if f | type == "string" then
# 		.. | select((._name | contains(f)) or (._value | contains(f)? // false))
# 	elif f | type == "number" then
# 		.. | select(._value == f)
# 	else
# 		.. | debug | select(f)?
# 	end;
