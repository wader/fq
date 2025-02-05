# bump: docker-golang /FROM golang:([\d.]+)/ docker:golang|^1
FROM golang:1.23.6-bookworm AS base

# expect is used to test cli
RUN \
    apt-get update -q && \
    apt-get install --no-install-recommends -qy \
    expect

# docker build --target dev -t fq-dev - < Dockerfile && docker run --rm -ti -v "$PWD:/$PWD" -w "$PWD" fq-dev
FROM base AS dev

FROM base AS builder

WORKDIR $GOPATH/src/fq
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile *.go ./
COPY pkg pkg
COPY internal internal
COPY format format
RUN make test fq
RUN cp fq /fq

FROM scratch
COPY --from=builder /fq /fq
RUN ["/fq", "--version"]
ENTRYPOINT ["/fq"]
