Supports reading Avro Object Container Format (OCF) files based on the 1.11.0 specification.

Capable of handling null, deflate, and snappy codecs for data compression.

Limitations:

- Schema does not support self-referential types, only built-in types.
- Decimal logical types are not supported for decoding, will just be treated as their primitive type

### References
- https://avro.apache.org/docs/current/spec.html#Object+Container+Files

### Authors
- Xentripetal
xentripetal@fastmail.com
[@xentripetal](https://github.com/xentripetal)
