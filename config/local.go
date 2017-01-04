package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

//Local describes a configuration for a specific repository
type Local struct {
	AccessTokenString string `yaml:"access_token"`
	Remote            string `yaml:"remote,omitempty"`
	RepoOwner         string `yaml:"repository_owner,omitempty"`
	Repository        string `yaml:"repository_name,omitempty"`
}

//LoadLocal returns the local configuration
func LoadLocal() (cfg Local, cfgExists bool, err error) {
	ex := configExists(localConfig)
	l := Local{}
	f, err := os.Open(localConfig)
	if err != nil {
		return l, ex, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return l, ex, err
	}
	err = yaml.Unmarshal(b, &l)
	if err != nil {
		return l, ex, err
	}
	return l, ex, nil
}

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
