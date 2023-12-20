// +build linux,go1.7

// put stuff in here that you only want compiled under linux

package shell

import "log"
import "os"
import "os/signal"
import "syscall"

// import "runtime"
// import "time"
// import "reflect"

// import "go.wit.com/shell"
// import "github.com/davecgh/go-spew/spew"

import "github.com/wercker/journalhook"

var sigChan	chan os.Signal

func handleSignal(err interface{}, ret int) {
	log.Println("handleSignal() only should be compiled on linux")
	sigChan = make(chan os.Signal, 3)
	signal.Notify(sigChan, syscall.SIGUSR1)
}

func UseJournalctl() {
	journalhook.Enable()
}
