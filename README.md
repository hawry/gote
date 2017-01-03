# gote
This project replaces the existing [git-note][1] application for creating issues and notes. The basic functionality is still the same, but this application only works for local repositories instead of using a global one (which means that the 'note'-functionality is gone and it's much more of a tool to create issues in a quick manner).

The tool relies on the use of personal access tokens with read and write access to public and private repositories where your user account have access.

## Usage
```
$ gote note
> The brown fox jumps over something something dark side. I think we have cookies!
```

When the input receives a newline character, the body is shortened and formatted and then pushed to the repository where the command was called from, creating an issue with a title and a body.

## TODO
* [ ] Markdown support
* [ ] Using $EDITOR instead of own secondary-prompt for title and body creation of issues
* [ ] Close/Re-open issues directly from the prompt
* [ ] Custom labels
* [ ] Using environment variables instead of files for access tokens (customizable)

[1]: https://github.com/hawry/git-note
