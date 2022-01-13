Can also decode CER and BER but with no extra validation. Currently does not support specifying a schema.

Supports `torepr` but without schema support it's not that useful:

```
fq -d asn1_ber torepr file.ber
```
