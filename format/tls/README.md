## Dev notes

TLS deflate compression seems to actually be zlib, so zlib header + deflate. Also each record is compressed with a flush (trailing 0x00 0x00 0xff 0xff) so that they can be uncompressed individually.

https://lekensteyn.nl/files/wireshark-ssl-tls-decryption-secrets-sharkfest18eu.pdf

```
tshark -x -V -o tls.keylog_file:file.keylog -r file.pcap
```

Wireshark gui has TLS debug option to write key/iv etc

```
tcpdump -i en0 -w file.pcap
SSLKEYLOGFILE=file.keylog /path/to/sslkey-able/curl --http1.1 -tlsv1.2 --tls-max 1.2 -v https://host/path
```

TLS 1.3 dumps https://gitlab.com/wireshark/wireshark/-/issues/12779
