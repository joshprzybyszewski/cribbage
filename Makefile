.PHONY: vendor
vendor:
	go mod vendor

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
	echo "spinning up mongo in replica mode"
	sudo run-rs -v 4.2.1 --shell
	# else you should do mongod

.PHONY: install
install:
	./scripts/install.sh
	$(MAKE) vendor

.PHONY: client
client:
	go run localclient/main/main.go

.PHONY: dockerbuild
dockerbuild:
	docker build -t cribbage .