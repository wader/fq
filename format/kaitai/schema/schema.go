//nolint:tagliatelle
package schema

// TODO: error line number? in error via Node?

import (
	"fmt"
	"io"

	"github.com/wader/fq/format/kaitai/ksexpr"
	"github.com/wader/fq/format/kaitai/schema/primitive"
	"gopkg.in/yaml.v3"
)

type ValueError struct {
	Node *yaml.Node
	Err  error
}

func (v ValueError) Unwrap() error { return v.Err }

func (v ValueError) Error() string {
	return fmt.Sprintf("%d: %s", v.Node.Line, v.Err)
}

func valueErrorf(n *yaml.Node, format string, a ...any) ValueError {
	return ValueError{
		Node: n,
		Err:  fmt.Errorf(format, a...),
	}
}

// meta:
//   id: ...
//   endian: le | be
//   bit-endian: le | be              # for b#(|be|le)
//   endian:
//     switch-on:
//     cases:
//       <expr>: le | be              # can use _ for default
//
// seq:
//   - id: <name>
//
//     type: u1 | str | <type> | ...  # builtin scalar type or <type> (seq)
//     type:
//       switch-in <expr>             # switch on cases using expr
//       cases:
//          <expr>: <type>            # can use _ for default
//
//     size: <expr>                      # size many bytes (even for UTF-16 etc)
//     size-eos: <bool>                  # true is size if rest of stream (not an expression)
//     contents: string | [1,0x1,string] # read content size of bytes and validate (constant)
//
//     enum: <string>                 # enum mapping (not an expression)
//
//     if: <expr>                     # should be skipped?
//
//     repeat: expr | until           # id is an array of type
//     repeat-expr: <expr>            # number of times
//     repeat-until: <expr>           # until expr is true
//
// types:
//   <type>:
//     seq:
//       - id: <name>
//         ...
//
// enums:
//   <name>:
//     0: <name>
//     0:
//       - id: name>
//         doc: <doc>
//
// instances:
//   <id>:
//      value: <expr>
//   <id>:
//      pos: <expr
//      type: <string>  # TODO: switch-on?

type Expr struct {
	Str    string
	KSExpr ksexpr.Node
}

func (e *Expr) UnmarshalYAML(value *yaml.Node) error {
	if err := value.Decode(&e.Str); err != nil {
		return err
	}

	ke, err := ksexpr.Parse(e.Str)
	if err != nil {
		return fmt.Errorf("failed to parse '%s': %w", e.Str, err)
	}
	e.KSExpr = ke

	return nil
}

type Endian primitive.Endianess

func (e *Endian) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	switch s {
	case "le":
		(*e) = Endian(primitive.LE)
		return nil
	case "be":
		(*e) = Endian(primitive.BE)
		return nil
	default:
		return valueErrorf(value, "unknown endian %q", s)
	}
}

type CaseExpr[T any] struct {
	Expr  *Expr
	Value T
}

type ValueOrSwitch[T any] struct {
	Value    *T
	SwitchOn *Expr        `yaml:"switch-on"`
	Cases    map[string]T `yaml:"cases"`

	CasesExprs []CaseExpr[T]
}

func (t *ValueOrSwitch[T]) UnmarshalYAML(value *yaml.Node) error {
	type tt ValueOrSwitch[T]
	var tv tt

	if err := value.Decode(&t.Value); err == nil {
		// field: string
		return nil
	} else if err := value.Decode(&tv); err == nil {
		// field:
		//   switch-in: <string>
		//   cases:
		//     <string>: <string>
		//     ...
		(*t) = ValueOrSwitch[T](tv)

		for es, v := range t.Cases {
			e, err := ksexpr.Parse(es)
			if err != nil {
				return fmt.Errorf("failed to parse case expr: %w", err)
			}
			t.CasesExprs = append(t.CasesExprs, CaseExpr[T]{
				Expr:  &Expr{KSExpr: e, Str: es},
				Value: v,
			})
		}

		return nil
	}

	return fmt.Errorf("failed to parse as value or switch-on")
}

type Meta struct {
	ID        string  `yaml:"id"`
	Title     string  `yaml:"title"`
	Endian    *Endian `yaml:"endian"`
	BitEndian *Endian `yaml:"bit-endian"`
}

type Enum struct {
	ToID   map[any]string
	FromID map[string]any
}

type EnumEntry struct {
	ID string `yaml:"id"`
}

func (e *EnumEntry) UnmarshalYAML(value *yaml.Node) error {
	type et EnumEntry
	var ev et

	if err := value.Decode(&e.ID); err == nil {
		return nil
	} else if err := value.Decode(&ev); err == nil {
		// TODO: fix this
		(*e) = EnumEntry(ev)
		return nil
	}
	return fmt.Errorf("failed to parse enum entry as string or id/doc")
}

func (e *Enum) UnmarshalYAML(value *yaml.Node) error {
	var em map[string]EnumEntry
	if err := value.Decode(&em); err != nil {
		return err
	}

	e.ToID = map[any]string{}
	e.FromID = map[string]any{}

	for es, v := range em {
		en, err := ksexpr.Parse(es)
		if err != nil {
			return fmt.Errorf("failed to parse enum expr: %w", err)
		}
		ev, err := en.Eval(0)
		if err != nil {
			return fmt.Errorf("failed to eval enum expr: %w", err)
		}

		// TODO: ev should be a number

		// TODO: types
		e.ToID[ev] = v.ID
		e.FromID[v.ID] = ev
	}

	return nil
}

type Type struct {
	Meta *Meta                `yaml:"meta"`
	ID   string               `yaml:"id"`
	Type *ValueOrSwitch[Expr] `yaml:"type"`

	Size     *Expr `yaml:"size"`
	SizeEOS  bool  `yaml:"size-eos"`
	Contents any   `yaml:"contents"`

	Repeat      string `yaml:"repeat"`
	RepeatExpr  *Expr  `yaml:"repeat-expr"`
	RepeatUntil *Expr  `yaml:"repeat-until"`

	Enum *Expr `yaml:"enum"`

	If *Expr `yaml:"if"`

	// only instance
	Value *Expr `yaml:"value"`
	Pos   *Expr `yaml:"pos"`

	Seq       []*Type          `yaml:"seq"`
	Types     map[string]*Type `yaml:"types"`
	Enums     map[string]*Enum `yaml:"enums"`
	Instances map[string]*Type `yaml:"instances"`

	// used by primitive types like u1, str etc
	Primitive *primitive.Type `yaml:"-"`

	Root   *Type `yaml:"-"` // TODO: not needed?
	Parent *Type `yaml:"-"`
}

func (t *Type) resolveTypeUp(s string) (*Type, bool) {
	// 1. check types
	if tt, ok := t.Types[s]; ok {
		return tt, true
	}
	// 2. check if current type
	if t.ID == s {
		return t, true
	}
	// 3. look in parent
	if t.Parent != nil {
		return t.Parent.resolveTypeUp(s)
	}
	return nil, false
}

// type and enum resolve algorithm described here:
// https://github.com/kaitai-io/kaitai_struct/issues/1019#issuecomment-1503769699
// https://github.com/kaitai-io/kaitai_struct_doc/blob/88c32183b26e4d2265bfecdce0a3dfeaea6975ff/ksy_reference.adoc#type
// https://github.com/kaitai-io/kaitai_struct_compiler/blob/829a14f1e33e8e48eeae726c8a287a5967bcb668/shared/src/main/scala/io/kaitai/struct/ClassTypeProvider.scala
func (t *Type) ResolveType(ns []string) (*Type, bool) {
	lns := len(ns)
	if lns == 0 {
		return nil, false
	} else if lns == 1 {
		if pt, ok := primitive.Types[ns[0]]; ok {
			return &Type{
				ID:        ns[0],
				Primitive: pt,
			}, true
		}
	}

	tt, ok := t.resolveTypeUp(ns[0])
	if !ok {
		return nil, false
	}

	// now resolve ns[1:] parts (could be none) in found type
	for _, n := range ns[1:] {
		var ok bool
		tt, ok = tt.Types[n]
		if !ok {
			return nil, false
		}
	}

	return tt, true
}

func (t *Type) resolveEnumUp(s string) (*Enum, bool) {
	// 1. check enums
	if te, ok := t.Enums[s]; ok {
		return te, true
	}
	// 2. look in parent
	if t.Parent != nil {
		return t.Parent.resolveEnumUp(s)
	}
	return nil, false
}

func (t *Type) ResolveEnum(ns []string) (*Enum, bool) {
	lns := len(ns)
	if lns == 0 {
		return nil, false
	} else if lns == 1 {
		// enum: a
		return t.resolveEnumUp(ns[0])
	}

	// enum (t::)+a
	tt, ok := t.ResolveType(ns[:lns-1])
	if !ok {
		return nil, false
	}
	e, ok := tt.Enums[ns[lns-1]]
	if !ok {
		return nil, false
	}

	return e, true
}

func (t *Type) assignParentAndId(id string, root, parent *Type) {
	t.ID = id
	t.Root = root
	t.Parent = parent
	for _, ct := range t.Seq {
		ct.assignParentAndId(ct.ID, root, t)
	}
	for id, ct := range t.Types {
		ct.assignParentAndId(id, root, t)
	}
	for id, ct := range t.Instances {
		ct.assignParentAndId(id, root, t)
	}
}

func Parse(r io.Reader) (*Type, error) {
	t := &Type{}
	err := yaml.NewDecoder(r).Decode(t)
	if err != nil {
		return nil, err
	}

	// set parent and type id:s
	t.assignParentAndId(t.Meta.ID, t, nil)

	return t, nil
}
