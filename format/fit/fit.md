### Limitations

- "compressed_speed_distance" field on globalMessageNumber 20 is not represented correctly. 
  The field is read as 3 separate bytes where the first 12 bits are speed and the last 12 bits are distance.
- There are still lots of UNKOWN fields due to gaps in Garmins SDK Profile documentation. (Currently FIT SDK 21.126)
- Dynamically referenced fields are named incorrectly and lacks scaling, offset and units (just raw values)

### Convert stream of data messages to JSON array

```
$ fq '[.dataRecords[] | select(.dataRecordHeader.messageType == 0).dataMessage]' file.fit 
```

### Authors
- Mikael Lofjärd mikael.lofjard@gmail.com, original author

### References
- https://developer.garmin.com/fit/protocol/
