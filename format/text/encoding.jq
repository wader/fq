def to_iso8859_1: _to_strencoding({encoding: "ISO8859_1"});
def from_iso8859_1: _from_strencoding({encoding: "ISO8859_1"});
def to_utf8: _to_strencoding({encoding: "UTF8"});
def from_utf8: _from_strencoding({encoding: "UTF8"});
def to_utf16: _to_strencoding({encoding: "UTF16"});
def from_utf16: _from_strencoding({encoding: "UTF16"});
def to_utf16le: _to_strencoding({encoding: "UTF16LE"});
def from_utf16le: _from_strencoding({encoding: "UTF16LE"});
def to_utf16be: _to_strencoding({encoding: "UTF16BE"});
def from_utf16be: _from_strencoding({encoding: "UTF16BE"});

def from_base64($opts): _from_base64({encoding: "std"} + $opts);
def from_base64: _from_base64(null);
def to_base64($opts): _to_base64({encoding: "std"} + $opts);
def to_base64: _to_base64(null);

# TODO: compat: remove at some point
def hex: _binary_or_orig(to_hex; from_hex);
def base64: _binary_or_orig(to_base64; from_base64);
def tohex: to_hex;
def fromhex: from_hex;
