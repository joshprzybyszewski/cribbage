FROM golang:1.14.3-alpine3.11 as build

ENV GOPATH=/go
WORKDIR $GOPATH/src/github.com/joshprzybyszewski/cribbage

# vendor'ed dependencies are unlikely to change, so download them first
COPY go.mod go.sum ./
RUN go mod download

# Copy the specific directories/files we need to build our binaries
COPY utils utils
COPY model model
COPY network network
COPY logic logic
COPY jsonutils jsonutils
COPY server server
COPY wasm wasm
COPY main.go main.go

# Build our golang binaries (client wasm output and the gin server binary)
RUN CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -o /bin/wa_output.wasm github.com/joshprzybyszewski/cribbage/wasm
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/cribbageServer main.go

# Start a new image that only holds the bare minimum files so that we don't build too much into our image
FROM scratch

WORKDIR /prod
COPY templates templates
COPY assets assets
COPY --from=build /bin/wa_output.wasm assets/wasm/wa_output.wasm
COPY --from=build /bin/cribbageServer .

# Define the gin server binary as the entry point
ENTRYPOINT ["/prod/cribbageServer"]
# We're gonna need to read these from an INI or something instead of trying to pass them in as flags
CMD ["-restPort=8081", "-dsn_host=host.docker.internal", "-dsn_user=root", "-dsn_password="]