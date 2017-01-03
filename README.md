# gote
This project replaces the existing [git-note][1] application for creating issues and notes. The basic functionality is still the same, but this application only works for local repositories instead of using a global one (which means that the 'note'-functionality is gone and it's much more of a tool to create issues in a quick manner).

The tool relies on the use of personal access tokens with read and write access to public and private repositories where your user account have access.

## Usage
```
$ gote note
> The brown fox jumps over something something dark side. I think we have cookies!
```

When the input receives a newline character, the body is shortened and formatted and then pushed to the repository where the command was called from, creating an issue with a title and a body.

## Installation
Either clone this repository and build it from source or download any of the precompiled binaries.

## Configuration
By using the `init` command, gote will create a configuration file in the current working directory, assuming it's a git repository. Gote will parse your .git/config and add remote address, username and repository name from it and add to a configuration file (default name `.gote`). Currently, gote assumes that your remote is named `origin`, otherwise gote will not be able to find the information. If you have any other remote name than origin, you can still create the configuration file manually.

### Example configuration
```
access_token: <40 char access token>
remote: git@github.com:hawry/gote
user: hawry
repository: gote
```

You will have to supply your [personal access token][2] manually if you didn't provide it during the init process.



## TODO
* [ ] Markdown support
* [ ] Using $EDITOR instead of own secondary-prompt for title and body creation of issues
* [ ] Close/Re-open issues directly from the prompt
* [ ] Custom labels
* [ ] Using environment variables instead of files for access tokens (customizable)

[1]: https://github.com/hawry/git-note
[2]: https://help.github.com/articles/creating-an-access-token-for-command-line-use/
