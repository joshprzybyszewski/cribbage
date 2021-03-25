package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rakyll/globalconf"
	ini "gopkg.in/ini.v1"
)

// loadVarsFromINI will check the INI for the given environment and any environment variables
// It will give the following priority:
// 1. Command-line value
// 2. Environment Variable
// 3. INI value
// 4. Default value of flag
// NOTE: All envvars need the prefix `CRIBBAGE_` and should exist in ALLCAPS.
// NOTE: The INI file is pulled from the inis/ directory and is determined by the
//       environment var `deploy` (set to `docker`, `prod`, etc.)
func loadVarsFromINI() {
	parseFlagsFromConfigFile(getConfigFile())
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
		panic(fmt.Sprintf("getConfigFile ini.LooseLoad err: %+v", err))
	}

	_, err = f.WriteTo(tmpFile)
	if err != nil {
		panic(fmt.Sprintf("getConfigFile f.WriteTo(tmpFile) err: %+v", err))
	}

	return tmpFile.Name()
}

// getEnvironment returns the environment variable set for `deploy`
func getEnvironment() string {
	v := os.Getenv(`deploy`)
	if v == `` {
		return `default`
	}
	return v
}
