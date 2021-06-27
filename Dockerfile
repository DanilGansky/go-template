FROM golang:1.16 AS base
WORKDIR /go/src/template

COPY go.mod /go/src/template
COPY Makefile /go/src/template
RUN make deps

COPY . /go/src/template
RUN make build-release

FROM scratch
WORKDIR /go/src/template
COPY --from=base /go/src/template/bin/server /go/src/template
ENTRYPOINT [ "./server" ]