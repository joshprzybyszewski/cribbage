// +build js,wasm

package callbacks

import (
	"honnef.co/go/js/dom/v2"
)

func goToPath(newPath string) {
	loc := dom.GetWindow().Location()
	newUrl := loc.Host() + newPath
	prefix := `http://`
	if len(newUrl) < len(prefix) || newUrl[:len(prefix)] != prefix {
		newUrl = prefix + newUrl
	}
	loc.Call("replace", newUrl)
}
