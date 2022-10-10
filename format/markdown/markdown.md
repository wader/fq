### Array with all level 1 and 2 headers
```sh
$ fq -d markdown '[.. | select(.type=="heading" and .level<=2)?.children[0]]' file.md
```