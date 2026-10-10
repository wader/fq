def tobits: _tobits({unit: 1, keep_range: false, pad_to_units: 0});
def tobytes: _tobits({unit: 8, keep_range: false, pad_to_units: 0});
def tobitsrange: _tobits({unit: 1, keep_range: true, pad_to_units: 0});
def tobytesrange: _tobits({unit: 8, keep_range: true, pad_to_units: 0});
def tobits($pad): _tobits({unit: 1, keep_range: false, pad_to_units: $pad});
def tobytes($pad): _tobits({unit: 8, keep_range: false, pad_to_units: $pad});

# same as regexp.QuoteMeta
def _re_quote_meta:
  gsub("(?<c>[\\.\\+\\*\\?\\(\\)\\|\\[\\]\\{\\}\\^\\$\\)])"; "\\\(.c)");

# TODO:
# maybe implode, join. but what would it mean?
# "abc" | tobits | explode | implode would not work

# helper for overloading regex/string functions to support binary
def _binary_or_orig(bfn; fn):
  if _exttype == "binary" then bfn
  else fn
  end;
def _bytes_or_orig(bfn; fn):
  _binary_or_orig(
    # convert to bytes if bits
    ( if .unit != 8 then tobytesrange end
    | bfn
    );
    fn
  );

def _orig_explode: explode;
def explode: _binary_or_orig([.[range(.size)]]; _orig_explode);

def _orig_split($val): split($val);
def _orig_split($regex; $flags): split($regex; $flags);
def _orig_splits($val): splits($val);
def _orig_splits($regex; $flags): splits($regex; $flags);

def _splits_binary($regex; $flags):
  ( . as $b
  # last null output is to do a last iteration that output from end of last match to end of binary
  | foreach (_match_binary($regex; $flags), null) as $m (
      {prev: null, current: null};
      ( .prev = .current
      | .current = $m
      );
      if .prev == null then $b[0:.current.offset]
      elif .current == null then $b[.prev.offset+.prev.length:]
      else $b[.prev.offset+.prev.length:.current.offset]
      end
    )
  );
def splits($val): _bytes_or_orig(_splits_binary($val; "g"); _orig_splits($val));
def splits($regex; $flags): _bytes_or_orig(_splits_binary($regex; "g"+$flags); _orig_splits($regex; $flags));

# split/1 splits on string not regexp
def split($val): _bytes_or_orig([_splits_binary($val; "g")]; _orig_split($val));
def split($regex; $flags): [splits($regex; $flags)];

def _orig_test($val): test($val);
def _orig_test($regex; $flags): test($regex; $flags);
def _test_binary($regex; $flags):
  ( isempty(_match_binary($regex; $flags))
  | not
  );
def test($val): _bytes_or_orig(_test_binary($val; ""); _orig_test($val));
def test($regex; $flags): _bytes_or_orig(_test_binary($regex; $flags); _orig_test($regex; $flags));

def _orig_match($val): match($val);
def _orig_match($regex; $flags): match($regex; $flags);
def match($val): _bytes_or_orig(_match_binary($val; ""); _orig_match($val));
def match($regex; $flags): _bytes_or_orig(_match_binary($regex; $flags); _orig_match($regex; $flags));

def _orig_capture($val): capture($val);
def _orig_capture($regex; $flags): capture($regex; $flags);
def _capture_binary($regex; $flags):
  ( . as $b
  | _match_binary($regex; $flags)
  | .captures
  | map(
      ( select(.name)
      | {key: .name, value: .string}
      )
    )
  | from_entries
  );
def capture($val): _bytes_or_orig(_capture_binary($val; ""); _orig_capture($val));
def capture($regex; $flags): _bytes_or_orig(_capture_binary($regex; $flags); _orig_capture($regex; $flags));

def _orig_scan($val): scan($val);
def _orig_scan($regex; $flags): scan($regex; $flags);
def _scan_binary($regex; $flags):
  ( . as $b
  | _match_binary($regex; $flags)
  | $b[.offset:.offset+.length]
  );
def scan($val): _bytes_or_orig(_scan_binary($val; "g"); _orig_scan($val));
def scan($regex; $flags): _bytes_or_orig(_scan_binary($regex; "g"+$flags); _orig_scan($regex; $flags));

# The functions below were left out when the overloads above were written, so
# they still reach a binary through a coercion that replaces every byte it
# cannot read as UTF-8. Each one now works on the bytes and gives back a binary
# wherever the string version gives back a string.

def _orig_startswith($v): startswith($v);
def startswith($v): _bytes_or_orig(_binary_starts_with($v); _orig_startswith($v));

def _orig_endswith($v): endswith($v);
def endswith($v): _bytes_or_orig(_binary_ends_with($v); _orig_endswith($v));

def _orig_indices($v): indices($v);
def indices($v): _bytes_or_orig(_binary_indices($v); _orig_indices($v));

def _orig_index($v): index($v);
def index($v):
  _bytes_or_orig(
    ( _binary_indices($v)
    | if length == 0 then null else .[0] end
    );
    _orig_index($v)
  );

def _orig_rindex($v): rindex($v);
def rindex($v):
  _bytes_or_orig(
    ( _binary_indices($v)
    | if length == 0 then null else .[-1] end
    );
    _orig_rindex($v)
  );

def _orig_contains($v): contains($v);
def contains($v):
  _bytes_or_orig(
    ( if _binary_needle_length($v) == 0 then true
      else _binary_indices($v) | length > 0
      end
    );
    _orig_contains($v)
  );

def _orig_ltrimstr($v): ltrimstr($v);
def ltrimstr($v):
  _bytes_or_orig(
    ( _binary_needle_length($v) as $n
    | if _binary_starts_with($v) then .[$n:]
      else .
      end
    );
    _orig_ltrimstr($v)
  );

def _orig_rtrimstr($v): rtrimstr($v);
def rtrimstr($v):
  _bytes_or_orig(
    ( _binary_needle_length($v) as $n
    | if $n > 0 and _binary_ends_with($v) then .[0:length-$n]
      else .
      end
    );
    _orig_rtrimstr($v)
  );

def _orig_ascii_downcase: ascii_downcase;
def ascii_downcase: _bytes_or_orig(_binary_ascii_case(false); _orig_ascii_downcase);

def _orig_ascii_upcase: ascii_upcase;
def ascii_upcase: _bytes_or_orig(_binary_ascii_case(true); _orig_ascii_upcase);

# sub and gsub were left out of the overloads above even though match, scan,
# split and capture were done, so they still went through the coercion. The
# replacement is evaluated once for each match with that match's capture object
# as its input, the way the string versions evaluate it, and its first value is
# used. A replacement may be a string or a binary.

def _capture_object_of($m):
  ( $m.captures
  | map(select(.name) | {key: .name, value: .string})
  | from_entries
  );

def _sub_binary($regex; str; $flags):
  ( . as $b
  | [ foreach (_match_binary($regex; $flags), null) as $m (
        {prev: 0, parts: []};
        if $m == null then
          {prev: .prev, parts: (.parts + [$b[.prev:]])}
        else
          { prev: ($m.offset + $m.length),
            parts: (.parts + [$b[.prev:$m.offset], first(_capture_object_of($m) | str)])
          }
        end
      )
    ]
  | last.parts
  | tobytes
  );

def _orig_sub($regex; str): sub($regex; str);
def _orig_sub($regex; str; $flags): sub($regex; str; $flags);
def sub($regex; str): _bytes_or_orig(_sub_binary($regex; str; ""); _orig_sub($regex; str));
def sub($regex; str; $flags): _bytes_or_orig(_sub_binary($regex; str; $flags); _orig_sub($regex; str; $flags));

def _orig_gsub($regex; str): gsub($regex; str);
def _orig_gsub($regex; str; $flags): gsub($regex; str; $flags);
def gsub($regex; str): _bytes_or_orig(_sub_binary($regex; str; "g"); _orig_gsub($regex; str));
def gsub($regex; str; $flags): _bytes_or_orig(_sub_binary($regex; str; "g"+$flags); _orig_gsub($regex; str; $flags));
