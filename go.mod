module github.com/wader/fq

go 1.17

require (
	// bump: gomod-gopacket /github\.com\/google\/gopacket v(.*)/ https://github.com/google/gopacket.git|^1
	// bump: gomod-gopacket command go get -d github.com/google/gopacket@v$LATEST && go mod tidy
	// bump: gomod-gopacket link "Release notes" https://github.com/google/gopacket/releases/tag/v$LATEST
	github.com/google/gopacket v1.1.19
	// bump: gomod-mapstructure /github.com\/mitchellh\/mapstructure v(.*)/ https://github.com/mitchellh/mapstructure.git|^1
	// bump: gomod-mapstructure command go get -d github.com/mitchellh/mapstructure@v$LATEST && go mod tidy
	// bump: gomod-mapstructure link "CHANGELOG" https://github.com/mitchellh/mapstructure/blob/master/CHANGELOG.md
	github.com/mitchellh/mapstructure v1.4.3
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
	// fork of github.com/itchyny/gojq, see github.com/wader/gojq fq branch
	github.com/wader/gojq v0.12.1-0.20211211101122-3894ded312be
	// fork of github.com/chzyer/readline, see github.com/wader/readline fq branch
	github.com/wader/readline v0.0.0-20210920124728-5a81f7707bac
)

require (
	github.com/cilium/ebpf v0.7.0 // indirect
	github.com/cosiner/argv v0.1.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/derekparker/trie v0.0.0-20200317170641-1fdf38b7b0e9 // indirect
	github.com/go-delve/delve v1.8.0 // indirect
	github.com/google/go-dap v0.6.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/itchyny/timefmt-go v0.1.3 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/peterh/liner v1.2.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.3.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.starlark.net v0.0.0-20211203141949-70c0e40ae128 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
