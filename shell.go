package shell

import (
	"strings"
	"time"
	"os"
	"os/exec"
	"bufio"
	"io/ioutil"

	"go.wit.com/log"
	"github.com/svent/go-nbreader"
)


// TODO: look at https://github.com/go-cmd/cmd/issues/20
// use go-cmd instead here?

var callback func(interface{}, int)

var shellStdout *os.File
var shellStderr *os.File

var spewOn      bool = false
var quiet       bool = false
// var msecDelay   int  = 20	// number of milliseconds to delay between reads with no data

// var bytesBuffer bytes.Buffer
// var bytesSplice []byte

func handleError(c interface{}, ret int) {
	log.Log(INFO, "shell.Run() Returned", ret)
	if (callback != nil) {
		callback(c, ret)
	}
}

func init() {
	callback  = nil
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
		log.Log(INFO, "LINE:", line)
		time.Sleep(1)
		Run(line)
	}
	return 0
}

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

func Unlink(filename string) {
	os.Remove(Path(filename))
}

func RM(filename string) {
	os.Remove(Path(filename))
}

/*
	err := process.Wait()

	if err != nil {
		if (spewOn) {
			// this panics: spew.Dump(err.(*exec.ExitError))
			spew.Dump(process.ProcessState)
		}
		// stuff := err.(*exec.ExitError)
		log.Log(INFO, "ERROR ", err.Error())
		log.Log(INFO, "END   ", cmdline)
		handleError(err, -1)
		return ""
*/

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
				log.Log(INFO, "count, err =", count, err)
				handleError(err, -1)
				return
			}
			totalCount += count
			if (count == 0) {
				time.Sleep(time.Duration(msecDelay) * time.Millisecond)   // without this delay this will peg the CPU
				if (totalCount != 0) {
					log.Log(INFO, "STDERR: totalCount = ", totalCount)
					totalCount = 0
				}
			} else {
				log.Log(INFO, "STDERR: count = ", count)
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
	log.Log(INFO, "shell.Run() START " + cmdline)

	cmd            := Chomp(cmdline) // this is like 'chomp' in perl
	cmdArgs        := strings.Fields(cmd)

	process        := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	process.Stderr  = os.Stderr
	process.Stdin   = os.Stdin
	process.Stdout  = os.Stdout
	process.Start()
	err := process.Wait()
	log.Log(INFO, "shell.Exec() err =", err)
	os.Exit(0)
}

// return true if the filename exists (cross-platform)
func Exists(filename string) bool {
	_, err := os.Stat(Path(filename))
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// return true if the filename exists (cross-platform)
func Dir(dirname string) bool {
	info, err := os.Stat(Path(dirname))
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// Cat a file into a string
func Cat(filename string) string {
	buffer, err := ioutil.ReadFile(Path(filename))
	// log.Log(INFO, "buffer =", string(buffer))
	if err != nil {
		return ""
	}
	return Chomp(buffer)
}
