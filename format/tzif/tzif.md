### Get last transition time
```sh
fq '.v2plusdatablock.transition_times[-1] | tovalue' tziffile
```

### Count leap second records
```sh
fq '.v2plusdatablock.leap_second_records | length' tziffile
```

### Authors
- Takashi Oguma
[@bitbears-dev](https://github.com/bitbears-dev)
[@0xb17bea125](https://twitter.com/0xb17bea125)

### References
- https://datatracker.ietf.org/doc/html/rfc8536
