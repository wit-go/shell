package main

/*
import "log"
import "reflect"
*/

import "os"

// import "github.com/davecgh/go-spew/spew"

import "git.wit.org/wit/shell"

func main() {
	shell.Run("ls /tmp")

	shell.Run("ping -c 3 localhost")

	// slow down the polling to every 2 seconds
	shell.SetDelayInMsec(2000)

	shell.Run("ping -c 4 localhost")

	// capture ping output into a file
	fout, _ := os.Create("/tmp/example1.ping.stdout")
	ferr, _ := os.Create("/tmp/example1.ping.stderr")
	shell.SetStdout(fout)
	shell.SetStderr(ferr)

	shell.Run("ping -c 5 localhost")

	// turn out process exit debugging
	shell.SpewOn()

	fout, _ = os.Create("/tmp/example1.fail.stdout")
	ferr, _ = os.Create("/tmp/example1.fail.stderr")
	shell.SetStdout(fout)
	shell.SetStderr(ferr)

	// TODO: this might not be working
	// check error handling
	shell.Run("ls /tmpthisisnothere")
}
