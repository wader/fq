module github.com/wader/fq

go 1.17

require (
	// bump: gomod-gopacket /github\.com\/google\/gopacket v(.*)/ https://github.com/google/gopacket.git|^1
	// bump: gomod-gopacket command go get -d github.com/google/gopacket@v$LATEST && go mod tidy
	github.com/google/gopacket v1.1.19
	// bump: gomod-mapstructure /github.com\/mitchellh\/mapstructure v(.*)/ https://github.com/mitchellh/mapstructure.git|^1
	// bump: gomod-mapstructure command go get -d github.com/mitchellh/mapstructure@v$LATEST && go mod tidy
	github.com/mitchellh/mapstructure v1.4.2
	// bump: gomod-go-difflib /github.com\/pmezard\/go-difflib v(.*)/ https://github.com/pmezard/go-difflib.git|^1
	// bump: gomod-go-difflib command go get -d github.com/pmezard/go-difflib@v$LATEST && go mod tidy
	github.com/pmezard/go-difflib v1.0.0
	// bump: gomod-golang/text /golang\.org\/x\/text v(.*)/ https://github.com/golang/text.git|^0
	// bump: gomod-golang/text command go get -d golang.org/x/text@v$LATEST && go mod tidy
	golang.org/x/text v0.3.7
)

require (
	// fork of github.com/itchyny/gojq, see github.com/wader/gojq fq branch
	github.com/wader/gojq v0.12.1-0.20211105163429-4313a117784f
	// fork of github.com/chzyer/readline, see github.com/wader/readline fq branch
	github.com/wader/readline v0.0.0-20210920124728-5a81f7707bac
)

require (
	github.com/itchyny/timefmt-go v0.1.3 // indirect
	golang.org/x/sys v0.0.0-20210831042530-f4d43177bf5e // indirect
)
