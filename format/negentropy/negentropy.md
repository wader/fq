### View a full Negentropy message

```
$ fq -d negentropy dd file
```

### Or from hex

```
$ echo '6186b7abb47c0001108e4206828ee3bf34258465809a337c6c00019a68e37b177a50b3ae7164ccc628b962020114019c1381281c9e3849d5fbd514b7bb65ad0101e601fbf7451f5d22e7fa36ae3e910e9f5215020157014a1b26853e06e9c32eb41b1df4f9ab300201e6011840e273c84bb1344f1d4e15d9aa67920200016f12ee2340888653f10b0ec2d438ac9f0101840156d2d796f4dff004ab369b9bcfa4d81e020187013f1b3c8a019800d5764e2de6bdfd2785020114017caaf0acb5dfe249aa0f7f742402168a01018301e7b8c4decb1eae455ca5714281e3245302017a01409c22636b097362df125ddffb6d944302015b01f332208bee82acf8ed922853ee54057f020001fc3e51fdb0b92966e38017f7959903850101cc01428ce0c96d49f15b50143e4fb228cb9300000131712d30e5296a7a45d07bba452d61cd' | fq -R 'from_hex | negentropy | dd'
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
