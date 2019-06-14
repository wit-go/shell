package shell

import "runtime"

func Execname(filename string) string {
	if runtime.GOOS == "windows" {
		return filename + ".exe"
	}
	return filename
}
