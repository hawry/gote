package config

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

//Configuration describes a configuration for a repository
type Configuration struct {
	AccessTokenString string `yaml:"access_token"`
	Remote            string `yaml:"remote,omitempty"`
	User              string `yaml:"user,omitempty"`
	Repository        string `yaml:"repository,omitempty"`
}

//GlobalConfiguration describes the global configuration used for gote application wide and can also be used to set some standard values, such as personal access tokens through init
type GlobalConfiguration struct {
	Editor      string `yaml:"editor,omitempty"`
	UseInline   bool   `yaml:"use_inline,omitempty"`
	GlobalToken string `yaml:"global_token,omitempty"`
}

const (
	//DefaultConfigName is the default name of the Gote configuration file that should exist in the repository root
	DefaultConfigName = ".gote"
	tokenPlaceholder  = "<please insert your personal access token here>"
	rawConfig         = `access_token: %s
remote: %s
user: %s
repository: %s`
	//DefaultGlobalConfigName is the default name of the global configuration file used for application wide settings
	DefaultGlobalConfigName = ".gote_global"
)

//Default creates and returns a default configuration (used by initialization command)
func Default() (Configuration, error) {
	if configExists() {
		return Configuration{}, fmt.Errorf("configuration already exists")
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("could not get current working directory (%v)", err)
		return Configuration{}, err
	}
	b, r := isGitDir("./")
	if !b {
		return Configuration{}, &notGitDirError{arg: wd}
	}
	usr, rep := parseRemoteInformation(r)
	defaultConfig := fmt.Sprintf(rawConfig, askForAccessToken(), r, usr, rep)

	f, err := os.Create(DefaultConfigName)
	if err != nil {
		return Configuration{}, err
	}
	if _, err = f.WriteString(defaultConfig); err != nil {
		return Configuration{}, err
	}
	return Unmarshal([]byte(defaultConfig))
}

func (c *Configuration) clean() {
	trim := func(r rune) bool {
		return r == '\n'
	}
	c.AccessTokenString = strings.TrimFunc(c.AccessTokenString, trim)
	c.Remote = strings.TrimFunc(c.Remote, trim)
	c.User = strings.TrimFunc(c.User, trim)
	c.Repository = strings.TrimFunc(c.Repository, trim)
}

//Create saves a configuration in yaml-format, and makes sure that all fields are valid
func Create(c *Configuration) (Configuration, error) {
	//Validate and make sure all needed variables are set, which is remote, user and repo
	c.clean()
	if c.Remote == "" {
		wd, err := os.Getwd()
		if err != nil {
			return *c, err
		}
		b, r := isGitDir("./")
		if !b {
			return *c, &notGitDirError{arg: wd}
		}
		c.Remote = r
	}

	if c.User == "" {
		usr, _ := parseRemoteInformation(c.Remote)
		c.User = usr
	}

	if c.Repository == "" {
		_, rep := parseRemoteInformation(c.Remote)
		c.Repository = rep
	}

	configData, err := yaml.Marshal(&c)
	if err != nil {
		return *c, err
	}
	f, err := os.Create(DefaultConfigName)
	if err != nil {
		return *c, err
	}
	_, err = f.WriteString(string(configData))
	if err != nil {
		return *c, err
	}
	return *c, nil
}

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
func (c *Configuration) AccessToken() string {
	if strings.HasPrefix(c.AccessTokenString, "$") {
		//Use environment variable for access token, fetch env
		return os.Getenv(c.AccessTokenString[1:])
	}
	return c.AccessTokenString
}

//Load tries to load a given configuration file
func Load(path string) (Configuration, error) {
	f, err := os.Open(path)
	if err != nil {
		return Configuration{}, err
	}
	configData, err := ioutil.ReadAll(f)
	if err != nil {
		return Configuration{}, err
	}
	return Unmarshal(configData)
}

//Unmarshal will return a Configuration-struct from any given string input, or return an error if it couldn't do it
func Unmarshal(data []byte) (Configuration, error) {
	c := Configuration{}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, err
	}
	return c, nil
}

//LoadDefault tries to load the default configuration file as specified by the DefaultConfigName constant
func LoadDefault() (Configuration, error) {
	return Load(DefaultConfigName)
}

//configExists will make a quick check if the configuration file already exists and return true if that is the case
func configExists() bool {
	if _, err := os.Stat(DefaultConfigName); os.IsNotExist(err) {
		return false
	}
	return true
}
