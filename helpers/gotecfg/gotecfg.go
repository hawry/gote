package gotecfg

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	localCfg        localConfiguration
	globalCfg       globalConfiguration
	hasGlobalConfig bool
}

type localConfiguration struct {
	AccessTokenString string `yaml:"access_token,omitempty"`
	Provider          string `yaml:"provider,omitempty"`
	Remote            string `yaml:"remote"`
}

type globalConfiguration struct {
	AccessTokenString string `yaml:"access_token,omitempty"`
	Editor            string `yaml:"editor,omitempty"`
}

var localConfigPath, globalConfigPath string

func init() {
	log.Printf("init")
	var homePath string
	if runtime.GOOS == "windows" {
		homePath = os.Getenv("USERPROFILE")
	} else {
		homePath = os.Getenv("HOME")
	}

	wd, err := os.Getwd()
	if err != nil {
		wd = "./"
	}
	setConfigPaths(fmt.Sprintf("%s/.gote/.config", strings.TrimSuffix(homePath, "/")), fmt.Sprintf("%s/%s", strings.TrimSuffix(wd, "/"), ".gote"))
}

//to enable tests, we need to be able to redirect these
func setConfigPaths(globalPath, localPath string) {
	globalConfigPath = globalPath
	localConfigPath = localPath
}

func New() *Configuration {
	f, err := os.Open(globalConfigPath)
	if err != nil {
		log.Panic(err)
	}
	c := Configuration{}
	c.globalCfg = globalConfiguration{}
	return &c
}

func (c *Configuration) AccessToken() string {
	if c.globalCfg.AccessTokenString != "" {
		if strings.HasPrefix(c.globalCfg.AccessTokenString, "$") {
			return os.Getenv(c.globalCfg.AccessTokenString[1:])
		}
		return c.globalCfg.AccessTokenString
	}

	if strings.HasPrefix(c.localCfg.AccessTokenString, "$") {
		return os.Getenv(c.localCfg.AccessTokenString[1:])
	}
	return c.localCfg.AccessTokenString
}

func load(r io.Reader, v interface{}) error {
	rd, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(rd, &v); err != nil {
		log.Panic(err)
	}
	return nil
}
