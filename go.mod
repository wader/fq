module github.com/wader/fq

go 1.18

require (
	// fork of github.com/itchyny/gojq, see github.com/wader/gojq fq branch
	github.com/wader/gojq v0.12.1-0.20220816100300-d38cb82d00bf
	// fork of github.com/chzyer/readline, see github.com/wader/readline fq branch
	github.com/wader/readline v0.0.0-20220704090837-31be50517a56
)

require (
	// bump: gomod-BurntSushi/toml /github\.com\/BurntSushi\/toml v(.*)/ https://github.com/BurntSushi/toml.git|^1
	// bump: gomod-BurntSushi/toml command go get -d github.com/BurntSushi/toml@v$LATEST && go mod tidy
	// bump: gomod-BurntSushi/toml link "Source diff $CURRENT..$LATEST" https://github.com/BurntSushi/toml/compare/v$CURRENT..v$LATEST
	github.com/BurntSushi/toml v1.2.0
	// bump: gomod-creasty-defaults /github\.com\/creasty\/defaults v(.*)/ https://github.com/creasty/defaults.git|^1
	// bump: gomod-creasty-defaults command go get -d github.com/creasty/defaults@v$LATEST && go mod tidy
	// bump: gomod-creasty-defaults link "Source diff $CURRENT..$LATEST" https://github.com/creasty/defaults/compare/v$CURRENT..v$LATEST
	github.com/creasty/defaults v1.6.0
	// bump: gomod-golang-snappy /github\.com\/golang\/snappy v(.*)/ https://github.com/golang/snappy.git|^0
	// bump: gomod-golang-snappy command go get -d github.com/golang/snappy@v$LATEST && go mod tidy
	// bump: gomod-golang-snappy link "Source diff $CURRENT..$LATEST" https://github.com/golang/snappy/compare/v$CURRENT..v$LATEST
	github.com/golang/snappy v0.0.4
	// bump: gomod-gopacket /github\.com\/google\/gopacket v(.*)/ https://github.com/google/gopacket.git|^1
	// bump: gomod-gopacket command go get -d github.com/google/gopacket@v$LATEST && go mod tidy
	// bump: gomod-gopacket link "Release notes" https://github.com/google/gopacket/releases/tag/v$LATEST
	github.com/google/gopacket v1.1.19
	// bump: gomod-copystructure /github\.com\/mitchellh\/copystructure v(.*)/ https://github.com/mitchellh/copystructure.git|^1
	// bump: gomod-copystructure command go get -d github.com/mitchellh/copystructure@v$LATEST && go mod tidy
	// bump: gomod-copystructure link "CHANGELOG" https://github.com/mitchellh/copystructure/blob/master/CHANGELOG.md
	github.com/mitchellh/copystructure v1.2.0
	// bump: gomod-mapstructure /github\.com\/mitchellh\/mapstructure v(.*)/ https://github.com/mitchellh/mapstructure.git|^1
	// bump: gomod-mapstructure command go get -d github.com/mitchellh/mapstructure@v$LATEST && go mod tidy
	// bump: gomod-mapstructure link "CHANGELOG" https://github.com/mitchellh/mapstructure/blob/master/CHANGELOG.md
	github.com/mitchellh/mapstructure v1.5.0
	// bump: gomod-go-difflib /github\.com\/pmezard\/go-difflib v(.*)/ https://github.com/pmezard/go-difflib.git|^1
	// bump: gomod-go-difflib command go get -d github.com/pmezard/go-difflib@v$LATEST && go mod tidy
	// bump: gomod-go-difflib link "Source diff $CURRENT..$LATEST" https://github.com/pmezard/go-difflib/compare/v$CURRENT..v$LATEST
	github.com/pmezard/go-difflib v1.0.0
	// has no tags
	// go get -d golang.org/x/crypto@master && go mod tidy
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	// has no tags
	// go get -d golang.org/x/exp@master && go mod tidy
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e
	// has no tags
	// go get -d golang.org/x/net@master && go mod tidy
	golang.org/x/net v0.0.0-20220630215102-69896b714898
	// bump: gomod-golang/text /golang\.org\/x\/text v(.*)/ https://github.com/golang/text.git|^0
	// bump: gomod-golang/text command go get -d golang.org/x/text@v$LATEST && go mod tidy
	// bump: gomod-golang/text link "Source diff $CURRENT..$LATEST" https://github.com/golang/text/compare/v$CURRENT..v$LATEST
	golang.org/x/text v0.3.7
	// bump: gomod-gopkg.in/yaml.v3 /gopkg\.in\/yaml\.v3 v(.*)/ https://github.com/go-yaml/yaml.git|^3
	// bump: gomod-gopkg.in/yaml.v3 command go get -d gopkg.in/yaml.v3@v$LATEST && go mod tidy
	// bump: gomod-gopkg.in/yaml.v3 link "Source diff $CURRENT..$LATEST" https://github.com/go-yaml/yaml/compare/v$CURRENT..v$LATEST
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/itchyny/timefmt-go v0.1.3 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	golang.org/x/sys v0.0.0-20220627191245-f75cf1eec38b // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)
