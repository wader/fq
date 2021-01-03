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
		["digraph formats {"] + \
		[ \
			[formats[] | . as $$f | .dependencies | flatten? | .[] | [$$f.name, .]] + \
			[formats[] | . as $$f | .groups | flatten? | .[] | [., $$f.name]] | \
			.[] | \
			join(" -> ") \
		] + \
		["}"] | \
		join("\n")' \
		| dot -Tsvg -o formatdgraph.svg
