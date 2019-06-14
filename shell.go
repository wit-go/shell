package shell

import "fmt"
import "strings"
import "time"
import "os"
import "os/exec"
import "bufio"
import "bytes"
import "io"
// import "io/ioutil"

import "github.com/davecgh/go-spew/spew"
import "github.com/svent/go-nbreader"

// import "log"
import log "github.com/sirupsen/logrus"
// TODO this journalhook to be cross platform
// import "github.com/wercker/journalhook"

// TODO: look at https://github.com/go-cmd/cmd/issues/20
// use go-cmd instead here?

var callback func(interface{}, int)

var shellStdout *os.File
var shellStderr *os.File

var spewOn      bool = false
var quiet       bool = false
var msecDelay   int  = 20	// number of milliseconds to delay between reads with no data

var bytesBuffer bytes.Buffer
var bytesSplice []byte

func handleError(c interface{}, ret int) {
	log.Debug("shell.Run() Returned", ret)
	if (callback != nil) {
		callback(c, ret)
	}
}

func init() {
	callback = nil
}

func InitCallback(f func(interface{}, int)) {
	callback = f
}

// this means it won't copy all the output to STDOUT
func Quiet(q bool) {
	quiet = q
}

func Script(cmds string) int {
	// split on new lines (while we are at it, handle stupid windows text files
	lines := strings.Split(strings.Replace(cmds, "\r\n", "\n", -1), "\n")

	for _, line := range lines {
		line = Chomp(line) // this is like 'chomp' in perl
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

// NOTE: this might cause problems:
// always remove the newlines at the end ?
func Run(cmdline string) string {
	log.Println("shell.Run() START " + cmdline)

	cmd := Chomp(cmdline) // this is like 'chomp' in perl
	cmdArgs := strings.Fields(cmd)
	if (len(cmdArgs) == 0) {
		handleError(fmt.Errorf("cmdline == ''"), 0)
		log.Debug("END   ", cmd)
		return "" // nothing to do
	}
	if (cmdArgs[0] == "cd") {
		if (len(cmdArgs) > 1) {
			log.Println("os.Chdir()", cmd)
			os.Chdir(cmdArgs[1])
		}
		handleError(nil, 0)
		log.Debug("END   ", cmd)
		return "" // nothing to do
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
	go nonBlockingReader(tmp, shellStderr, f)

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
				log.Debug("Read() count = ", count, "err = ", err)
				oneByte = make([]byte, 1024)
				count, err = nbr.Read(oneByte)
				log.Debug("STDOUT: count = ", count)
				if (quiet == false) {
					f.Write(oneByte[0:count])
					f.Flush()
				}
				empty = true
				dead = true
			}
			// f.Write([]byte(string(oneByte)))
			if (count == 0) {
				empty = true
			} else {
				log.Debug("STDOUT: count = ", count)
				io.WriteString(&bytesBuffer, string(oneByte))
				if (quiet == false) {
					f.Write(oneByte[0:count])
					f.Flush()
				}
			}
		}

		if (totalCount != 0) {
			log.Debug("STDOUT: totalCount = ", totalCount)
			totalCount = 0
		}
	}

	err := process.Wait()

	if err != nil {
		if (spewOn) {
			// this panics: spew.Dump(err.(*exec.ExitError))
			spew.Dump(process.ProcessState)
		}
		// stuff := err.(*exec.ExitError)
		log.Debug("ERROR ", err.Error())
		log.Debug("END   ", cmdline)
		handleError(err, -1)
		return ""
	}

	// log.Println("shell.Run() END buf =", bytesBuffer)
	// convert this to a byte array and then trip NULLs
	// WTF this copies nulls with b.String() is fucking insanly stupid
	byteSlice := bytesBuffer.Bytes()
	b := bytes.Trim(byteSlice, "\x00")

	log.Debug("shell.Run() END b =", b)

	// reset the bytesBuffer
	bytesBuffer.Reset()

	// NOTE: this might cause problems:
	// this removes the newlines at the end
	tmp2 := string(b)
	tmp2  = strings.TrimSuffix(tmp2, "\n")
	handleError(nil, 0)
	log.Println("shell.Run() END   ", cmdline)
	return Chomp(b)
}

func Daemon(cmdline string, timeout time.Duration) int {
	for {
		Run(cmdline)
		time.Sleep(timeout)
	}
}

// pass in two file handles (1 read, 1 write)
func nonBlockingReader(buffReader *bufio.Reader, writeFileHandle *os.File, stdout *bufio.Writer) {
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
				log.Debug("count, err =", count, err)
				handleError(err, -1)
				return
			}
			totalCount += count
			if (count == 0) {
				time.Sleep(time.Duration(msecDelay) * time.Millisecond)   // without this delay this will peg the CPU
				if (totalCount != 0) {
					log.Debug("STDERR: totalCount = ", totalCount)
					totalCount = 0
				}
			} else {
				log.Debug("STDERR: count = ", count)
				writeFileHandle.Write(oneByte[0:count])
				if (quiet == false) {
					stdout.Write(oneByte[0:count])
					stdout.Flush()
				}
			}
		}
	}
}

// run something and never return from it
// TODO: pass STDOUT, STDERR, STDIN correctly
// TODO: figure out how to nohup the process and exit
func Exec(cmdline string) {
	log.Println("shell.Run() START " + cmdline)

	cmd     := Chomp(cmdline) // this is like 'chomp' in perl
	cmdArgs := strings.Fields(cmd)

	process := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	process.Start()
	err := process.Wait()
	log.Println("shell.Exec() err =", err)
	os.Exit(0)
}
