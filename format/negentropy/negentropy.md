### View a full Negentropy message

```
$ fq -d negentropy dd file
```

### Check how many ranges the message has and how many of those are of 'fingerprint' mode

```
$ fq -d negentropy '.bounds | length as $total | map(select(.mode == "fingerprint")) | length | {$total, fingerprint: .}' message
```

### Check get all ids in all idlists

```
$ fq -d negentropy '.bounds | map(select(.mode == "idlist") | .idlist | .ids) | flatten' message
```

### Authors
- fiatjaf, https://fiatjaf.com

### References
- https://github.com/hoytech/negentropy
