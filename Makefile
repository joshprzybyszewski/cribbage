.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

.PHONY: golint
golint:
	golangci-lint run -v ./...

.PHONY: gotest
gotest:
	go test ./...

.PHONY: install
install:
	./scripts/install.sh
	$(MAKE) vendor

.PHONY: wasm
wasm:
	GOOS=js GOARCH=wasm go build -o assets/wasm/wa_output.wasm github.com/joshprzybyszewski/cribbage/wasm
