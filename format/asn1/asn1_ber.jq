def _asn1_ber_torepr:
  def _f:
    if .class == "universal" then
      if .tag | . == "sequence" or . == "set" then
        .constructed | map(_f)
      else .value | tovalue
      end
    else .constructed | map(_f)
    end;
  _f;
