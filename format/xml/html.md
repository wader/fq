HTML is decoded in HTML5 mode and will always include `<html>`, `<body>` and `<head>` element.

See xml format for more examples and how to preserve element order and how to encode to xml.

There is no `to_html` function, see `to_xml` instead.

### Element as object

```sh
# decode as object is the default
$ echo '<a href="url">text</a>' | fq -d html
{
  "html": {
    "body": {
      "a": {
        "#text": "text",
        "@href": "url"
      }
    },
    "head": ""
  }
}
```

### Element as array

```sh
$ '<a href="url">text</a>' | fq -d html -o array=true
[
  "html",
  null,
  [
    [
      "head",
      null,
      []
    ],
    [
      "body",
      null,
      [
        [
          "a",
          {
            "#text": "text",
            "href": "url"
          },
          []
        ]
      ]
    ]
  ]
]

# decode html files to a {file: "title", ...} object
$ fq -n -d html '[inputs | {key: input_filename, value: .html.head.title?}] | from_entries' *.html

# <a> href:s in file
$ fq -r -o array=true -d html '.. | select(.[0] == "a" and .[1].href)?.[1].href' file.html
```
