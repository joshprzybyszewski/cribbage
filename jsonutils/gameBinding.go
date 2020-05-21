package jsonutils

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin/binding"

	"github.com/joshprzybyszewski/cribbage/model"
)

var _ binding.Binding = GameBinding
var _ binding.BindingBody = GameBinding
var (
	GameBinding = gameBinding{}
)

type gameBinding struct{}

func (gb gameBinding) Name() string {
	return `gameBinding`
}

func (gb gameBinding) Bind(r *http.Request, obj interface{}) error {
	if r == nil || r.Body == nil {
		return errors.New(`invalid request`)
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return bindGame(bs, obj)
}

func (gb gameBinding) BindBody(body []byte, obj interface{}) error {
	return bindGame(body, obj)
}

func bindGame(bs []byte, obj interface{}) error {
	gIn, ok := obj.(*model.Game)
	if !ok {
		return errors.New(`gameBinding only works with model.Game`)
	}
	g, err := UnmarshalGame(bs)
	if err != nil {
		return err
	}
	*gIn = g
	return nil
}
