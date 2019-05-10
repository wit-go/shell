package shell

// import "log"
import "strings"
import "time"
import "os"
import "os/exec"
import "bufio"
import "github.com/davecgh/go-spew/spew"
import "github.com/svent/go-nbreader"

import log "github.com/sirupsen/logrus"
// import "github.com/wercker/journalhook"

var shellStdout *os.File
var shellStderr *os.File

var spewOn bool = false

func Script(cmds string) int {
	// split on new lines (while we are at it, handle stupid windows text files
	lines := strings.Split(strings.Replace(cmds, "\r\n", "\n", -1), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line) // this is like 'chomp' in perl
		log.Println("LINE:", line)
		time.Sleep(1)
		Run(line)
	}
	return 0
}

/*
func UseJournalctl() {
	journalhook.Enable()
}
*/

func SpewOn() {
	spewOn = true
}

func SetStdout(newout *os.File) {
	shellStdout = newout
}

func SetStderr(newerr *os.File) {
	shellStderr = newerr
}

func Run(cmdline string) int {
	log.Println("START " + cmdline)

	cmd := strings.TrimSpace(cmdline) // this is like 'chomp' in perl
	cmdArgs := strings.Fields(cmd)
	if (len(cmdArgs) == 0) {
		log.Println("END   ", cmd)
		return 0 // nothing to do
	}
	if (cmdArgs[0] == "cd") {
		if (len(cmdArgs) > 1) {
			log.Println("os.Chdir()", cmd)
			os.Chdir(cmdArgs[1])
		}
		log.Println("END   ", cmd)
		return 0 // nothing to do
	}

	process := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	pstdout, _ := process.StdoutPipe()
	pstderr, _ := process.StderrPipe()

	if (spewOn) {
		spew.Dump(pstdout)
	}

	process.Start()

	if (shellStdout == nil) {
		shellStdout = os.Stdout
	}

	f := bufio.NewWriter(shellStdout)

	newreader := bufio.NewReader(pstdout)
	nbr := nbreader.NewNBReader(newreader, 1024)

	newerrreader := bufio.NewReader(pstderr)
	nbrerr := nbreader.NewNBReader(newerrreader, 1024)

	totalCount := 0

	for {
		time.Sleep(2 * time.Millisecond)   // only check the buffer 500 times a second
		// log.Println("sleep done")

		// tight loop that reads 1K at a time until buffer is empty
		for {
			oneByte := make([]byte, 1024)
			count, err := nbr.Read(oneByte)
			totalCount += count

			if (err != nil) {
				// log.Println("Read() count = ", count, "err = ", err)
				oneByte = make([]byte, 1024)
				count, err = nbr.Read(oneByte)
				f.Write([]byte(string(oneByte)))
				f.Flush()
			}
			f.Write([]byte(string(oneByte)))
			f.Flush()
			if (count == 0) {
				break
			}
		}

		if (totalCount != 0) {
			log.Println("totalCount = ", totalCount)
		}

		//
		// HANDLE STDERR
		// HANDLE STDERR
		//
		oneByte := make([]byte, 1024)
		count, err := nbrerr.Read(oneByte)

		if (err != nil) {
			oneByte = make([]byte, 1024)
			count, err = nbrerr.Read(oneByte)
			f.Write([]byte(string(oneByte)))
			f.Flush()

			log.Println("Read() count = ", count, "err = ", err)
			// spew.Dump(process.Process)
			// spew.Dump(process.ProcessState)
			err := process.Wait()
			if err != nil {
				if (spewOn) {
					spew.Dump(err.(*exec.ExitError))
					spew.Dump(process.ProcessState)
				}
				stuff := err.(*exec.ExitError)
				log.Println("ERROR ", stuff)
				log.Println("END   ", cmdline)
				return -1
			}
			log.Println("END   ", cmdline)
			return 0
		} else {
			f.Write([]byte(string(oneByte)))
			f.Flush()
		}

		// spew.Dump(reflect.ValueOf(cmd.Process).Elem())
	}

	err := process.Wait()

	if err != nil {
		if (spewOn) {
			spew.Dump(err.(*exec.ExitError))
			spew.Dump(process.ProcessState)
		}
		stuff := err.(*exec.ExitError)
		log.Println("ERROR ", stuff)
		log.Println("END   ", cmdline)
		return -1
	}
	log.Println("END   ", cmdline)
	return 0
}

func Daemon(cmdline string, timeout time.Duration) int {
	for {
		Run(cmdline)
		time.Sleep(timeout)
	}
}
