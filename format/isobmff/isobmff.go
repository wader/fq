package isobmff

import (
	"iter"
	"strings"

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

// TODO: maybe in some future when go supports method type parameters
// these could methods
func findData[T any](n *box, path string) T {
	if n == nil {
		var zero T
		return zero
	}
	if b := n.find(path); b != nil {
		if v, ok := b.data.(T); ok {
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
	allowTruncated bool
	root           *box
	current        *box
	boxCount       int
	ftypSeen       bool
	brands         []string
}

func isobmffDecode(d *decode.D, allowTruncated bool, brands []string) *decodeContext {
	root := &box{typ: ""}
	ctx := &decodeContext{
		allowTruncated: allowTruncated,
		root:           root,
		current:        root,
		ftypSeen:       false,
		brands:         brands,
	}

	decodeBoxes(ctx, d)
	if ctx.boxCount == 0 {
		d.Fatalf("no boxes found")
	}
	if !ctx.ftypSeen {
		d.Errorf("no ftyp box found")
	}

	return ctx
}
