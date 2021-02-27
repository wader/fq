module fq

go 1.16

require (
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/itchyny/gojq v0.12.1-0.20210219205417-8d3017ec07d3
	github.com/pmezard/go-difflib v1.0.0
)

//github.com/chzyer/readline => /Users/wader/src/readline
//replace github.com/itchyny/gojq => /Users/wader/src/gojq

replace github.com/itchyny/gojq => github.com/wader/gojq v0.12.1-0.20210227093348-ffb6d053b8bc
