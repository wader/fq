Supports decoding BER, CER and DER ([X.690]([X.690_1297.pdf)).

- Currently no extra validation is done for CER and DER.
- Does not support specifying a schema.
- Supports `torepr` but without schema all sequences and sets will be arrays.

```
fq -d asn1_ber torepr file.ber
```

Functions `frompem` and `topem` can help working with PEM format:

```
fq -d raw 'frompem | asn1_ber | d' cert.pem
```

If the schema is known and not that complicated it can be reproduced:

```
fq -d asn1_ber 'torepr as $r | ["version", "modulus", "private_exponent", "private_exponen", "prime1", "prime2", "exponent1", "exponent2", "coefficient"] | with_entries({key: .value, value: $r[.key]})' pkcs1.der
```

Can be used to decode nested parts:

```
fq -d asn1_ber '.constructed[1].value | asn1_ber' file.ber
```

References and tools:
- https://www.itu.int/ITU-T/studygroups/com10/languages/X.690_1297.pdf
- https://en.wikipedia.org/wiki/X.690
- https://letsencrypt.org/docs/a-warm-welcome-to-asn1-and-der/
- https://lapo.it/asn1js/
