package isobmff

import (
	"iter"
	"slices"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
)

type box struct {
	typ      string
	data     any
	parent   *box
	children []*box
}

// findAll traverses the box tree using a "/"-separated path syntax.
// Only the last component matches are yielded:
//
//	"name"     — immediate child named name
//	"<name"    — parent named name
//	"<<name"   — nearest ancestor named name (walks up)
//	">>name"   — any descendant named name (walks all children recursively)
//
//	findAll("moov/trak")  — all trak under moov
//	find("<moof")         — parent moof
//	find("<<traf/tfhd")   — nearest traf ancestor, then its child tfhd
//	find("<<trak/>>tenc") — nearest trak ancestor, then descendant tenc
//	find("<<stbl/stsd")   — nearest stbl ancestor, then child stsd
func (n *box) findAll(path string) iter.Seq[*box] {
	return func(yield func(*box) bool) {
		next := n
		parts := strings.Split(path, "/")
		for i, p := range parts {
			isLast := i == len(parts)-1
			switch {
			case strings.HasPrefix(p, "<<"):
				for a := next.parent; a != nil; a = a.parent {
					if a.typ == p[2:] {
						if isLast {
							if !yield(a) {
								return
							}
						} else {
							next = a
						}
						break
					}
				}
			case strings.HasPrefix(p, "<"):
				if next.parent != nil && next.parent.typ == p[1:] {
					if isLast {
						if !yield(next.parent) {
							return
						}
					} else {
						next = next.parent
					}
				}
			case strings.HasPrefix(p, ">>"):
				var walk func(*box) bool
				walk = func(n *box) bool {
					for _, child := range n.children {
						if child.typ == p[2:] {
							if isLast {
								if !yield(child) {
									return false
								}
							} else {
								next = child
								break
							}
						}
						if !walk(child) {
							return false
						}
					}
					return true
				}
				if !walk(next) {
					return
				}
			default:
				for _, child := range next.children {
					if child.typ == p {
						if isLast {
							if !yield(child) {
								return
							}
						} else {
							next = child
						}
					}
				}
			}
			if isLast {
				return
			}
		}
	}
}

func (n *box) find(path string) *box {
	for m := range n.findAll(path) {
		return m
	}
	return nil
}

func findData[T any](n *box, path string) T {
	if n == nil {
		var zero T
		return zero
	}
	if node := n.find(path); node != nil {
		if v, ok := node.data.(T); ok {
			return v
		}
	}
	var zero T
	return zero
}

func findAllData[T any](n *box, path string) iter.Seq[T] {
	return func(yield func(T) bool) {
		if n == nil {
			return
		}
		for m := range n.findAll(path) {
			if v, ok := m.data.(T); ok {
				if !yield(v) {
					return
				}
			}
		}
	}
}

type decodeContext struct {
	opts    format.MP4_In
	root    *box
	current *box
}

func isobmffDecode(d *decode.D, brandsFn func(firstType string, ftyp ftypBox)) any {
	var mi format.MP4_In
	d.ArgAs(&mi)

	root := &box{typ: ""}
	ctx := &decodeContext{
		opts:    mi,
		root:    root,
		current: root,
	}

	// TODO: nicer, validate functions without field?
	d.AssertLeastBytesLeft(16)
	size := d.U32()
	if size < 8 {
		d.Fatalf("first box size too small < 8")
	}
	var ftyp ftypBox
	firstType := strings.TrimSpace(d.UTF8(4))
	// this is to make it possible to force decode when the first box is not ftyp or styp
	switch firstType {
	case "ftyp", "styp":
		ftyp.majorBrand = strings.TrimSpace(d.UTF8(4))
		minorCount := (size - (4 * 4)) / 4 /* size,type,major,minor_version */
		ftyp.minorVersion = uint32(d.U32())
		for i := 0; i < int(minorCount); i++ {
			ftyp.minorBrands = append(ftyp.minorBrands, strings.TrimSpace(d.UTF8(4)))
		}
	}

	brandsFn(firstType, ftyp)
	d.SeekAbs(0)
	decodeBoxes(ctx, d)

	trakNodes := slices.Collect(ctx.root.findAll("moov/trak"))
	moofNodes := slices.Collect(ctx.root.findAll("moof"))
	if len(trakNodes) > 0 || len(moofNodes) > 0 {
		mp4Tracks(d, ctx, trakNodes, moofNodes)
	}

	return nil
}
