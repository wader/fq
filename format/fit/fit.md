### Limitations

- Fields with subcomponents, such as "compressed_speed_distance" field on globalMessageNumber 20 is not represented correctly. 
  The field is read as 3 separate bytes where the first 12 bits are speed and the last 12 bits are distance.
- There are still lots of UNKOWN fields due to gaps in Garmins SDK Profile documentation. (Currently FIT SDK 21.126)
- Compressed timestamp messages are not accumulated against last known full timestamp.

### Convert stream of data messages to JSON array

```
$ fq '[.data_records[] | select(.record_header.message_type == "data").data_message]' file.fit 
```

### Authors
- Mikael Lofjärd mikael.lofjard@gmail.com, original author

### References
- https://developer.garmin.com/fit/protocol/
- https://developer.garmin.com/fit/cookbook/decoding-activity-files/
