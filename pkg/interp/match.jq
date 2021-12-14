def _buffer_fn(f):
  ( . as $c
  | tobytesrange
  | f
  );

def _buffer_try_orig(bfn; fn):
  ( . as $c
  | if type == "string" and (_is_buffer | not) then fn
    else
      ( $c
      | tobytesrange
      | bfn
      )
    end
  );

# overloads to support buffer

def _orig_test($val): test($val);
def _orig_test($regex; $flags): test($regex; $flags);
def _test_buffer($regex; $flags):
  ( isempty(_match_buffer($regex; $flags))
  | not
  );
def test($val): _buffer_try_orig(_test_buffer($val; ""); _orig_test($val));
def test($regex; $flags): _buffer_try_orig(_test_buffer($regex; $flags); _orig_test($regex; $flags));

def _orig_match($val): match($val);
def _orig_match($regex; $flags): match($regex; $flags);
def match($val): _buffer_try_orig(_match_buffer($val); _orig_match($val));
def match($regex; $flags): _buffer_try_orig(_match_buffer($regex; $flags); _orig_match($regex; $flags));

def _orig_capture($val): capture($val);
def _orig_capture($regex; $flags): capture($regex; $flags);
def _capture_buffer($regex; $flags):
  ( . as $b
  | _match_buffer($regex; $flags)
  | .captures
  | map(
      ( select(.name)
      | {key: .name, value: .string}
      )
    )
  | from_entries
  );
def capture($val): _buffer_try_orig(_capture_buffer($val; ""); _orig_capture($val));
def capture($regex; $flags): _buffer_try_orig(_capture_buffer($regex; $flags); _orig_capture($regex; $flags));

def _orig_scan($val): scan($val);
def _orig_scan($regex; $flags): scan($regex; $flags);
def _scan_buffer($regex; $flags):
  ( . as $b
  | _match_buffer($regex; $flags)
  | $b[.offset:.offset+.length]
  );
def scan($val): _buffer_try_orig(_scan_buffer($val; "g"); _orig_scan($val));
def scan($regex; $flags): _buffer_try_orig(_scan_buffer($regex; "g"+$flags); _orig_scan($regex; $flags));

def _orig_splits($val): splits($val);
def _orig_splits($regex; $flags): splits($regex; $flags);
def _splits_buffer($regex; $flags):
  ( . as $b
  # last null output is to do a last iteration that output from end of last match to end of buffer
  | foreach (_match_buffer($regex; $flags), null) as $m (
      {prev: null, curr: null};
      ( .prev = .curr
      | .curr = $m
      );
      if .prev == null then $b[0:.curr.offset]
      elif .curr == null then $b[.prev.offset+.prev.length:]
      else $b[.prev.offset+.prev.length:.curr.offset+.curr.length]
      end
    )
  );
def splits($val): _buffer_try_orig(_splits_buffer($val; "g"); _orig_splits($val));
def splits($regex; $flags): _buffer_try_orig(_splits_buffer($regex; "g"+$flags); _orig_splits($regex; $flags));

# same as regexp.QuoteMeta
def _quote_meta:
  gsub("(?<c>[\\.\\+\\*\\?\\(\\)\\|\\[\\]\\{\\}\\^\\$\\)])"; "\\\(.c)");

def _orig_split($val): split($val);
def _orig_split($regex; $flags): split($regex; $flags);
# split/1 splits on string not regexp
def split($val): [splits($val | _quote_meta)];
def split($regex; $flags): [splits($regex; $flags)];

# TODO: rename
# same as scan but outputs buffer from start of match to end of buffer
def _scan_toend($regex; $flags):
  ( . as $b
  | _match_buffer($regex; $flags)
  | $b[.offset:]
  );
def scan_toend($val): _buffer_fn(_scan_toend($val; "g"));
def scan_toend($regex; $flags):  _buffer_fn(_scan_toend($regex; "g"+$flags));
