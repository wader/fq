GO_BUILD_FLAGS=-trimpath
GO_BUILD_LDFLAGS=-s -w

all: test fq

.PHONY: fq
fq: VERSION=$(shell git describe --all --long --dirty 2>/dev/null || echo nogit)
fq:
	CGO_ENABLED=0 go build -o fq -ldflags "${GO_BUILD_LDFLAGS} -X main.version=${VERSION}" ${GO_BUILD_FLAGS} .

.PHONY: test
test: testgo testjq testcli

.PHONY: testgo
# figure out all go pakges with test files
testgo: PKGS=$(shell find . -name "*_test.go" | xargs -n 1 dirname | sort | uniq)
testgo:
	go test ${VERBOSE} ${COVER} ${PKGS}

.PHONY: testgov
testgov: export VERBOSE=-v
testgov: testgo

.PHONY: testjq
testjq: fq
	@dev/testjq.sh ./fq pkg/interp/*_test.jq

.PHONY: testcli
testcli: fq
	@pkg/cli/test.sh ./fq pkg/cli/test.exp

.PHONY: actual
actual: export WRITE_ACTUAL=1
actual: testgo

.PHONY: cover
cover: COVER=-cover -race -coverpkg=./... -coverprofile=cover.out
cover: test
	go tool cover -html=cover.out -o cover.out.html
	cat cover.out.html | grep '<option value="file' | sed -E 's/.*>(.*) \((.*)%\)<.*/\2 \1/' | sort -rn

.PHONY: doc
doc: fq doc/file.mp3 doc/file.mp4 doc/formats.svg doc/demo.svg
	@doc/mdsh.sh ./fq *.md doc/*.md

.PHONY: doc/demo.svg
doc/demo.svg: fq
	(cd doc ; ./demo.sh ../fq) | go run github.com/wader/ansisvg@master > doc/demo.svg

.PHONY: doc/formats.svg
doc/formats.svg: fq
	./fq -rnf doc/formats_diagram.jq | dot -Tsvg -o doc/formats.svg

doc/file.mp3: Makefile
	ffmpeg -y -f lavfi -i sine -f lavfi -i testsrc -map 0:0 -map 1:0 -t 20ms "$@"

doc/file.mp4: Makefile
	ffmpeg -y -f lavfi -i sine -f lavfi -i testsrc -c:a aac -c:v h264 -f mp4 -t 20ms "$@"

.PHONY: gogenerate
gogenerate:
	go generate -x ./...

.PHONY: lint
lint:
# bump: make-golangci-lint /golangci-lint@v([\d.]+)/ git:https://github.com/golangci/golangci-lint.git|^1
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1 run

.PHONY: depgraph.svg
depgraph.svg:
	go run github.com/kisielk/godepgraph@latest github.com/wader/fq | dot -Tsvg -o godepgraph.svg

# make memprof ARGS=". test.mp3"
# make cpuprof ARGS=". test.mp3"
.PHONY: prof
prof:
	go build -tags profile -o fq.prof main.go
	CPUPROFILE=fq.cpu.prof MEMPROFILE=fq.mem.prof ./fq.prof ${ARGS}
.PHONY: memprof
memprof: prof
	go tool pprof -http :5555 fq.prof fq.mem.prof

.PHONY: cpuprof
cpuprof: prof
	go tool pprof -http :5555 fq.prof fq.cpu.prof

.PHONY: update-gomodreplace
update-gomod:
	GOPROXY=direct go get -d github.com/wader/readline@fq
	GOPROXY=direct go get -d github.com/wader/gojq@fq
	go mod tidy
