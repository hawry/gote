package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

//Global describes a global configuration for gote
type Global struct {
	Editor            string `yaml:"editor,omitempty"`
	AccessTokenString string `yaml:"access_token,omitempty"`
}

//LoadGlobal loads the global configuration
func LoadGlobal() (cfg Global, cfgExists bool, err error) {
	exists := configExists(globalConfig)
	g := Global{}
	f, err := os.Open(globalConfig)
	if err != nil {
		return g, exists, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return g, exists, err
	}
	log.Printf("debug: raw global; %v", string(b))
	err = yaml.Unmarshal(b, &g)
	if err != nil {
		return g, exists, err
	}
	log.Printf("debug: global loaded (%v)", g)
	return g, exists, nil
}

//AccessToken implements the Configuration interface
func (g *Global) AccessToken() string {
	if strings.HasPrefix(g.AccessTokenString, "$") {
		return os.Getenv(g.AccessTokenString[1:])
	}
	return g.AccessTokenString
}
