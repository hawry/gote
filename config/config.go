package config

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var globalConfig string
var localConfig string

func init() {
	var homePath string
	if runtime.GOOS == "windows" {
		homePath = os.Getenv("USERPROFILE")
	} else {
		homePath = os.Getenv("HOME")
	}
	globalConfig = fmt.Sprintf("%s/.gote/.config", strings.TrimSuffix(homePath, "/"))
	wd, err := os.Getwd()
	if err != nil {
		wd = "./"
	}
	localConfig = fmt.Sprintf("%s/%s", strings.TrimSuffix(wd, "/"), ".gote")
}

//Configuration is an interface to gather configuration structs with similar methods together
type Configuration interface {
	AccessToken() string
}

//AccessToken returns either the global or the local access token, depending on which is available
func AccessToken() string {
	if configExists(globalConfig) {
		cfg, _, err := LoadGlobal()
		if err != nil {
			return ""
		}
		return cfg.AccessToken()
	}

	if configExists(localConfig) {
		cfg, _, err := LoadLocal()
		if err != nil {
			return ""
		}
		return cfg.AccessToken()
	}

	return ""
}

const (
	//DefaultConfigName is the default name of the Gote configuration file that should exist in the repository root
	DefaultConfigName = ".gote"
	tokenPlaceholder  = "<please insert your personal access token here>"
	rawConfig         = `access_token: %s
remote: %s
repository_owner: %s
repository_name: %s`
	//DefaultGlobalConfigName is the default name of the global configuration file used for application wide settings
	DefaultGlobalConfigName = ".gote_global"
)

//askForAccessToken will ask the user to provide their personal access token for the init-process. If they only give an empty line, it will return a default placeholder instead
func askForAccessToken() string {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Please provide the personal access token for this repository (just press enter if you wish to do this manually later): ")
	txt, err := r.ReadString('\n')
	if err != nil {
		return tokenPlaceholder
	}
	txt = strings.TrimSpace(txt)
	if len(txt) > 0 {
		return txt
	}
	return tokenPlaceholder
}

//AccessToken should be used instead of directly access through the AccessToken attribute. If the user have specified that the token should be taken from an environment variable, this will ensure that the token is updated if the environment variable is changed (and the raw token will not be saved into a new configuration file by accident)
func (c *Local) AccessToken() string {
	if strings.HasPrefix(c.AccessTokenString, "$") {
		//Use environment variable for access token, fetch env
		return os.Getenv(c.AccessTokenString[1:])
	}
	return c.AccessTokenString
}

//configExists will make a quick check if the configuration file already exists and return true if that is the case
func configExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
