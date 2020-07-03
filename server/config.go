package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rakyll/globalconf"
	ini "gopkg.in/ini.v1"
)

func loadVarsFromINI() error {
	parseFlagsFromConfigFile(getConfigFile())
	return nil
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

func getConfigFile() string {
	tmpFile, err := ioutil.TempFile(``, `config.ini`)
	if err != nil {
		panic(fmt.Sprintf("getConfigFile TempFile err: %+v", err))
	}

	iniPath := `inis/` + getEnvironment() + `/cribbage.ini`
	log.Printf("ini.LooseLoad from %q\n", iniPath)

	f, err := ini.LooseLoad(iniPath)
	if err != nil {
		panic(fmt.Sprintf("getConfigFile LooseLoad err: %+v", err))
	}

	_, err = f.WriteTo(tmpFile)
	if err != nil {
		panic(fmt.Sprintf("getConfigFile WriteTo err: %+v", err))
	}

	return tmpFile.Name()
}

func getEnvironment() string {
	v := os.Getenv(`deployed`)
	if v == `` {
		return `local`
	}
	return v
}
