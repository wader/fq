# overload match to support buffer
def _orig_match($val): match($val);
def _orig_match($regex; $flags): match($regex; $flags);
def _match_try_buffer(bfn; fn):
  ( . as $c
  | if _is_buffer then bfn
    else
      try fn
      catch
        ( $c
        | tobytesrange
        | bfn
        )
    end
  );
def match($val): _match_try_buffer(_match_buffer($val); _orig_match($val));
def match($regex; $flags): _match_try_buffer(_match_buffer($regex; $flags); _orig_match($regex; $flags));
