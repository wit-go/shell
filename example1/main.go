package main

/*
import "log"
import "reflect"
import "os"
*/

// import "github.com/davecgh/go-spew/spew"

import "git.wit.com/jcarr/shell"

func main() {
	shell.SpewOn()

	shell.Run("ls /tmp")

	shell.Run("ping -c 4 localhost")

	// slow down the polling to every 2 seconds
	shell.SetDelayInMsec(2000)
	shell.Run("ping -c 4 localhost")

	// TODO: this might not be working
	// check error handling
	shell.Run("ls /tmpthisisnothere")
}
