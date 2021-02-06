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
README.md:
	PATH=${PWD}/_dev:${PATH} go run _dev/mdsh.go < README.md > README.md.new
	mv README.md.new README.md
