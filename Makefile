all: test

.PHONY: fq
fq:
	CGO_ENABLED=0 go build -o fq -ldflags '-extldflags "-static"' .

.PHONY: test
test:
	go test -v -cover -race -coverpkg=./... -coverprofile=cover.out ./pkg/format ./pkg/interp
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
	go run main.go -rn 'formats | _formats_dot' | dot -Tsvg -o formats.svg

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

.PHONY: README.md
README.md: _doc/file.mp3 _doc/file.mp4
	$(eval REPODIR=$(shell pwd))
	$(eval TEMPDIR=$(shell mktemp -d))
	cp -a _doc/* "${TEMPDIR}"
	go build -o "${TEMPDIR}/fq" main.go
	cd "${TEMPDIR}" ; \
	        cat "${REPODIR}/$@" | PATH="${TEMPDIR}:${PATH}" go run "${REPODIR}/_doc/mdsh.go" > "${TEMPDIR}/$@"
	mv "${TEMPDIR}/$@" "${REPODIR}/$@"
	rm -rf "${TEMPDIR}"

_doc/file.mp3: Makefile
	ffmpeg -y -f lavfi -i sine -f lavfi -i testsrc -map 0:0 -map 1:0 -t 20ms "$@"

_doc/file.mp4: Makefile
	ffmpeg -y -f lavfi -i sine -f lavfi -i testsrc -c:a aac -c:v h264 -f mp4 -t 20ms "$@"
