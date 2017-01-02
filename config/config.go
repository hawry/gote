package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

//Configuration describes a configuration for a repository
type Configuration struct {
	AccessToken string
	Remote      string
	Username    string
	Repository  string
}

const (
	//DefaultConfigName is the default name of the Gote configuration file that should exist in the repository root
	DefaultConfigName = ".gote"
)

//Default creates and returns a default configuration (used by initialization command)
func Default() (Configuration, error) {
	var defaultConfig = `github:
    access_token: <insert your github access token here>
  `
	if configExists() {
		return Configuration{}, fmt.Errorf("configuration already exists")
	}
	f, err := os.Create(DefaultConfigName)
	if err != nil {
		return Configuration{}, err
	}
	if _, err = f.WriteString(defaultConfig); err != nil {
		return Configuration{}, err
	}
	return Marshal(defaultConfig)
}

//Load tries to load a given configuration file
func Load(path string) (Configuration, error) {
	return Configuration{}, fmt.Errorf("not implemented yet")
}

//Marshal will return a Configuration-struct from any given string input, or return an error if it couldn't do it
func Marshal(data string) (Configuration, error) {
	c := Configuration{}
	if err := yaml.Unmarshal([]byte(data), &c); err != nil {
		return c, err
	}
	return c, nil
}

//LoadDefault tries to load the default configuration file as specified by the DefaultConfigName constant
func LoadDefault() (Configuration, error) {
	return Load(DefaultConfigName)
}

func configExists() bool {
	if _, err := os.Stat(DefaultConfigName); os.IsNotExist(err) {
		return false
	}
	return true
}
