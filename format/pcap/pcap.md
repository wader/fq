### Build object with number of (reassembled) TCP bytes sent to/from client IP
```sh
# for a pcapng file you would use .[0].tcp_connections for first section
$ fq '.tcp_connections | group_by(.client.ip) | map({key: .[0].client.ip, value: map(.client.stream, .server.stream | tobytes.size) | add}) | from_entries'
{
  "10.1.0.22": 15116,
  "10.99.12.136": 234,
  "10.99.12.150": 218
}
```