### View a full OpenTimestamps file

```
$ fq dd file.ots
```

### List the names of the Calendar servers used

```
$ fq '.operations | map(select(.attestation_type == "calendar") | .url)' file.ots
```

### Check if there are Bitcoin attestations present

```
$ fq '.operations | map(select(.attestation_type == "bitcoin")) | length > 0' file.ots
```

### Authors
- fiatjaf, https://fiatjaf.com

### References
- https://opentimestamps.org/
- https://github.com/opentimestamps/python-opentimestamps
