// +build js,wasm

package callbacks

import (
	"honnef.co/go/js/dom/v2"
)

func goToPath(newPath string) {
	loc := dom.GetWindow().Location()
	// all of the wasm paths exist under the wasm "root directory"
	newUrl := loc.Host() + `/wasm` + newPath
	prefix := loc.Protocol() + `//`
	if len(newUrl) < len(prefix) || newUrl[:len(prefix)] != prefix {
		newUrl = prefix + newUrl
	}
	loc.Call("replace", newUrl)
}
