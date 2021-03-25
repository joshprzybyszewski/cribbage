// +build js,wasm

package callbacks

type Releaser interface {
	Release()
}
