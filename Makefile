.PHONY: help
help: ## Prints out the options available in this makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: vendor
vendor: ## Gets vendored depenencies for golang
	go mod vendor

.PHONY: golint
golint: ## Runs linters (via golangci-lint) on golang code
	golangci-lint run -v ./...

.PHONY: gotest
gotest: ## Runs all of the golang unit tests
	go test ./...

.PHONY: mongo
mongo: ## Sets up the mongo database in replica mode
	# See http://thecodebarbarian.com/introducing-run-rs-zero-config-mongodb-runner
	# sudo npm install --unsafe-perm run-rs -g
	# See https://www.npmjs.com/package/run-rs
	echo "spinning up mongo in replica mode"
	sudo run-rs -v 4.2.1 --shell
	# else you should do mongod

.PHONY: install
install: ## Runs the install script and vendors golang dependencies
	./scripts/install.sh
	$(MAKE) vendor

.PHONY: goclient
goclient: ## Runs the old golang survey client to play cribbage
	go run localclient/main/main.go

tailwind: client/src/styles.css

client/src/styles.css: client/src/tailwind.css
	cd client/ && npx tailwindcss build src/tailwind.css -o src/styles.css

.PHONY: client
client: tailwind ## Sets up the react client
	cd client/ && npm run client

.PHONY: serve
serve: ## Sets up the server locally with default options
	go run main.go

.PHONY: dockerbuild
dockerbuild: ## Builds the docker image
	docker build -t cribbage .

.PHONY: dockerrunlocal
dockerrunlocal: ## Runs the latest tag of the built docker image locally on port :8081
	docker run -t -i --env deploy=docker -p 8081:8081 cribbage

.PHONY: wasm
wasm: ## Builds the wasm output for the gowasm client
	GOOS=js GOARCH=wasm go build -o assets/wasm/wa_output.wasm github.com/joshprzybyszewski/cribbage/wasm
