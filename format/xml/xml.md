XML can be decoded and encoded into jq values in two ways, elements as object or array.
Which variant to use depends a bit what you want to do. The object variant might be easier
to query for a specific value but array might be easier to use to generate xml or to query
after all elements of some kind etc.

Encoding is done using the `to_xml` function and it will figure what variant that is used based on the input value.
Is has two optional options `indent` and `attribute_prefix`.

### Elements as object

Element can have different shapes depending on body text, attributes and children:

- `<a key="value">text</a>` is `{"a":{"#text":"text","@key":"value"}}`, has text (`#text`) and attributes (`@key`)
- `<a>text</a>` is `{"a":"text"}`
- `<a><b>text</b></a>` is `{"a":{"b":"text"}}` one child with only text and no attributes
- `<a><b/><b>text</b></a>` is `{"a":{"b":["","text"]}}` two children with same name end up in an array
- `<a><b/><b key="value">text</b></a>` is `{"a":{"b":["",{"#text":"text","@key":"value"}]}}`

If there is `#seq` attribute it encodes the child element order. Use `-o seq=true` to include sequence number when decoding,
otherwise order might be lost.

```sh
# decode as object is the default
$ echo '<a><b/><b>bbb</b><c attr="value">ccc</c></a>' | fq -d xml -o seq=true
{
  "a": {
    "b": [
      {
        "#seq": 0
      },
      {
        "#seq": 1,
        "#text": "bbb"
      }
    ],
    "c": {
      "#seq": 2,
      "#text": "ccc",
      "@attr": "value"
    }
  }
}

# access text of the <c> element
$ echo '<a><b/><b>bbb</b><c attr="value">ccc</c></a>' | fq '.a.c["#text"]'
"ccc"

# decode to object and encode to xml
$ echo '<a><b/><b>bbb</b><c attr="value">ccc</c></a>' | fq -r -d xml -o seq=true 'to_xml({indent:2})'
<a>
  <b></b>
  <b>bbb</b>
  <c attr="value">ccc</c>
</a>
```

### Elements as array

Elements are arrays of the shape `["#text": "body text", "attr_name", {key: "attr value"}|null, [<child element>, ...]]`.

```sh
# decode as array
$ echo '<a><b/><b>bbb</b><c attr="value">ccc</c></a>' | fq -d xml -o array=true
[
  "a",
  null,
  [
    [
      "b",
      null,
      []
    ],
    [
      "b",
      {
        "#text": "bbb"
      },
      []
    ],
    [
      "c",
      {
        "#text": "ccc",
        "attr": "value"
      },
      []
    ]
  ]
]

# decode to array and encode to xml
$ echo '<a><b/><b>bbb</b><c attr="value">ccc</c></a>' | fq -r -d xml -o array=true -o seq=true 'to_xml({indent:2})'
<a>
  <b></b>
  <b>bbb</b>
  <c attr="value">ccc</c>
</a>

# access text of the <c> element, the object variant above is probably easier to use
$ echo '<a><b/><b>bbb</b><c attr="value">ccc</c></a>' | fq -o array=true '.[2][2][1]["#text"]'
"ccc"
```

### References
- [xml.com's Converting Between XML and JSON](https://www.xml.com/pub/a/2006/05/31/converting-between-xml-and-json.html)
