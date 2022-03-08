def _asn1_ber_torepr:
  if .class == "universal" then
    if .tag | . == "sequence" or . == "set" then
      .constructed | map(_asn1_ber_torepr)
    else .value | tovalue
    end
  else .constructed | map(_asn1_ber_torepr)
  end;
