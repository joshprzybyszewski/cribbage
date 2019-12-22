.PHONY: wasm
wasm:
	GOOS=js GOARCH=wasm go build -o wa_output.wasm wasm/wasm_main.go