.PHONY: vendor
vendor:
	go mod vendor

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