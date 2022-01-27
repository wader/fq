%{

// https://github.com/kaitai-io/kaitai_struct_compiler/blob/master/shared/src/main/scala/io/kaitai/struct/exprlang

package ksexpr

%}

%union{
    token Token
    node Node
    nodes []Node
    ns []Token
}

%type <node> expr
%type <node> term
%type <node> trailer
%type <nodes> trailers
%type <nodes> arglist
%type <node> array
%type <nodes> arraylist
%type <ns> ns

%token <token> tokNumber
%token <token> tokIdent
%token <token> tokString
%token <token> tokLessEq
%token <token> tokGreaterEq
%token <token> tokEqEq
%token <token> tokNotEq

%token <token> tokBSL
%token <token> tokBSR
%token <token> tokBAnd
%token <token> tokBOr
%token <token> tokBXor
%token <token> tokNot
%token <token> tokAnd
%token <token> tokOr
%token <token> tokTrue
%token <token> tokFalse
%token <token> tokColonColon
%token <token> tokUnterminatedString
%token <token> tokError

%left tokNot
%left '~'
%left '|'
%left '^'
%right '?' ':'
%left tokBSL tokBSR tokBXor tokAnd tokOr
%nonassoc '<' tokLessEq '>' tokGreaterEq tokEqEq tokNotEq
%left '&'
%left '+' '-'
%left '*' '/' '%'

%%

start:
    expr { yylex.(*yyLex).result = $1 }

expr:
	  expr '+' expr          { $$ = InfixOpNode{LHS: $1, Op: InfixOpAdd, RHS: $3} }
	| expr '-' expr          { $$ = InfixOpNode{LHS: $1, Op: InfixOpSub, RHS: $3} }
	| expr '/' expr          { $$ = InfixOpNode{LHS: $1, Op: InfixOpDiv, RHS: $3} }
	| expr '*' expr          { $$ = InfixOpNode{LHS: $1, Op: InfixOpMul, RHS: $3} }
	| expr '%' expr          { $$ = InfixOpNode{LHS: $1, Op: InfixOpMod, RHS: $3} }
	| expr '<' expr          { $$ = InfixOpNode{LHS: $1, Op: InfixOpLT, RHS: $3} }
	| expr tokLessEq expr    { $$ = InfixOpNode{LHS: $1, Op: InfixOpLTEQ, RHS: $3} }
	| expr '>' expr          { $$ = InfixOpNode{LHS: $1, Op: InfixOpGT, RHS: $3} }
	| expr tokGreaterEq expr { $$ = InfixOpNode{LHS: $1, Op: InfixOpGTEQ, RHS: $3} }
	| expr tokEqEq expr      { $$ = InfixOpNode{LHS: $1, Op: InfixOpEQ, RHS: $3} }
	| expr tokNotEq expr     { $$ = InfixOpNode{LHS: $1, Op: InfixOpNotEQ, RHS: $3} }
	| expr tokBSL expr       { $$ = InfixOpNode{LHS: $1, Op: InfixOpBSL, RHS: $3} }
	| expr tokBSR expr       { $$ = InfixOpNode{LHS: $1, Op: InfixOpBSR, RHS: $3} }
	| expr '&' expr          { $$ = InfixOpNode{LHS: $1, Op: InfixOpBAnd, RHS: $3} }
	| expr '|' expr          { $$ = InfixOpNode{LHS: $1, Op: InfixOpBOr, RHS: $3} }
  | expr '^' expr          { $$ = InfixOpNode{LHS: $1, Op: InfixOpBXor, RHS: $3} }
  | expr tokAnd expr       { $$ = InfixOpNode{LHS: $1, Op: InfixOpAnd, RHS: $3} }
  | expr tokOr expr        { $$ = InfixOpNode{LHS: $1, Op: InfixOpOr, RHS: $3} }
	| '~' expr               { $$ = PrefixOpNode{Expr: $2, Op: PrefixOpInv} }
	| '-' expr               { $$ = PrefixOpNode{Expr: $2, Op: PrefixOpNeg} }
	| tokNot expr            { $$ = PrefixOpNode{Expr: $2, Op: PrefixOpNot} }
  | expr '?' expr ':' expr { $$ = TernaryNode{Expr: $1, TrueExpr: $3, FalseExpr: $5} }
  | term                   { $$ = $1 }
  | term trailers {
      t, _ := $1.(TermNode)
      t.Trailers = $2
      $$ = t
  }

trailers:
      trailer          { $$ = []Node{$1} }
    | trailers trailer { $$ = append($1, $2) }
trailer:
      '.' tokIdent '(' arglist ')' { $$ = TrailerCallNode{Name: $2, Args: $4} }
    | '.' tokIdent                 { $$ = TrailerCallNode{Name: $2} }
    | '[' expr ']'                 { $$ = TrailerIndexNode{Expr: $2} }

arglist:
      expr             { $$ = []Node{$1} }
    | arglist ',' expr { $$ = append($1, $3) }

term:
      tokTrue                         { $$ = TermNode{T: ConstNode($1)} }
    | tokFalse                        { $$ = TermNode{T: ConstNode($1)} }
    | '(' expr ')'                    { $$ = TermNode{T: $2} }
    | tokNumber                       { $$ = TermNode{T: ConstNode($1)} }
    | tokIdent                        { $$ = TermNode{T: IdentNode{Name:$1}} }
    | ns tokColonColon tokIdent       { $$ = TermNode{T: IdentNode{NS: $1, Name: $3}} }
    | tokString                       { $$ = TermNode{T: ConstNode($1)} }
    | array                           { $$ = TermNode{T: $1 } }

ns:
      tokIdent                   { $$ = []Token{$1} }
    | ns tokColonColon tokIdent  { $$ = append($1, $3) }

array:
      '[' ']'           { $$ = ArrayNode{} }
    | '[' arraylist ']' { $$ = ArrayNode($2) }
arraylist:
      expr               { $$ = []Node{$1} }
    | arraylist ',' expr { $$ = append($1, $3) }

%%