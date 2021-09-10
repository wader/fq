# bump: golang /FROM golang:([\d.]+)/ docker:golang|^1
FROM golang:1.17.1 AS base

# docker build --target dev -t fq-dev - < Dockerfile && docker run --rm -ti -v "$PWD:/$PWD" -w "$PWD" fq-dev
FROM base AS dev

# bump: golangci-lint /GOLANGCILINT_VERSION=([\d.]+)/ git:https://github.com/golangci/golangci-lint.git|^1
ARG GOLANGCILINT_VERSION=1.42.1
RUN \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b /usr/local/bin v$GOLANGCILINT_VERSION

FROM base AS builder

WORKDIR $GOPATH/src/fq
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile main.go ./
COPY pkg pkg
COPY internal internal
RUN make test fq
RUN cp fq /fq

FROM scratch
COPY --from=builder /fq /fq
RUN ["/fq", "--version"]
ENTRYPOINT ["/fq"]
