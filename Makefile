all: test

.PHONY: test
test:
	go test -v -cover -race -coverpkg=./... -coverprofile=cover.out ./pkg/format ./pkg/query
	#go tool cover -func=cover.out
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
	godepgraph -s cmd/fq/main.go | dot -Tsvg -o godepgraph.svg

.PHONY: formats.svg
formats.svg:
	go run cmd/fq/main.go -n 'formats | _formats_dot' | dot -Tsvg -o formats.svg

.PHONY: README.md
README.md: _doc/file.mp3 _doc/file.mp4
	$(eval REPODIR=$(shell pwd))
	$(eval TEMPDIR=$(shell mktemp -d))
	cp -a _doc/* "${TEMPDIR}"
	go build -o "${TEMPDIR}/fq" cmd/fq/main.go
	cd "${TEMPDIR}" ; \
	        cat "${REPODIR}/$@" | PATH="${TEMPDIR}:${PATH}" go run "${REPODIR}/_doc/mdsh.go" > "${TEMPDIR}/$@"
	mv "${TEMPDIR}/$@" "${REPODIR}/$@"
	rm -rf "${TEMPDIR}"

_doc/file.mp3: Makefile
	ffmpeg -y -f lavfi -i sine -f lavfi -i testsrc -map 0:0 -map 1:0 -t 20ms "$@"

_doc/file.mp4: Makefile
	ffmpeg -y -f lavfi -i sine -f lavfi -i testsrc -c:a aac -c:v h264 -f mp4 -t 20ms "$@"
