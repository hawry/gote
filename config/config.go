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
	AccessToken string `yaml:"access_token"`
	Remote      string `yaml:"remote,omitempty"`
	User        string `yaml:"user,omitempty"`
	Repository  string `yaml:"repository,omitempty"`
}

const (
	//DefaultConfigName is the default name of the Gote configuration file that should exist in the repository root
	DefaultConfigName = ".gote"
	tokenPlaceholder  = "<please insert your personal access token here>"
	rawConfig         = `access_token: %s
remote: %s
user: %s
repository: %s`
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
