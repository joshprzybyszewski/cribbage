// +build js,wasm

package callbacks

import (
	"honnef.co/go/js/dom/v2"
)

func getClickHandlerForID(id string, cb func(e dom.Event)) Releaser {
	elem := dom.GetWindow().Document().GetElementByID(id)
	return elem.AddEventListener(`click`, false, cb)
}

func getChangeHandlerForID(id string, cb func(e dom.Event)) Releaser {
	elem := dom.GetWindow().Document().GetElementByID(id)
	return elem.AddEventListener(`change`, false, cb)
}

func getInputHandlerForID(id string, cb func(e dom.Event)) Releaser {
	elem := dom.GetWindow().Document().GetElementByID(id)
	return elem.AddEventListener(`input`, false, cb)
}

func getEnterKeyHandlerForID(id string, cb func(e dom.Event)) Releaser {
	elem := dom.GetWindow().Document().GetElementByID(id)
	return elem.AddEventListener(`keyup`, false, func(e dom.Event) {
		ke, ok := e.(*dom.KeyboardEvent)
		if ok && ke.Key() == `Enter` {
			cb(e)
		}
	})
}
