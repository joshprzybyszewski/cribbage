//+build lambda

package server

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/rakyll/globalconf"
)

// for lambda, we want to only use env vars. Pass in an empty temp file.
func loadVarsFromINI() {
	tmpFile, err := ioutil.TempFile(``, `config.ini`)
	if err != nil {
		panic(fmt.Sprintf("getConfigFile TempFile err: %+v", err))
	}

	parseFlagsFromConfigFile(tmpFile.Name())
}

func parseFlagsFromConfigFile(confFileName string) {
	log.Printf("parseFlagsFromConfigFile from %q\n", confFileName)

	options := &globalconf.Options{
		EnvPrefix: `CRIBBAGE_`,
		Filename:  confFileName,
	}

	conf, err := globalconf.NewWithOptions(options)
	if err != nil {
		panic(fmt.Sprintf("globalconf.NewWithOptions error (from %s): %+v", confFileName, err))
	}

	conf.ParseAll()
}
