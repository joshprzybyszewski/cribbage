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