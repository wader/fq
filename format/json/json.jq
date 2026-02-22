# overrides jq's standard tojson
def tojson($opts): _to_json($opts);
def tojson: _to_json(null);
# overrides jq's standard fromjson
# NOTE:
# should be kept in sync with format_decode.jq and can't use from_json as
# it's not defined yet.
# also uses tovalue on it's input to care of the where the input is a decode_value
# string which without would end up decoding the "backing" binary instead.
# Ex:
# $ fq -n '"\"1,2,3\"" | fromjson | tobytes | tostring'
# "\"1,2,3\""
def fromjson: tovalue | decode("json") | if ._error then error(._error.error) end;

def _json__todisplay: tovalue;
