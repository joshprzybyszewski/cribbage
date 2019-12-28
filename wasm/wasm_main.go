// +build js,wasm

package main

import (
	"fmt"
	"net/url"
	"syscall/js"

	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
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

	// TODO populate game
	g := model.Game{}

	cardClickCB := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("this, args := %+v, %+v\n", this, args)
		return nil
	})

	fmt.Printf("dom.GetWindow(): %+v\n", dom.GetWindow())
	fmt.Printf("dom.GetWindow().Document(): %+v\n", dom.GetWindow().Document())
	elems := dom.GetWindow().Document().GetElementsByClassName(".card")
	fmt.Printf("elems: %+v\n", elems)

	ret := js.Global().Get("document").
		Call("getElementsByClassName", "card")
	// Call("addEventListener", "change", cardClickCB)
	fmt.Printf("ret: %+v\n", ret)

	fmt.Printf("g: %+v\n", g)

	<-done

	cardClickCB.Release()

}
