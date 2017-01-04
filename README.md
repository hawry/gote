# gote
This project replaces the existing [git-note][1] application for creating issues and notes. The basic functionality is still the same, but this application only works for local repositories instead of using a global one (which means that the 'note'-functionality is gone and it's much more of a tool to create issues in a quick manner).

The tool relies on the use of personal access tokens with read and write access to public and private repositories where your user account have access.

## Usage

Make sure you are in an initialized git-repository root folder. Then run the `init`-command in the project root to bootstrap gote and answer the access question. This command will fetch as much information as it can from your `.git/config`. If you wish to enter the information manually, use the `--interactive` flag.

```
/my/awesome/project$ gote init
Please provide the personal access token for this repository (just press enter if you wish to do this manually later):
```

```
$ gote note
> The brown fox jumps over something something dark side. I think we have cookies!
```

When the input receives a newline character, the body is shortened and formatted and then pushed to the repository where the command was called from, creating an issue with a title and a body.

If the $EDITOR environment variable is set, gote asks for input from the editor instead of through the command line. Gote will always take the first line and turn it into the Issue-title, and the rest of the text will be the body. If no body is provided, the title is repeated as the body. Using the editor enables you to use markdown in your issue body.

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

You will have to supply your [personal access token][2] manually if you didn't provide it during the init process. To allow gote to access issues in private repositories make sure to give the token full `repo`-access.

### Security
The access token can be placed in an environment variable instead of directly in the .gote-configuration file to reduce the risk of leaking sensitive data to a remote endpoint. To use an environment variable, replace your access token in the configuration file with your environment variable name (i.e. `$GOTE_ACCESS`) and export the access token string `export GOTE_ACCESS=<access token>`.

#### Example configuration with environment variable
```
access_token: $GOTE_ACCESS
remote: git@github.com:hawry/gote
user: hawry
repository: gote
```

## Upcoming features
* [x] Using $EDITOR instead of own secondary-prompt for title and body creation of issues
* [ ] Close/Re-open issues directly from the prompt
* [ ] Custom labels
* [ ] Global configuration to set default access tokens, specify editor in other means than through environment variables and also default user (for collabs and such)
* [ ] Using environment variables instead of files for access tokens (customizable)

[1]: https://github.com/hawry/git-note
[2]: https://help.github.com/articles/creating-an-access-token-for-command-line-use/
