# overload match to support buffer
def _orig_match($val): match($val);
def _orig_match($regex; $flags): match($regex; $flags);
def match($val): if _is_buffer then _match_buffer($val) else _orig_match($val) end;
def match($regex; $flags): if _is_buffer then _match_buffer($regex; $flags) else _orig_match($regex; $flags) end;
