all: test

.PHONY: test
test:
	go test -v -cover -race -coverpkg=./... -coverprofile=cover.out ./...
	go tool cover -func=cover.out

.PHONY: lint
lint:
	golangci-lint run

.PHONY: depgraph.svg
depgraph.svg:
	godepgraph -s cmd/fq/main.go | dot -Tsvg -o godepgraph.svg

.PHONY: formatdgraph.svg
formatdgraph.svg:
	go run cmd/fq/main.go -n - '\
		"digraph formats {", \
		"nodesep=0.5", \
		"ranksep=0.5", \
		"node [shape=\"box\",style=\"rounded,filled\"]", \
		"edge [arrowsize=\"0.7\"]", \
		(formats[] | "\(.name) -> {\(.dependencies | flatten? | join(" "))}"), \
		(formats[] | .name as $$name | .groups[]? | "\(.) -> \($$name)"), \
		(formats | keys[] | "\(.) [color=\"paleturquoise\"]"), \
		([formats[].groups[]?] | unique[] | "\(.) [color=\"palegreen\"]"), \
		"}" \
		' | dot -Tsvg -o formatdgraph.svg
