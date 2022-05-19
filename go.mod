module github.com/wader/fq

go 1.17

require (
	// fork of github.com/itchyny/gojq, see github.com/wader/gojq fq branch
	github.com/wader/gojq v0.12.1-0.20220503161759-eae34c73de13
	// fork of github.com/chzyer/readline, see github.com/wader/readline fq branch
	github.com/wader/readline v0.0.0-20220519145246-cc6b808370a7
)

require (
	// bump: gomod-golang-snappy /github.com\/golang\/snappy v(.*)/ https://github.com/golang/snappy.git|^0
	// bump: gomod-golang-snappy command go get -d github.com/golang/snappy@v$LATEST && go mod tidy
	// bump: gomod-golang-snappy link "Source diff $CURRENT..$LATEST" https://github.com/golang/snappy/compare/v$CURRENT..v$LATEST
	github.com/golang/snappy v0.0.4
	// bump: gomod-gopacket /github\.com\/google\/gopacket v(.*)/ https://github.com/google/gopacket.git|^1
	// bump: gomod-gopacket command go get -d github.com/google/gopacket@v$LATEST && go mod tidy
	// bump: gomod-gopacket link "Release notes" https://github.com/google/gopacket/releases/tag/v$LATEST
	github.com/google/gopacket v1.1.19
	// bump: gomod-copystructure /github.com\/mitchellh\/copystructure v(.*)/ https://github.com/mitchellh/copystructure.git|^1
	// bump: gomod-copystructure command go get -d github.com/mitchellh/copystructure@v$LATEST && go mod tidy
	// bump: gomod-copystructure link "CHANGELOG" https://github.com/mitchellh/copystructure/blob/master/CHANGELOG.md
	github.com/mitchellh/copystructure v1.2.0
	// bump: gomod-mapstructure /github.com\/mitchellh\/mapstructure v(.*)/ https://github.com/mitchellh/mapstructure.git|^1
	// bump: gomod-mapstructure command go get -d github.com/mitchellh/mapstructure@v$LATEST && go mod tidy
	// bump: gomod-mapstructure link "CHANGELOG" https://github.com/mitchellh/mapstructure/blob/master/CHANGELOG.md
	github.com/mitchellh/mapstructure v1.5.0
	// bump: gomod-go-difflib /github.com\/pmezard\/go-difflib v(.*)/ https://github.com/pmezard/go-difflib.git|^1
	// bump: gomod-go-difflib command go get -d github.com/pmezard/go-difflib@v$LATEST && go mod tidy
	// bump: gomod-go-difflib link "Source diff $CURRENT..$LATEST" https://github.com/pmezard/go-difflib/compare/v$CURRENT..v$LATEST
	github.com/pmezard/go-difflib v1.0.0
	// bump: gomod-golang/text /golang\.org\/x\/text v(.*)/ https://github.com/golang/text.git|^0
	// bump: gomod-golang/text command go get -d golang.org/x/text@v$LATEST && go mod tidy
	// bump: gomod-golang/text link "Source diff $CURRENT..$LATEST" https://github.com/golang/text/compare/v$CURRENT..v$LATEST
	golang.org/x/text v0.3.7
)

require (
	github.com/itchyny/timefmt-go v0.1.3 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
)
