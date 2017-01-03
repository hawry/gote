package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hawry/gote/helpers"

	gcfg "gopkg.in/gcfg.v1"
)

const (
	// httpsRegex = "\\A(?P<protocol>https://)(?P<host>[^/]*)/(?P<user>[^/]*)/(?P<repo>[^/]*)\\z"
	identifyRegex = "\\A(?P<protocol>git@|git://|https://)(?P<host>[^/:]*)[/|:]{1}(?P<user>[^/]*)/(?P<repo>[^/]*)\\z"
)

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

//parseRemoteInformation will return the user and the repository name from a given remote
func parseRemoteInformation(remote string) (user, repository string) {
	r := regexp.MustCompile(identifyRegex)
	result := helpers.ToMap(r.FindStringSubmatch(remote), r.SubexpNames())
	return result["user"], result["repo"]
}
