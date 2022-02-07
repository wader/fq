Can also decode CER and BER but with no extra validation. Currently does not support specifying a schema.

Supports `torepr` but without schema support it's not that useful:

```
fq -d asn1_ber torepr file.ber
```

There is also `frompem` and `topem` to help working with PEM format:

```
fq -d raw 'frompem | asn1_ber | d' cert.pem
```

If the schema is known and not that complicated it can be reproduced:

```
fq -d asn1_ber 'torepr as $r | ["version", "modulus", "private_exponent", "private_exponen", "prime1", "prime2", "exponent1", "exponent2", "coefficient"] | with_entries({key: .value, value: $r[.key]})' pkcs1.der
```