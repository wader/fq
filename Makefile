all: test

.PHONY: fq
fq:
	go build -buildmode=pie -ldflags '-linkmode external -extldflags "-static-pie"' -trimpath -o fq .

.PHONY: test
test:
	go test -v -cover -race -coverpkg=./... -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.out.html
	cat cover.out.html | grep '<option value="file' | sed -E 's/.*>(.*) \((.*)%\)<.*/\2 \1/' | sort -rn
testwrite: export WRITE_ACTUAL=1
testwrite: test

generate: README.md testwrite

.PHONY: gogenerate
gogenerate:
	go generate -x ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: depgraph.svg
depgraph.svg:
	godepgraph -s main.go | dot -Tsvg -o godepgraph.svg

.PHONY: formats.svg
formats.svg:
	dev/formats_dot.jq | dot -Tsvg -o formats.svg

.PHONY: prof
prof:
	go build -o fq.prof main.go
	CPUPROFILE=fq.cpu.prof MEMPROFILE=fq.mem.prof ./fq.prof "${ARGS}"
.PHONY: memprof
memprof: prof
	go tool pprof -http :5555 fq.prof fq.mem.prof

.PHONY: cpuprof
cpuprof: prof
	go tool pprof -http :5555 fq.prof fq.cpu.prof

.PHONY: update-gomodreplace
update-gomodreplace:
	go mod edit -replace github.com/chzyer/readline=github.com/wader/readline@fq && GOPROXY=direct go mod download github.com/chzyer/readline && go mod tidy
	go mod edit -replace github.com/itchyny/gojq=github.com/wader/gojq@fq && GOPROXY=direct go mod download github.com/itchyny/gojq && go mod tidy


.PHONY: doc
doc: doc/file.mp3 doc/file.mp4
	$(eval REPODIR=$(shell pwd))
	$(eval TEMPDIR=$(shell mktemp -d))
	cp -a doc/* "${TEMPDIR}"
	go build -o "${TEMPDIR}/fq" main.go
	for f in *.md doc/*.md ; do \
		cd "${TEMPDIR}" ; \
		echo $$f ; \
		mkdir -p $$(dirname "${TEMPDIR}/$$f") ; \
		cat "${REPODIR}/$$f" | PATH="${TEMPDIR}:${PATH}" go run "${REPODIR}/doc/mdsh.go" > "${TEMPDIR}/$$f" ; \
		mv "${TEMPDIR}/$$f" "${REPODIR}/$$f" ; \
	done
	rm -rf "${TEMPDIR}"

doc/file.mp3: Makefile
	ffmpeg -y -f lavfi -i sine -f lavfi -i testsrc -map 0:0 -map 1:0 -t 20ms "$@"

doc/file.mp4: Makefile
	ffmpeg -y -f lavfi -i sine -f lavfi -i testsrc -c:a aac -c:v h264 -f mp4 -t 20ms "$@"
