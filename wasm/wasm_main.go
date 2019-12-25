// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("hello wasm")

	listenToHome()
	listenForGamePage()

	println(`exiting main`)
}

func listenToHome() {
	println(`listenToHome`)

	cb := js.NewCallback(func(args []js.Value) {
		println(`wasm callback`)
	})
	// TODO cb.Release()

	js.Global().Get("document").
		Call("getElementById", "create-user-submit").
		Call("addEventListener", "submit", cb)
}

func listenForGamePage() {
	println(`listenForGamePage`)
}
