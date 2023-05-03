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

`bplist` files can be converted to a JSON representation using the `torepr` filter:
```sh
$ fq torepr com.apple.UIAutomation.plist
{
  "UIAutomationEnabled": true
}
```

### Decoding NSKeyedArchiver serialized objects

A common way that Swift and Objective-C libraries on macOS serialize objects
is through the NSKeyedArchiver API, which flattens objects into a list of elements
and class descriptions that are reconstructed into an object graph using CFUID
elements in the property list. `fq` includes a function, `from_ns_keyed_archiver`,
which will rebuild this object graph into a friendly representation. 

If no parameters are supplied, it will assume that there is a CFUID located at
`."$top".root` that specifies the root from which decoding should occur. If this
is not present, an error will be produced, asking the user to specify a root
object in the `.$objects` list from which to decode.

The following examples show how this might be used (in this case, within the `fq` REPL):
```
# Assume $top.root is present
bplist> from_ns_keyed_archiver

# Specify optional root
bplist> from_ns_keyed_archiver(1)
```

### Authors
- David McDonald
[@dgmcdona](https://github.com/dgmcdona)

### References
- http://fileformats.archiveteam.org/wiki/Property_List/Binary
- https://medium.com/@karaiskc/understanding-apples-binary-property-list-format-281e6da00dbd
- https://opensource.apple.com/source/CF/CF-550/CFBinaryPList.c
