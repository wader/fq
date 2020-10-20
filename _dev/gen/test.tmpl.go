package main

{{- range $s := array (map "short" "S" "desc" "signed") (map "short" "U" "desc" "unsigned")}}

{{- range $e := array (map "short" "" "long" "BigEndian" "desc" "big endian") (map "short" "BE" "long" "BigEndian" "desc" "big endian")  (map "short" "LE" "long" "LittleEndian" "desc" "little endian")}}
// S{{$e.short}} read a nBits {{$s.desc}} {{$e.desc}} integer
func (c *Common) {{$s.short}}{{$e.short}}(nBits int64) int64 { return c.SE(nBits, bitbuf.{{$e.long}}) }
{{- range $i := xrange 1 64}}
// FieldS{{$e.short}}{{$i}} read {{$i}} bit {{$s.desc}} {{$e.desc}} integer
func (c *Common) {{$s.short}}{{$e.short}}{{$i}}() int64 { return c.SE({{$i}}, bitbuf.{{$e.long}}) }
{{- end}}

// Field{{$s.short}}{{$e.short}} read a nBits {{$s}} {{$e.desc}} integer and add a field
func (c *Common) FieldS{{$e.short}}(name string, nBits int64) int64 { return c.FieldSE(name, nBits, bitbuf.{{$e.long}}) }
{{- range $i := xrange 1 64}}
// Field{{$s.short}}{{$e.short}}{{$i}} read {{$i}} bit {{$s.desc}} {{$e.desc}} integer and add a field
func (c *Common) FieldS{{$e.short}}{{$i}}(name string) int64 { return c.FieldSE(name, {{$i}}, bitbuf.{{$e.long}}) }
{{- end}}
{{- end}}
{{- end}}
