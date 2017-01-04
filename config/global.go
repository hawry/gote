package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var configPath string

func init() {
	var homePath string
	if runtime.GOOS == "windows" {
		homePath = os.Getenv("USERPROFILE")
	} else {
		homePath = os.Getenv("HOME")
	}
	configPath = fmt.Sprintf("%s/.gote/.config", strings.TrimSuffix(homePath, "/"))
}

//Global describes a global configuration for gote
type Global struct {
	Editor            string `yaml:"editor,omitempty"`
	AccessTokenString string `yaml:"global_token,omitempty"`
}

//LoadGlobal loads the global configuration
func LoadGlobal() (globalConfig Global, cfgExists bool, err error) {
	exists := configExists(configPath)
	g := Global{}
	f, err := os.Open(configPath)
	if err != nil {
		return g, exists, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return g, exists, err
	}
	err = yaml.Unmarshal(b, &g)
	if err != nil {
		return g, exists, err
	}
	return g, exists, nil
}

//AccessToken implements the Configuration interface
func (g *Global) AccessToken() string {
	if strings.HasPrefix(g.AccessTokenString, "$") {
		return os.Getenv(g.AccessTokenString[1:])
	}
	return g.AccessTokenString
}
