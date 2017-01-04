package config

import "fmt"

type notGitDirError struct {
	arg string
}

type configNotFound struct {
	arg string
}

func (e *notGitDirError) Error() string {
	return fmt.Sprintf("%s is not a valid git repository", e.arg)
}

func (e *configNotFound) Error() string {
	return fmt.Sprintf("%s could not be found or contains invalid syntax", e.arg)
}

//IsNotGitDir will make a check if the given error struct is of the type returned whenever a directory isn't a git-repository
func IsNotGitDir(err error) bool {
	if _, ok := err.(*notGitDirError); ok {
		return true
	}
	return false
}

//IsNotFound will check if the error is of type 'config not found'
func IsNotFound(err error) bool {
	if _, ok := err.(*configNotFound); ok {
		return true
	}
	return false
}
