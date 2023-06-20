package luajit

import (
	"github.com/wader/fq/pkg/scalar"
)

const (
	BcMnone = iota
	BcMdst
	BcMbase
	BcMvar
	BcMrbase
	BcMuv
	BcMlit
	BcMlits
	BcMpri
	BcMnum
	BcMstr
	BcMtab
	BcMfunc
	BcMjump
	BcMcdata
)

type BcDef struct {
	Name string
	MA   int
	MB   int
	MC   int
}

func (op *BcDef) HasD() bool {
	return op.MB == BcMnone
}

func (op *BcDef) IsJump() bool {
	return op.MC == BcMjump
}

type BcDefList []BcDef

var opcodes = BcDefList{
	{"ISLT", BcMvar, BcMnone, BcMvar},
	{"ISGE", BcMvar, BcMnone, BcMvar},
	{"ISLE", BcMvar, BcMnone, BcMvar},
	{"ISGT", BcMvar, BcMnone, BcMvar},

	{"ISEQV", BcMvar, BcMnone, BcMvar},
	{"ISNEV", BcMvar, BcMnone, BcMvar},
	{"ISEQS", BcMvar, BcMnone, BcMstr},
	{"ISNES", BcMvar, BcMnone, BcMstr},
	{"ISEQN", BcMvar, BcMnone, BcMnum},
	{"ISNEN", BcMvar, BcMnone, BcMnum},
	{"ISEQP", BcMvar, BcMnone, BcMpri},
	{"ISNEP", BcMvar, BcMnone, BcMpri},

	// Unary test and copy ops.
	{"ISTC", BcMdst, BcMnone, BcMvar},
	{"ISFC", BcMdst, BcMnone, BcMvar},
	{"IST", BcMnone, BcMnone, BcMvar},
	{"ISF", BcMnone, BcMnone, BcMvar},
	{"ISTYPE", BcMvar, BcMnone, BcMlit},
	{"ISNUM", BcMvar, BcMnone, BcMlit},

	// Unary ops.
	{"MOV", BcMdst, BcMnone, BcMvar},
	{"NOT", BcMdst, BcMnone, BcMvar},
	{"UNM", BcMdst, BcMnone, BcMvar},
	{"LEN", BcMdst, BcMnone, BcMvar},

	// Binary ops. ORDER OPR. VV last, POW must be next.
	{"ADDVN", BcMdst, BcMvar, BcMnum},
	{"SUBVN", BcMdst, BcMvar, BcMnum},
	{"MULVN", BcMdst, BcMvar, BcMnum},
	{"DIVVN", BcMdst, BcMvar, BcMnum},
	{"MODVN", BcMdst, BcMvar, BcMnum},

	{"ADDNV", BcMdst, BcMvar, BcMnum},
	{"SUBNV", BcMdst, BcMvar, BcMnum},
	{"MULNV", BcMdst, BcMvar, BcMnum},
	{"DIVNV", BcMdst, BcMvar, BcMnum},
	{"MODNV", BcMdst, BcMvar, BcMnum},

	{"ADDVV", BcMdst, BcMvar, BcMvar},
	{"SUBVV", BcMdst, BcMvar, BcMvar},
	{"MULVV", BcMdst, BcMvar, BcMvar},
	{"DIVVV", BcMdst, BcMvar, BcMvar},
	{"MODVV", BcMdst, BcMvar, BcMvar},

	{"POW", BcMdst, BcMvar, BcMvar},
	{"CAT", BcMdst, BcMrbase, BcMrbase},

	// Constant ops.
	{"KSTR", BcMdst, BcMnone, BcMstr},
	{"KCDATA", BcMdst, BcMnone, BcMcdata},
	{"KSHORT", BcMdst, BcMnone, BcMlits},
	{"KNUM", BcMdst, BcMnone, BcMnum},
	{"KPRI", BcMdst, BcMnone, BcMpri},
	{"KNIL", BcMbase, BcMnone, BcMbase},

	// Upvalue and function ops.
	{"UGET", BcMdst, BcMnone, BcMuv},
	{"USETV", BcMuv, BcMnone, BcMvar},
	{"USETS", BcMuv, BcMnone, BcMstr},
	{"USETN", BcMuv, BcMnone, BcMnum},
	{"USETP", BcMuv, BcMnone, BcMpri},
	{"UCLO", BcMrbase, BcMnone, BcMjump},
	{"FNEW", BcMdst, BcMnone, BcMfunc},

	// Table ops.
	{"TNEW", BcMdst, BcMnone, BcMlit},
	{"TDUP", BcMdst, BcMnone, BcMtab},
	{"GGET", BcMdst, BcMnone, BcMstr},
	{"GSET", BcMvar, BcMnone, BcMstr},
	{"TGETV", BcMdst, BcMvar, BcMvar},
	{"TGETS", BcMdst, BcMvar, BcMstr},
	{"TGETB", BcMdst, BcMvar, BcMlit},
	{"TGETR", BcMdst, BcMvar, BcMvar},
	{"TSETV", BcMvar, BcMvar, BcMvar},
	{"TSETS", BcMvar, BcMvar, BcMstr},
	{"TSETB", BcMvar, BcMvar, BcMlit},
	{"TSETM", BcMbase, BcMnone, BcMnum},
	{"TSETR", BcMvar, BcMvar, BcMvar},

	// Calls and vararg handling. T = tail call.
	{"CALLM", BcMbase, BcMlit, BcMlit},
	{"CALL", BcMbase, BcMlit, BcMlit},
	{"CALLMT", BcMbase, BcMnone, BcMlit},
	{"CALLT", BcMbase, BcMnone, BcMlit},
	{"ITERC", BcMbase, BcMlit, BcMlit},
	{"ITERN", BcMbase, BcMlit, BcMlit},
	{"VARG", BcMbase, BcMlit, BcMlit},
	{"ISNEXT", BcMbase, BcMnone, BcMjump},

	// Returns.
	{"RETM", BcMbase, BcMnone, BcMlit},
	{"RET", BcMrbase, BcMnone, BcMlit},
	{"RET0", BcMrbase, BcMnone, BcMlit},
	{"RET1", BcMrbase, BcMnone, BcMlit},

	// Loops and branches. I/J = interp/JIT, I/C/L = init/call/loop.
	{"FORI", BcMbase, BcMnone, BcMjump},
	{"JFORI", BcMbase, BcMnone, BcMjump},

	{"FORL", BcMbase, BcMnone, BcMjump},
	{"IFORL", BcMbase, BcMnone, BcMjump},
	{"JFORL", BcMbase, BcMnone, BcMlit},

	{"ITERL", BcMbase, BcMnone, BcMjump},
	{"IITERL", BcMbase, BcMnone, BcMjump},
	{"JITERL", BcMbase, BcMnone, BcMlit},

	{"LOOP", BcMrbase, BcMnone, BcMjump},
	{"ILOOP", BcMrbase, BcMnone, BcMjump},
	{"JLOOP", BcMrbase, BcMnone, BcMlit},

	{"JMP", BcMrbase, BcMnone, BcMjump},

	// Function headers. I/J = interp/JIT, F/V/C = fixarg/vararg/C func.
	{"FUNCF", BcMrbase, BcMnone, BcMnone},
	{"IFUNCF", BcMrbase, BcMnone, BcMnone},
	{"JFUNCF", BcMrbase, BcMnone, BcMlit},
	{"FUNCV", BcMrbase, BcMnone, BcMnone},
	{"IFUNCV", BcMrbase, BcMnone, BcMnone},
	{"JFUNCV", BcMrbase, BcMnone, BcMlit},
	{"FUNCC", BcMrbase, BcMnone, BcMnone},
	{"FUNCCW", BcMrbase, BcMnone, BcMnone},
}

func (opcodes BcDefList) MapUint(s scalar.Uint) (scalar.Uint, error) {
	listIdx := int(s.Actual)

	if listIdx < len(opcodes) {
		s.Sym = opcodes[listIdx].Name
	}

	return s, nil
}
