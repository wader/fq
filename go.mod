module fq

go 1.14

require (
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/d5/tengo/v2 v2.6.2
	github.com/itchyny/gojq v0.12.2-0.20210211042401-8acfc7d4b109
	github.com/ozanh/ugo v0.0.0-20201221143933-13b118e1828a
	github.com/pmezard/go-difflib v1.0.0
)

//github.com/chzyer/readline => /Users/wader/src/readline
//replace github.com/itchyny/gojq => /Users/wader/src/gojq
replace github.com/itchyny/gojq => github.com/wader/gojq v0.12.1-0.20210206152702-fef2a4facf55
