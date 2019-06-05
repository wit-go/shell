package main

import "log"
import "fmt"
import "git.wit.com/wit/shell"

func main() {
	tmp, output, err := shell.Run("cat /etc/issue")
	log.Println("cat /etc/issue returned", tmp, "error =", err)
	fmt.Print(output)
}
