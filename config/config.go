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

//Configuration is an interface to gather configuration structs with similar methods together
type Configuration interface {
	AccessToken() string
}

//Local describes a configuration for a specific repository
type Local struct {
	AccessTokenString string `yaml:"access_token"`
	Remote            string `yaml:"remote,omitempty"`
	RepoOwner         string `yaml:"repository_owner,omitempty"`
	Repository        string `yaml:"repository_name,omitempty"`
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

//Default creates and returns a default configuration (used by initialization command)
func Default() (Local, error) {
	if configExists(DefaultConfigName) {
		return Local{}, fmt.Errorf("configuration already exists")
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("could not get current working directory (%v)", err)
		return Local{}, err
	}
	b, r := isGitDir("./")
	if !b {
		return Local{}, &notGitDirError{arg: wd}
	}
	usr, rep := parseRemoteInformation(r)
	defaultConfig := fmt.Sprintf(rawConfig, askForAccessToken(), r, usr, rep)

	f, err := os.Create(DefaultConfigName)
	if err != nil {
		return Local{}, err
	}
	if _, err = f.WriteString(defaultConfig); err != nil {
		return Local{}, err
	}
	return Unmarshal([]byte(defaultConfig))
}

func (c *Local) clean() {
	trim := func(r rune) bool {
		return r == '\n'
	}
	c.AccessTokenString = strings.TrimFunc(c.AccessTokenString, trim)
	c.Remote = strings.TrimFunc(c.Remote, trim)
	c.RepoOwner = strings.TrimFunc(c.RepoOwner, trim)
	c.Repository = strings.TrimFunc(c.Repository, trim)
}

//Create saves a configuration in yaml-format, and makes sure that all fields are valid
func Create(c *Local) (Local, error) {
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

	if c.RepoOwner == "" {
		owner, _ := parseRemoteInformation(c.Remote)
		c.RepoOwner = owner
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
func (c *Local) AccessToken() string {
	if strings.HasPrefix(c.AccessTokenString, "$") {
		//Use environment variable for access token, fetch env
		return os.Getenv(c.AccessTokenString[1:])
	}
	return c.AccessTokenString
}

//Load tries to load a given configuration file
func Load(path string) (Local, error) {
	f, err := os.Open(path)
	if err != nil {
		return Local{}, err
	}
	configData, err := ioutil.ReadAll(f)
	if err != nil {
		return Local{}, err
	}
	return Unmarshal(configData)
}

//Unmarshal will return a Configuration-struct from any given string input, or return an error if it couldn't do it
func Unmarshal(data []byte) (Local, error) {
	c := Local{}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, err
	}
	return c, nil
}

//LoadDefault tries to load the default configuration file as specified by the DefaultConfigName constant
func LoadDefault() (Local, error) {
	return Load(DefaultConfigName)
}

//configExists will make a quick check if the configuration file already exists and return true if that is the case
func configExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
