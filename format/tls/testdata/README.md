dump.pcapng dump-broken.pcapng was created by Peter Wu and comes from https://bugs.wireshark.org/bugzilla/show_bug.cgi?id=9144.

dump.pcapng contains 73 tls connections with differens cipher suites. split.jq was used to split it into one pcap per connection named after cipher suit used.

dump-broken.pcapng is a broken SSL v3, uses extensions. dump-broken.pcapng.keylog not used yet.
