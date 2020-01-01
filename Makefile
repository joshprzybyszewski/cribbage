.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

.PHONY: golint
golint:
	golangci-lint run -v ./...

.PHONY: gotest
gotest:
	go test ./...

.PHONY: mongo
mongo:
	# See http://thecodebarbarian.com/introducing-run-rs-zero-config-mongodb-runner
	# sudo npm install --unsafe-perm run-rs -g
	# See https://www.npmjs.com/package/run-rs
	sudo run-rs -v 4.2.1 --shell

.PHONY: install
install:
	./scripts/install.sh
	$(MAKE) vendor