# gote
This project replaces the existing [git-note][1] application for creating issues and notes. The basic functionality is still the same, but this application only works for local repositories instead of using a global one (which means that the 'note'-functionality is gone and it's much more of a tool to create issues in a quick manner).

The tool relies on the use of personal access tokens with read and write access to public and private repositories where your user account have access.

## Download
Visit [the releases][3] page to download the source, or pre-compiled binaries for Linux and Windows.

## Usage

Make sure you are in an initialized git-repository root folder. Then run the `init`-command in the project root to bootstrap gote and answer the access question. This command will fetch as much information as it can from your `.git/config`. If you wish to enter the information manually, use the `--interactive` flag.

If you need any details regarding any command, just run `gote help` or `gote help <command>`. You can also find the generated help files in [the /docs directory](gote)

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

You will have to supply your [personal access token][2] manually if you didn't provide it during the init process.

### Security
The access token can be placed in an environment variable instead of directly in the .gote-configuration file to reduce the risk of leaking sensitive data to a remote endpoint. To use an environment variable, replace your access token in the configuration file with your environment variable name (i.e. `$GOTE_ACCESS`) and export the access token string `export GOTE_ACCESS=<access token>`.

#### Example configuration with environment variable
```
access_token: $GOTE_ACCESS
remote: git@github.com:hawry/gote
user: hawry
repository: gote
```

You will have to supply your [personal access token][2] manually if you didn't provide it during the init process.

### Global configuration
To avoid using access tokens in the repository folder, a global configuration can be used with a machine wide access token setting as well as your choice of editor. The global configuration resides in `$HOME/.gote/.config` or if you are on windows `%USERPROFILE%/.gote/.config` (normally `C:\Users\<username>\`). In the global configuration, it's still possible to use environment variables.

#### Examples
**Using hard coded access token**
```
access_token: 123...123
editor: vim
```

**Using environment variables**
```
access_token: $GOTE_ACCESS_TOKEN
editor: $EDITOR
```

If the editor variable is empty, the default is for gote to check the environment variable `$EDITOR` first, and if not found or invalid gote will use the default (somewhat primitime) secondary prompt.

[1]: https://github.com/hawry/git-note
[2]: https://help.github.com/articles/creating-an-access-token-for-command-line-use/
[3]: https://github.com/Hawry/gote/releases
[4]: gote.md
