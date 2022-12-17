### Show full decoding
```sh
$ fq d Info.plist
```

### Timestamps
Timestamps in Apple Binary Property Lists are encoded as Cocoa Core Data
timestamps, where the raw value is the floating point number of seconds since
January 1, 2001. By default, `fq` will render the raw floating point value. In
order to get the raw value or string description, use the `todescription`
function, you can use the `tovalue` and `todescription` functions:

```sh
$ fq 'torepr.SomeTimeStamp | tovalue' Info.plist
685135328

$ fq 'torepr.SomeTimeStamp | todescription' Info.plist
"2022-09-17T19:22:08Z"
```


### Get JSON representation
```sh
$ fq torepr com.apple.UIAutomation.plist
{
  "UIAutomationEnabled": true
}
```

### Authors
- David McDonald
[@dgmcdona](https://github.com/dgmcdona)

### References
- http://fileformats.archiveteam.org/wiki/Property_List/Binary
- https://medium.com/@karaiskc/understanding-apples-binary-property-list-format-281e6da00dbd
- https://opensource.apple.com/source/CF/CF-550/CFBinaryPList.c
