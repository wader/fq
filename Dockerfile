# bump: golang /FROM golang:([\d.]+)/ docker:golang|^1
FROM golang:1.16.0 AS builder

WORKDIR $GOPATH/src/fq
COPY go.mod go.sum ./
RUN go mod download
COPY version.go .
COPY pkg pkg
COPY internal internal
COPY cmd cmd
RUN go test -v -cover -race ./pkg/format ./pkg/query
RUN CGO_ENABLED=0 go build -o /fq -ldflags '-extldflags "-static"' ./cmd/fq

FROM scratch
COPY --from=builder /fq /fq
RUN ["/fq", "-version"]
ENTRYPOINT ["/fq"]
