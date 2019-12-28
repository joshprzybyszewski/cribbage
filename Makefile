.PHONY: wasm
wasm:
	GOOS=js GOARCH=wasm go build -o assets/wasm/wa_output.wasm wasm/wasm_main.go

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor