module fq

go 1.16

require (
	github.com/goinsane/readline v1.5.0
	github.com/itchyny/gojq v0.12.1-0.20210219205417-8d3017ec07d3

	// bump: go-difflib /github.com\/pmezard\/go-difflib v(.*)/ git://github.com/pmezard/go-difflib|^1
	// bump: go-difflib command go get -d github.com/pmezard/go-difflib@v$LATEST && go mod tidy
	github.com/pmezard/go-difflib v1.0.0
)

// go mod edit -replace github.com/goinsane/readline=github.com/wader/readline@develop-v2 && GOPROXY=direct go mod download github.com/goinsane/readline && go mod tidy
replace github.com/goinsane/readline => github.com/wader/readline v0.0.0-20210306181459-854482684b51

// go mod edit -replace github.com/itchyny/gojq=github.com/wader/gojq@fq && GOPROXY=direct go mod download github.com/itchyny/gojq && go mod tidy
replace github.com/itchyny/gojq => github.com/wader/gojq v0.12.1-0.20210615152814-6b798cc8396d

// replace github.com/itchyny/gojq => /Users/wader/src/gojq
