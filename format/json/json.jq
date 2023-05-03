# overrides jq's standard tojson
def tojson($opts): _to_json($opts);
def tojson: _to_json(null);
# overrides jq's standard fromjson
# NOTE: should be kept in sync with format_decode.jq
def fromjson: decode("json") | if ._error then error(._error.error) end;

def _json__todisplay: tovalue;
