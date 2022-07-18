def _asn1_ber_torepr:
  if .class == "universal" then
    if .tag | . == "sequence" or . == "set" then
      .constructed | map(_asn1_ber_torepr)
    else .value | tovalue
    end
  else .constructed | map(_asn1_ber_torepr)
  end;

def _asn1_ber__help:
  { notes: "Supports decoding BER, CER and DER (X.690).

- Currently no extra validation is done for CER and DER.
- Does not support specifying a schema.
- Supports `torepr` but without schema all sequences and sets will be arrays.",
    examples: [
      {comment: "`frompem` and `topem` can be used to work with PEM format", shell: "fq -d raw 'frompem | asn1_ber | d' cert.pem"},
      {comment: "Can be used to decode nested parts", shell: "fq -d asn1_ber '.constructed[1].value | asn1_ber' file.ber"},
      { comment: "If schema is known and not complicated it can be reproduced",
        shell: "fq -d asn1_ber 'torepr as $r | [\"version\", \"modulus\", \"private_exponent\", \"private_exponen\", \"prime1\", \"prime2\", \"exponent1\", \"exponent2\", \"coefficient\"] | with_entries({key: .value, value: $r[.key]})' pkcs1.der"
      }
    ],
    links: [
      {url: "https://www.itu.int/ITU-T/studygroups/com10/languages/X.690_1297.pdf"},
      {url: "https://en.wikipedia.org/wiki/X.690"},
      {url: "https://letsencrypt.org/docs/a-warm-welcome-to-asn1-and-der/"},
      {url: "https://lapo.it/asn1js/"}
    ]
  };
