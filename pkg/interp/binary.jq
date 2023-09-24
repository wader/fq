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

def _orig_split($val): split($val);
def _orig_split($regex; $flags): split($regex; $flags);
# split/1 splits on string not regexp
def split($val): [splits($val | _re_quote_meta)];
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
