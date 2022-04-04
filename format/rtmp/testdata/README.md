Some files based on pcap:s from https://wiki.wireshark.org/SampleCaptures

ffmpeg_*_stream files were generated using:

```sh
ffmpeg -hide_banner -f lavfi -i sine=frequency=500:sample_rate=44100 -acodec aac -ac 2 -f flv -t 200ms -metadata streamName="test stream" -listen 1 rtmp://0.0.0.0:1935/test_stream
ffmpeg -i rtmp://localhost:1935/test_stream -f null -
tcpdump -w rtmp.pcap -i lo0
fq '.tcp_connections[1].client_stream | tobytes' rtmp.pcap > ffmpeg_client_stream
fq '.tcp_connections[1].server_stream | tobytes' rtmp.pcap > ffmpeg_server_stream
```
