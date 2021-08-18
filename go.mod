module github.com/wader/fq

go 1.16

require (
	// bump: go-difflib /github.com\/pmezard\/go-difflib v(.*)/ git://github.com/pmezard/go-difflib|^1
	// bump: go-difflib command go get -d github.com/pmezard/go-difflib@v$LATEST && go mod tidy
	github.com/pmezard/go-difflib v1.0.0

	// fork of github.com/itchyny/gojq
	github.com/wader/gojq v0.12.1-0.20210818160419-24d2cb60f657
	// fork of github.com/chzyer/readline
	github.com/wader/readline v0.0.0-20210817095433-c868eb04b8b2
)
