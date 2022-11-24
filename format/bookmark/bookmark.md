## Apple bookmarkData format

Apple's `bookmarkData` format is used to encode information that can be resolved
into a `URL` object for a file even if the user moves or renames it. Can also
contain security scoping information for App Sandbox support.

These `bookmarkData` blobs are often found endcoded in data fields of Binary
Property Lists. Notable examples include:

- `com.apple.finder.plist` - contains an `FXRecentFolders` value, which is an
  array of ten objects, each of which consists of a `name` and `file-bookmark`
  field, which is a `bookmarkData` object for each recently accessed folder
  location.

- `com.apple.LSSharedFileList.RecentApplications.sfl2` - `sfl2` files are
  actually `plist` files of the `NSKeyedArchiver` format. They can be parsed the
  same as `plist` files, but they have a more complicated tree-like structure
  than would typically be found. For more information about these types of files,
  see [Sarah Edwards'](https://www.mac4n6.com/blog/2016/1/1/manual-analysis-of-nskeyedarchiver-formatted-plist-files-a-review-of-the-new-os-x-1011-recent-items)
  excellent research on the subject.

Locating `bookmarkData` objects in `NSKeyedArchiver` plist files is a place
where `fq`'s recursive searching really shines:
```
fq '.. | select(format=="bookmark") | .map(. | torepr)' com.apple.LSSharedFileList.RecentApplications.sfl2
```

### Authors
- David McDonald
[@dgmcdona](https://github.com/dgmcdona)

### References
- https://mac-alias.readthedocs.io/en/latest/bookmark_fmt.html
- https://www.mac4n6.com/blog/2016/1/1/manual-analysis-of-nskeyedarchiver-formatted-plist-files-a-review-of-the-new-os-x-1011-recent-items
- https://michaellynn.github.io/2015/10/24/apples-bookmarkdata-exposed/
