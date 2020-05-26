FROM golang:1.14.3-alpine3.11 as build

ENV GOPATH=/go
RUN mkdir -p $GOPATH/src/github.com/joshprzybyszewski/cribbage
WORKDIR $GOPATH/src/github.com/joshprzybyszewski/cribbage

COPY vendor vendor
COPY model model
COPY logic logic
COPY utils utils
COPY jsonutils jsonutils
COPY network network
COPY server server
COPY wasm wasm
COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -o assets/wasm/wa_output.wasm github.com/joshprzybyszewski/cribbage/wasm
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/cribbageServer main.go

FROM scratch

WORKDIR /prod
COPY templates templates
COPY assets assets
COPY --from=build /go/src/github.com/joshprzybyszewski/cribbage/assets/wasm/wa_output.wasm assets/wasm/wa_output.wasm
COPY --from=build /bin/cribbageServer .
ENTRYPOINT ["/prod/cribbageServer"]
# We're gonna need to read these from an INI or something instead of trying to pass them in as flags
CMD ["-restPort=8081", "-dsn_host=host.docker.internal", "-dsn_user=root", "-dsn_password="]