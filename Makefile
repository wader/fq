all: test

.PHONY: test
test:
	go test -v -cover -race -coverpkg=./... -coverprofile=cover.out ./...
	go tool cover -func=cover.out

.PHONY: generate
generate:
	go generate -x ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: depgraph.svg
depgraph.svg:
	godepgraph -s cmd/fq/main.go | dot -Tsvg -o godepgraph.svg

.PHONY: formatdgraph.svg
formatdgraph.svg:
	go run cmd/fq/main.go -n - 'formats | _formats_dot' | dot -Tsvg -o formatdgraph.svg

.PHONY: README.md
README.md: _doc/file.mp3  _doc/file.mp4
	$(eval REPODIR=$(shell pwd))
	$(eval TEMPDIR=$(shell mktemp -d))
	cp -a _doc/* "${TEMPDIR}"
	go build -o "${TEMPDIR}/fq" cmd/fq/main.go
	cd "${TEMPDIR}" ; \
	        cat "${REPODIR}/$@" | PATH="${TEMPDIR}:${PATH}" go run "${REPODIR}/_doc/mdsh.go" > "${TEMPDIR}/$@"
	mv "${TEMPDIR}/$@" "${REPODIR}/$@"
	rm -rf "${TEMPDIR}"

.PHONY: _doc/file.mp3
_doc/file.mp3:
	ffmpeg -y -f lavfi -i sine -f lavfi -i testsrc -map 0:0 -map 1:0 -t 20ms "$@"

.PHONY: _doc/file.mp4
_doc/file.mp4:
	ffmpeg -y -f lavfi -i sine -f lavfi -i testsrc -c:a aac -c:v h264 -f mp4 -t 20ms "$@"
