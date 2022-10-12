Current only supports plain RTMP (not RTMPT or encrypted variants etc) with AMF0 (not AMF3).

### Show rtmp streams in PCAP file
```sh
fq '.tcp_connections[] | select(.server.port=="rtmp") | d' file.cap
```

### References
- https://rtmp.veriskope.com/docs/spec/
- https://rtmp.veriskope.com/pdf/video_file_format_spec_v10.pdf
