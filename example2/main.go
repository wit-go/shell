package main

import "log"
// import "fmt"
import "git.wit.org/wit/shell"

func main() {
	err := shell.Run("cat /etc/issue")
	log.Println("cat /etc/issue returned", err)
	// fmt.Print(output)
}
