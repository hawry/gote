package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	gcfg "gopkg.in/gcfg.v1"
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
    remote: %s
  `
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
	defaultConfig = fmt.Sprintf(defaultConfig, r)

	f, err := os.Create(DefaultConfigName)
	if err != nil {
		return Configuration{}, err
	}
	if _, err = f.WriteString(defaultConfig); err != nil {
		return Configuration{}, err
	}
	return Unmarshal([]byte(defaultConfig))
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

func configExists() bool {
	if _, err := os.Stat(DefaultConfigName); os.IsNotExist(err) {
		return false
	}
	return true
}

//isGitDir will return true if the specified directory is a git-repository and if so, will also return the remote name
func isGitDir(path string) (bool, string) {
	gitConfig := struct {
		Remote map[string]*struct {
			URL string
		}
	}{}
	err := gcfg.ReadFileInto(&gitConfig, fmt.Sprintf("%s/.git/config", strings.TrimSuffix(path, "/")))
	err = gcfg.FatalOnly(err)
	if err != nil {
		//Assume directory isn't a git repository since a fatal error was thrown
		return false, ""
	}
	//If remote origin can be found, return the remote address of that remote
	if val, v := gitConfig.Remote["origin"]; v {
		return true, val.URL
	}
	return false, ""
}
