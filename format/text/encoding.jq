def toiso8859_1: _tostrencoding({encoding: "ISO8859_1"});
def fromiso8859_1: _fromstrencoding({encoding: "ISO8859_1"});
def toutf8: _tostrencoding({encoding: "UTF8"});
def fromutf8: _fromstrencoding({encoding: "UTF8"});
def toutf16: _tostrencoding({encoding: "UTF16"});
def fromutf16: _fromstrencoding({encoding: "UTF16"});
def toutf16le: _tostrencoding({encoding: "UTF16LE"});
def fromutf16le: _fromstrencoding({encoding: "UTF16LE"});
def toutf16be: _tostrencoding({encoding: "UTF16BE"});
def fromutf16be: _fromstrencoding({encoding: "UTF16BE"});

def frombase64($opts): _frombase64({encoding: "std"} + $opts);
def frombase64: _frombase64(null);
def tobase64($opts): _tobase64({encoding: "std"} + $opts);
def tobase64: _tobase64(null);

# TODO: compat: remove at some point
def hex: _binary_or_orig(tohex; fromhex);
def base64: _binary_or_orig(tobase64; frombase64);
