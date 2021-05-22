# bump: golang /FROM golang:([\d.]+)/ docker:golang|^1
FROM golang:1.16.4 AS builder

WORKDIR $GOPATH/src/fq
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile main.go ./
COPY pkg pkg
COPY internal internal
RUN make test fq
RUN cp fq /fq

FROM alpine
COPY --from=builder /fq /fq
RUN ["/fq", "--version"]
ENTRYPOINT ["/fq"]
