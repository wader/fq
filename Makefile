all: test

.PHONY: fq
fq: VERSION=$(shell git describe --all --long --dirty || echo nogit)
fq:
	go build -ldflags "-X main.version=${VERSION}" -trimpath -o fq .

.PHONY: test
# figure out all go pakges with test files
test: PKGS=$(shell find . -name "*_test.go" | xargs -n 1 dirname | sort | uniq)
test: testjq
	go test ${VERBOSE} ${COVER} ${PKGS}

testwrite: export WRITE_ACTUAL=1
testwrite: test

testv: export VERBOSE=-v
testv: test

cover: COVER=-cover -race -coverpkg=./... -coverprofile=cover.out
cover: test
	go tool cover -html=cover.out -o cover.out.html
	cat cover.out.html | grep '<option value="file' | sed -E 's/.*>(.*) \((.*)%\)<.*/\2 \1/' | sort -rn

.PHONY: testjq
testjq:
	@for f in $$(find . -name "*_test.jq"); do \
		echo $$f ; \
		go run main.go -L "$$(dirname $$f)" -f "$$f" -n -r ; \
	done

.PHONY: doc
doc: doc/file.mp3 doc/file.mp4
	$(eval REPODIR=$(shell pwd))
	$(eval TEMPDIR=$(shell mktemp -d))
	@cp -a doc/* "${TEMPDIR}"
	@go build -o "${TEMPDIR}/fq" main.go
	@for f in *.md doc/*.md ; do \
		cd "${TEMPDIR}" ; \
		echo $$f ; \
		mkdir -p $$(dirname "${TEMPDIR}/$$f") ; \
		cat "${REPODIR}/$$f" | PATH="${TEMPDIR}:${PATH}" go run "${REPODIR}/doc/mdsh.go" > "${TEMPDIR}/$$f" ; \
		mv "${TEMPDIR}/$$f" "${REPODIR}/$$f" ; \
	done
	@rm -rf "${TEMPDIR}"

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

.PHONY: formats.svg
formats.svg:
	dev/formats_dot.jq | dot -Tsvg -o formats.svg

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
