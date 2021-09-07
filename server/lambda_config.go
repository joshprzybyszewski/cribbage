//+build lambda

package server

import (
	"fmt"
	"io/ioutil"

	"github.com/rakyll/globalconf"
)

// for lambda, we want to only use env vars. Load globalconf with an empty temp file to just grab the env values
func loadVarsFromINI() {
	tmpFile, err := ioutil.TempFile(``, `config.ini`)
	if err != nil {
		panic(fmt.Sprintf("getConfigFile TempFile err: %+v", err))
	}

	options := &globalconf.Options{
		EnvPrefix: `CRIBBAGE_`,
		Filename:  tmpFile.Name(),
	}

	conf, err := globalconf.NewWithOptions(options)
	if err != nil {
		panic(fmt.Sprintf("globalconf.NewWithOptions error (from %s): %+v", tmpFile.Name(), err))
	}

	conf.ParseAll()
}
