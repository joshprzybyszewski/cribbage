// +build js,wasm

package main

import (
	"fmt"
	"net/url"
	"syscall/js"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/wasm/callbacks"
)

func main() {
	fmt.Println("hello wasm")

	path, err := getCurrentURLPath()
	if err != nil {
		return
	}

	if path == `/` {
		startHomePage()
	} else {
		startGamePage()
	}

	println(`exiting main`)
}

func getCurrentURLPath() (string, error) {
	href := js.Global().Get("location").Get("href")
	u, err := url.Parse(href.String())
	if err != nil {
		return ``, err
	}
	return u.Path, nil
}

func startHomePage() {
	println(`startHomePage`)
	var done chan struct{} = make(chan struct{}, 1)

	usernameChangeCB := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("this, args := %+v, %+v\n", this, args)
		return nil
	})

	js.Global().Get("document").
		Call("getElementById", "createUN").
		Call("addEventListener", "change", usernameChangeCB)

	displaynameChangeCB := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("this, args := %+v, %+v\n", this, args)
		return nil
	})

	js.Global().Get("document").
		Call("getElementById", "createDN").
		Call("addEventListener", "change", displaynameChangeCB)

	createUserCB := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("this, args := %+v, %+v\n", this, args)
		return nil
	})

	js.Global().Get("document").
		Call("getElementById", "createForm").
		Call("addEventListener", "submit", createUserCB)

	<-done

	usernameChangeCB.Release()
	displaynameChangeCB.Release()
	createUserCB.Release()
}

func startGamePage() {
	println(`startGamePage`)
	var done chan struct{} = make(chan struct{}, 1)

	// TODO populate game and myID
	g := model.Game{}
	var myID model.PlayerID

	callbacks.SetupGamePage(g, myID)

	<-done
}
