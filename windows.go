// +build windows

// put stuff in here that you only want compiled under windows

package shell

import "log"

// import "git.wit.org/wit/shell"
// import "github.com/davecgh/go-spew/spew"

func handleSignal(err interface{}, ret int) {
	log.Println("handleSignal() windows doesn't do signals")
}

func UseJournalctl() {
	log.Println("journalctl doesn't exist on windows")
}
