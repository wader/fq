package http

// TODO: pipeline
// TODO: range request reassembly?
// TODO: Trailer, only chunked?
// TODO: mime_multi_part decoder?
// TODO: content_type group? pass gzip pass along content_type?
// TODO: text/*  and encoding? ISO-8859-1?
// TODO: PRI * HTTP/2.0, h2?
// TODO: 101 Switch protocol, Connection: Upgrade

/*

echo reqbody | curl --trace bla -H "Transfer-Encoding: chunked" -d @- http://0:8080
while true ; do echo -e 'HTTP/1.0 200 OK\r\nrespbody' | nc -v -l 8080 ; done

split("\n") | reduce .[] as $l ({state: "send",  send: [], recv: []}; if $l | startswith("=>") then .state="send" elif $l | startswith("<=") then .state="recv" elif $l | test("^\\d") then .[.state] += [$l] end) | .["send", "recv"] |= (map(capture(": (?<hex>.{1,47})").hex | gsub(" "; "")) | add | hex) | .send | http | d

*/

/*


   Content-Type: multipart/form-data; boundary=AaB03x

   --AaB03x
   Content-Disposition: form-data; name="submit-name"

   Larry
   --AaB03x
   Content-Disposition: form-data; name="files"; filename="file1.txt"
   Content-Type: text/plain

   ... contents of file1.txt ...
   --AaB03x--

   Multi file:

   Content-Type: multipart/form-data; boundary=AaB03x

   --AaB03x
   Content-Disposition: form-data; name="submit-name"

   Larry
   --AaB03x
   Content-Disposition: form-data; name="files"
   Content-Type: multipart/mixed; boundary=BbC04y

   --BbC04y
   Content-Disposition: file; filename="file1.txt"
   Content-Type: text/plain

   ... contents of file1.txt ...
   --BbC04y
   Content-Disposition: file; filename="file2.gif"
   Content-Type: image/gif
   Content-Transfer-Encoding: binary

   ...contents of file2.gif...
   --BbC04y--
   --AaB03x--

*/

import (
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/lazyre"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var httpContentTypeGroup decode.Group
var httpTextprotoGroup decode.Group
var httpHttpChunkedGroup decode.Group
var httpGzipGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.HTTP,
		&decode.Format{
			Description: "Hypertext Transfer Protocol 1 and 1.1", // TODO: and v1.1?
			Groups:      []*decode.Group{format.TCP_Stream},
			DecodeFn:    httpDecode,
			RootArray:   true,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Content_Type}, Out: &httpContentTypeGroup},
				{Groups: []*decode.Group{format.TextProto}, Out: &httpTextprotoGroup},
				{Groups: []*decode.Group{format.HTTP_Chunked}, Out: &httpHttpChunkedGroup},
				{Groups: []*decode.Group{format.Gzip}, Out: &httpGzipGroup},
			},
		})
}

func headersFirst(m map[string][]string, key string) string {
	for k, vs := range m {
		if strings.EqualFold(k, key) {
			return vs[0]
		}
	}
	return ""
}

// https://www.rfc-editor.org/rfc/rfc6750#section-3
type Pairs struct {
	Scheme string
	Params map[string]string
}

// quoteSplit splits but respects quotes and escapes, and can mix quotes
func quoteSplit(s string, sep rune) ([]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	// allows mix quotes and explicit ","
	r.LazyQuotes = true
	r.Comma = sep
	return r.Read()
}

// multipart/form-data; boundary=...
// form-data; name="aaa_file"; filename="aaa"
func parsePairs(s string) (Pairs, error) {
	var w Pairs
	parts := strings.SplitN(s, ";", 2)
	if len(parts) < 1 {
		return Pairs{}, fmt.Errorf("invalid params")
	}
	w.Scheme = parts[0]
	if len(parts) < 2 {
		return w, nil
	}

	pairsStr := strings.TrimSpace(parts[1])
	pairs, pairsErr := quoteSplit(pairsStr, ';')
	if pairsErr != nil {
		return Pairs{}, pairsErr
	}

	w.Params = map[string]string{}
	for _, p := range pairs {
		kv, kvErr := quoteSplit(p, '=')
		if kvErr != nil {
			return Pairs{}, kvErr
		}
		if len(kv) != 2 {
			return Pairs{}, fmt.Errorf("invalid pair")
		}
		w.Params[kv[0]] = kv[1]
	}

	return w, nil
}

// "GET /path HTTP/1.1"
// note that version can end with "\r\n" or EOF
var requestLineRE = &lazyre.RE{S: `^(?P<method>[^ ]*[ ]+)(?P<uri>[^ ]*[ ]+)(?P<version>.*?(?:\r\n|$))`}

// "HTTP/1.1 200 OK"
// note that text can end with "\r\n" or EOF
var statusLineRE = &lazyre.RE{S: `^(?P<version>[^ ]*[ ]+)(?P<code>[^ ]*[ ]*)(?P<text>.*?(?:\r\n|$))`}
var headersEndLineRE = &lazyre.RE{S: `^(?P<headers_end>.*\r?\n)`}

// TODO: more methods?
var probePrefixRE = &lazyre.RE{S: `` +
	`^` +
	`(?:` +
	// response
	`HTTP/1` +
	`|` +
	// request
	`(?P<method>` +
	// http methods
	`CONNECT` +
	`|DELETE` +
	`|GET` +
	`|HEAD` +
	`|PATCH` +
	`|POST` +
	`|PUT` +
	`|TRACE` +
	`|OPTIONS` +
	// dav methods
	`|COPY` +
	`|LOCK` +
	`|MKCOL` +
	`|MOVE` +
	`|PROPFIND` +
	`|PROPPATCH` +
	`|UNLOCK` +
	`)` +
	` [[:graph:]]` + // <space><path-etc>
	`)`,
}

func httpDecodeMessage(d *decode.D, isRequest bool) {
	matches := map[string]string{}
	if isRequest {
		d.FieldStruct("request_line", func(d *decode.D) {
			d.FieldRE(requestLineRE.Must(), &matches, scalar.ActualTrimSpace)
		})
	} else {
		d.FieldStruct("status_line", func(d *decode.D) {
			d.FieldRE(statusLineRE.Must(), &matches, scalar.ActualTrimSpace)
		})
	}
	log.Printf("matches: %#+v\n", matches)
	// no body, seems to happen
	if d.End() {
		return
	}

	isHTTPv11 := matches["version"] == "HTTP/1.1"
	isHEAD := matches["method"] == "HEAD"

	_, tpoV := d.FieldFormat("headers", &httpTextprotoGroup, format.TextProto_In{Name: "header"})
	tpo, ok := tpoV.(format.TextProto_Out)
	if !ok {
		panic(fmt.Sprintf("expected TextProtoOut got %#+v", tpoV))
	}
	headers := tpo.Pairs
	d.FieldRE(headersEndLineRE.Must(), nil)

	contentLength := headersFirst(headers, "content-length")
	connection := headersFirst(headers, "connection")
	transferEncoding := headersFirst(headers, "transfer-encoding")
	contentEncoding := headersFirst(headers, "content-encoding")
	contentType := headersFirst(headers, "content-type")

	bodyLen := int64(-1)

	if connection == "Upgrade" {
		upgrade := headersFirst(headers, "upgrade")
		// TODO: h2, h2c
		// TODO: h2c would need HTTP2-Settings from request?
		// h2 => http2 over tls
		// h2c => http2 cleartext
		_ = upgrade

	} else {
		if isHEAD {
			// assume zero content-length for HEAD
			bodyLen = 0
		} else {
			if contentLength != "" {
				if n, err := strconv.ParseInt(contentLength, 10, 64); err == nil {
					bodyLen = n
				}
			} else {
				if isHTTPv11 && connection != "closed" {
					// http 1.1 is persistent by default
					bodyLen = 0
				} else {
					// TODO: assume reset?
				}
			}
		}
	}

	if bodyLen < 0 {
		bodyLen = d.BitsLeft() / 8
	}

	// log.Printf("headers: %#+v\n", headers)

	// log.Printf("contentType: %#+v\n", contentType)
	// TODO: content-range
	// TODO: Transfer-Encoding
	//   chunked + trailer

	// TODO: gzip format hint for subformat?

	contentTypeValues, _ := parsePairs(contentType)

	switch transferEncoding {
	case "chunked":
		d.FieldFormat("body", &httpHttpChunkedGroup, format.Http_Chunked_In{
			ContentEncoding: contentEncoding,
			ContentType:     contentTypeValues.Scheme,
			Pairs:           contentTypeValues.Params,
		})
	default:
		bodyGroup := &httpContentTypeGroup
		bodyGroupInArg := format.Content_Type_In{
			ContentType: contentTypeValues.Scheme,
			Pairs:       contentTypeValues.Params,
		}

		d.FramedFn(bodyLen*8, func(d *decode.D) {
			switch contentEncoding {
			case "gzip":
				if dv, _, _ := d.TryFieldFormat("body", &httpGzipGroup, nil); dv == nil {
					d.FieldRawLen("body", d.BitsLeft())
				}
			default:
				if bodyGroup != nil {
					log.Printf("bodyGroup: %#+v\n", bodyGroup)
					log.Printf("http bodyGroupInArg: %#+v\n", bodyGroupInArg)
					d.FieldFormatOrRawLen("body", d.BitsLeft(), bodyGroup, bodyGroupInArg)

				} else {
					d.FieldRawLen("body", d.BitsLeft())
				}
			}

			// Transfer-Encoding: chunked
			// Transfer-Encoding: compress
			// Transfer-Encoding: deflate
			// Transfer-Encoding: gzip

			// // Several values can be listed, separated by a comma
			// Transfer-Encoding: gzip, chunked

		})
	}
}

func httpDecode(d *decode.D) any {
	var isRequest bool
	var tsi format.TCP_Stream_In

	if d.ArgAs(&tsi) {
		m := d.RE(probePrefixRE.Must())
		if m == nil {
			d.Fatalf("no request or response prefix found")
		}
		isRequest = tsi.IsClient
	} else {
		isRequest = string(d.PeekBytes(5)) != "HTTP/"
	}

	name := "response"
	if isRequest {
		name = "request"
	}
	for !d.End() {
		d.FieldStruct(name, func(d *decode.D) {
			httpDecodeMessage(d, isRequest)
		})
	}

	return nil
}
