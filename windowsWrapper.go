package shell

import "log"
import "runtime"
import "strings"

func Execname(filename string) string {
	if runtime.GOOS != "windows" {
		return filename
	}
	return Path(filename) + ".exe"
}

func Path(filename string) string {
	log.Println("shell.Path() START filename =", filename)
	if runtime.GOOS != "windows" {
		log.Println("shell.Path() END filename =", filename)
		return filename
	}
	filename = strings.Replace(filename, "/", "\\", -1)
	log.Println("shell.Path() END filename =", filename)
	return filename
}
