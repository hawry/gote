package format

import "runtime"

//TrimNewlines is supposed to be used with strings.TrimLeftFunc, or strings.TrimRightFunc instead of using anonymous functions (and this one takes care of different OS:es as well)
func TrimNewlines(r rune) bool {
	if runtime.GOOS == "windows" {
		return r == '\n' || r == '\r'
	}
	return r == '\n'
}
