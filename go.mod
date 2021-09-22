module github.com/wader/fq

go 1.17

require (
	// bump: mapstructure /github.com\/mitchellh\/mapstructure v(.*)/ git://github.com/mitchellh/mapstructure|^1
	// bump: mapstructure command go get -d github.com/mitchellh/mapstructure@v$LATEST && go mod tidy
	github.com/mitchellh/mapstructure v1.4.2
	// bump: go-difflib /github.com\/pmezard\/go-difflib v(.*)/ git://github.com/pmezard/go-difflib|^1
	// bump: go-difflib command go get -d github.com/pmezard/go-difflib@v$LATEST && go mod tidy
	github.com/pmezard/go-difflib v1.0.0

	// fork of github.com/itchyny/gojq, see github.com/wader/gojq fq branch
	github.com/wader/gojq v0.12.1-0.20210922165705-abd9cceb5f74
	// fork of github.com/chzyer/readline, see github.com/wader/readline fq branch
	github.com/wader/readline v0.0.0-20210920124728-5a81f7707bac

	// bump: golang/text /golang\.org\/x\/text v(.*)/ git://github.com/golang/text|^0
	// bump: golang/text command go get -d golang.org/x/text@v$LATEST && go mod tidy
	golang.org/x/text v0.3.7
)

require (
	github.com/itchyny/timefmt-go v0.1.3 // indirect
	golang.org/x/sys v0.0.0-20210831042530-f4d43177bf5e // indirect
)
