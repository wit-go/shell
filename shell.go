package shell

import "fmt"
import "strings"
import "time"
import "os"
import "os/exec"
import "bufio"
import "bytes"
import "io"

import "github.com/davecgh/go-spew/spew"
import "github.com/svent/go-nbreader"

import log "github.com/sirupsen/logrus"
// import "github.com/wercker/journalhook"

var shellStdout *os.File
var shellStderr *os.File

var spewOn bool = false
var msecDelay int = 20	// number of milliseconds to delay between reads with no data

var buf bytes.Buffer

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

func SetDelayInMsec(msecs int) {
	msecDelay = msecs
}

func SetStdout(newout *os.File) {
	shellStdout = newout
}

func SetStderr(newerr *os.File) {
	shellStderr = newerr
}

/*
func Capture(cmdline string) (int, string) {
	val, _, _ := Run(cmdline)

	if (val != 0) {
		log.Println("shell.Capture() ERROR")
	}

	return val, buf.String()
}
*/

func Run(cmdline string) (int, string, error) {
	log.Println("START " + cmdline)

	cmd := strings.TrimSpace(cmdline) // this is like 'chomp' in perl
	cmdArgs := strings.Fields(cmd)
	if (len(cmdArgs) == 0) {
		log.Println("END   ", cmd)
		return 0, "", fmt.Errorf("") // nothing to do
	}
	if (cmdArgs[0] == "cd") {
		if (len(cmdArgs) > 1) {
			log.Println("os.Chdir()", cmd)
			os.Chdir(cmdArgs[1])
		}
		log.Println("END   ", cmd)
		return 0, "", fmt.Errorf("") // nothing to do
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

	tmp := bufio.NewReader(pstderr)
	go NonBlockingReader(tmp, shellStderr)

	totalCount := 0

	var dead bool = false
	for (dead == false) {
		time.Sleep(time.Duration(msecDelay) * time.Millisecond)   // only check the buffer 500 times a second
		// log.Println("sleep done")

		var empty bool = false
		// tight loop that reads 1K at a time until buffer is empty
		for (empty == false) {
			oneByte := make([]byte, 1024)
			count, err := nbr.Read(oneByte)
			totalCount += count

			if (err != nil) {
				log.Println("Read() count = ", count, "err = ", err)
				oneByte = make([]byte, 1024)
				count, err = nbr.Read(oneByte)
				log.Println("STDOUT: count = ", count)
				f.Write(oneByte[0:count])
				f.Flush()
				empty = true
				dead = true
			}
			// f.Write([]byte(string(oneByte)))
			if (count == 0) {
				empty = true
			} else {
				log.Println("STDOUT: count = ", count)
				io.WriteString(&buf, string(oneByte))
				f.Write(oneByte[0:count])
				f.Flush()
			}
		}

		if (totalCount != 0) {
			log.Println("STDOUT: totalCount = ", totalCount)
			totalCount = 0
		}
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
		return -1, "", err
	}
	// log.Println("shell.Run() END buf =", buf)
	// log.Println("shell.Run() END string(buf) =", string(buf))
	// log.Println("shell.Run() END buf.String() =", buf.String())
	// log.Println("shell.Run() END string(buf.Bytes()) =", string(buf.Bytes()))
	log.Println("shell.Run() END   ", cmdline)
	return 0, buf.String(), fmt.Errorf("") // nothing to do
}

func Daemon(cmdline string, timeout time.Duration) int {
	for {
		Run(cmdline)
		time.Sleep(timeout)
	}
}

// pass in two file handles (1 read, 1 write)
func NonBlockingReader(buffReader *bufio.Reader, writeFileHandle *os.File) {
	// newreader := bufio.NewReader(readFileHandle)

	// create a nonblocking GO reader
        nbr := nbreader.NewNBReader(buffReader, 1024)

	for {
		// defer buffReader.Close()
		// defer writeFileHandle.Flush()
		defer writeFileHandle.Close()
		totalCount := 0
		for {
			oneByte := make([]byte, 1024)
			count, err := nbr.Read(oneByte)
			if (err != nil) {
				log.Println("count, err =", count, err)
				return
			}
			totalCount += count
			if (count == 0) {
				time.Sleep(time.Duration(msecDelay) * time.Millisecond)   // without this delay this will peg the CPU
				if (totalCount != 0) {
					log.Println("STDERR: totalCount = ", totalCount)
					totalCount = 0
				}
			} else {
				log.Println("STDERR: count = ", count)
				writeFileHandle.Write(oneByte[0:count])
			}
		}
	}
}
