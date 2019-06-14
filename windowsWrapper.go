package shell

import "runtime"
import "strings"

func Execname(filename string) string {
	if runtime.GOOS != "windows" {
		return filename
	}
	return Path(filename) + ".exe"
}

func Path(filename string) string {
	if runtime.GOOS != "windows" {
		return filename
	}
	filename = strings.Replace(filename, "/", "\\", -1)
	return filename
}
