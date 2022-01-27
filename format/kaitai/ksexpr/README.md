# ksexpr

Kaitai struct expression evaluator.

Parts of the operator code inspired by gojq's operator code.

```sh
# evaluate
$ go run ./cmd/ksexpr '1+2'
3
$ go run ./cmd/ksexpr -input '{"a":1,"b":2}' 'a+b'
3
$ go run ./cmd/ksexpr -input '{"a":1,"b":2}' 'a > b ? "yes" : "no"'
"no"
$ go run ./cmd/ksexpr '(1+("0x2").to_i).to_s'
"3"

# see tokens
$ go run ./cmd/ksexpr --lex 'a+b'
tokIdent a (0-1) <nil>
'+' + (1-2) <nil>
tokIdent b (2-3) <nil>
end  (3-3) <nil>

# see ast
$ go run ./cmd/ksexpr --parse 'a+b'
{
  "LHS": {
    "T": {
      "NS": null,
      "Name": {
        "Str": "a",
        "Span": {
          "Start": 0,
          "Stop": 1
        },
        "V": null
      }
    },
    "Trailers": null
  },
  "Op": 0,
  "RHS": {
    "T": {
      "NS": null,
      "Name": {
        "Str": "b",
        "Span": {
          "Start": 2,
          "Stop": 3
        },
        "V": null
      }
    },
    "Trailers": null
  }
}
```
