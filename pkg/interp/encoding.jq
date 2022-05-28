include "internal";
include "binary";

# convert all scalars to strings, null as empty string (same as @csv)
def _walk_tostring:
  walk(
    if _is_null then ""
    elif _is_scalar then tostring
    end
  );
# overloads builtin tojson to have options
def tojson($opts): _tojson({} + $opts);
def tojson: tojson(null);

def fromxml($opts): _fromxml({} + $opts);
def fromxml: _fromxml(null);
def toxml($opts): _walk_tostring | _toxml({} + $opts);
def toxml: toxml(null);

def fromhtml($opts): _fromhtml({} + $opts);
def fromhtml: fromhtml(null);

def fromyaml: _fromyaml;
def toyaml: _toyaml;

def fromtoml: _fromtoml;
def totoml: _totoml;

def fromcsv($opts): _fromcsv({comma: ",", comment: "#"} + $opts);
def fromcsv: fromcsv(null);
def tocsv($opts): _walk_tostring | _tocsv({comma: ","} + $opts);
def tocsv: tocsv(null);

def fromxmlentities: _fromxmlentities;
def toxmlentities: _toxmlentities;

def fromurlpath: _fromurlpath;
def tourlpath: _tourlpath;

def fromurlencode: _fromurlencode;
def tourlencode: _tourlencode;

def fromurlquery: _fromurlquery;
def tourlquery: _tourlquery;

def fromurl: _fromurl;
def tourl: _tourl;

def fromhex: _fromhex;
def tohex: _tohex;

def frombase64($opts): _frombase64({encoding: "std"} + $opts);
def frombase64: _frombase64(null);
def tobase64($opts): _tobase64({encoding: "std"} + $opts);
def tobase64: _tobase64(null);

def tomd4: _tohash({name: "md4"});
def tomd5: _tohash({name: "md5"});
def tosha1: _tohash({name: "sha1"});
def tosha256: _tohash({name: "sha256"});
def tosha512: _tohash({name: "sha512"});
def tosha3_224: _tohash({name: "sha3_224"});
def tosha3_256: _tohash({name: "sha3_256"});
def tosha3_384: _tohash({name: "sha3_384"});
def tosha3_512: _tohash({name: "sha3_512"});

# _tostrencoding/_fromstrencoding can do more but not exposed as functions yet
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

# https://en.wikipedia.org/wiki/Privacy-Enhanced_Mail
# TODO: add test
def frompem:
  ( tobytes
  | tostring
  | capture("-----BEGIN(.*?)-----(?<s>.*?)-----END(.*?)-----"; "mg").s
  | frombase64
  ) // error("no pem header or footer found");

def topem($label):
  ( tobytes
  | tobase64
  | ($label | if $label != "" then " " + $label end) as $label
  | [ "-----BEGIN\($label)-----"
    , .
    , "-----END\($label)-----"
    , ""
    ]
  | join("\n")
  );
def topem: topem("");

def fromradix($base; $table):
  ( if _is_string | not then error("cannot fromradix convert: \(.)") end
  | split("")
  | reverse
  | map($table[.])
  | if . == null then error("invalid char \(.)") end
  # state: [power, ans]
  | reduce .[] as $c ([1,0];
      ( (.[0] * $base) as $b
      | [$b, .[1] + (.[0] * $c)]
      )
    )
  | .[1]
  );
def fromradix($base):
  fromradix($base; {
    "0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
    "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16,
    "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23,
    "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30,
    "v": 31, "w": 32, "x": 33, "y": 34, "z": 35,
    "A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42,
    "H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49,
    "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56,
    "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61,
    "@": 62, "_": 63,
  });

def toradix($base; $table):
  ( if type != "number" then error("cannot toradix convert: \(.)") end
  | if . == 0 then "0"
    else
      ( [ recurse(if . > 0 then _intdiv(.; $base) else empty end) | . % $base]
      | reverse
      | .[1:]
      | if $base <= ($table | length) then
          map($table[.]) | join("")
        else
          error("base too large")
        end
      )
    end
  );
def toradix($base):
  toradix($base; "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@_");

# to jq-flavoured json
def _tojq($opts):
  def _is_ident: test("^[a-zA-Z_][a-zA-Z_0-9]*$");
  def _key: if _is_ident | not then tojson end;
  def _f($opts; $indent):
    def _r($prefix):
      ( type as $t
      | if $t == "null" then tojson
        elif $t == "string" then tojson
        elif $t == "number" then tojson
        elif $t == "boolean" then tojson
        elif $t == "array" then
          [ "[", $opts.compound_newline
          , ( [ .[]
              | $prefix, $indent
              , _r($prefix+$indent), $opts.array_sep
              ]
            | .[0:-1]
            )
          , $opts.compound_newline
          , $prefix, "]"
          ]
        elif $t == "object" then
          [ "{", $opts.compound_newline
          , ( [ to_entries[]
              | $prefix, $indent
              , (.key | _key), $opts.key_sep
              , (.value | _r($prefix+$indent)), $opts.value_sep
              ]
            | .[0:-1]
            )
          , $opts.compound_newline
          , $prefix, "}"
          ]
        else error("unknown type \($t)")
        end
      );
    _r("");
  ( _f($opts; $opts.indent * " ")
  | if _is_array then flatten | join("") end
  );
def tojq($opts):
  _tojq(
    ( { indent: 0,
        key_sep: ":",
        value_sep: ",",
        array_sep: ",",
        compound_newline: "",
      } + $opts
    | if .indent > 0  then
        ( .key_sep = ": "
        | .value_sep = ",\n"
        | .array_sep = ",\n"
        | .compound_newline = "\n"
        )
      end
    )
  );
def tojq: tojq(null);

# from jq-flavoured json
def fromjq:
  def _f:
    ( . as $v
    | .term.type
    | if . == "TermTypeNull" then null
      elif . == "TermTypeTrue" then true
      elif . == "TermTypeFalse" then false
      elif . == "TermTypeString" then $v.term.str.str
      elif . == "TermTypeNumber" then $v.term.number | tonumber
      elif . == "TermTypeObject" then
        ( $v.term.object.key_vals
        | map(
            { key: (.key // .key_string.str),
              value: (.val.queries[0] | _f)
            }
          )
        | from_entries
        )
      elif . == "TermTypeArray" then
        ( def _a: if .op then .left, .right | _a end;
          [$v.term.array.query | _a | _f]
        )
      else error("unknown term")
      end
    );
  try
    (_query_fromstring | _f)
  catch
    error("fromjq only supports constant literals");

# TODO: compat remove at some point
def hex: _binary_or_orig(tohex; fromhex);
def base64: _binary_or_orig(tobase64; frombase64);
