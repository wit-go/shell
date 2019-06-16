package shell

import "strings"
import "time"
import "os/exec"
import "bytes"
import "io"
import "fmt"
import "os"
import "bufio"

import "github.com/svent/go-nbreader"
// import "github.com/davecgh/go-spew/spew"

import "log"
// import "git.wit.com/wit/log"

var msecDelay int = 20 // check every 20 milliseconds

// TODO: look at https://github.com/go-cmd/cmd/issues/20
// use go-cmd instead here?
// exiterr.Sys().(syscall.WaitStatus)

// var newfile *shell.File

func Run(cmdline string) string {
	test := New()
	test.Exec(cmdline)
	return Chomp(test.Buffer)
}

func (cmd *Shell) Run(cmdline string) string {
	cmd.InitProcess(cmdline)
	if (cmd.Error != nil) {
		return ""
	}
	cmd.Exec(cmdline)
	return Chomp(cmd.Buffer)
}

func (cmd *Shell) InitProcess(cmdline string) {
	log.Println("shell.InitProcess() START " + cmdline)

	cmd.Cmdline = Chomp(cmdline) // this is like 'chomp' in perl
	cmdArgs := strings.Fields(cmd.Cmdline)
	if (len(cmdArgs) == 0) {
		cmd.Error = fmt.Errorf("cmdline == ''")
		cmd.Done = true
		return
	}
	if (cmdArgs[0] == "cd") {
		if (len(cmdArgs) > 1) {
			log.Println("os.Chdir()", cmd)
			os.Chdir(cmdArgs[1])
		}
		handleError(nil, 0)
		cmd.Done = true
		return
	}

	cmd.Process = exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
}

func (cmd *Shell) FileCreate(out string) {
	var newfile File

	var iof io.ReadCloser
	if (out == "STDOUT") {
		iof, _   = cmd.Process.StdoutPipe()
	} else {
		iof, _   = cmd.Process.StderrPipe()
	}

	newfile.Fio = iof
	newfile.Fbufio = bufio.NewReader(iof)
	newfile.Fnbreader = nbreader.NewNBReader(newfile.Fbufio, 1024)

	if (out == "STDOUT") {
		cmd.STDOUT = &newfile
	} else {
		cmd.STDERR = &newfile
	}
}

// NOTE: this might cause problems:
// always remove the newlines at the end ?
func (cmd *Shell) Exec(cmdline string) {
	log.Println("shell.Run() START " + cmdline)

	cmd.InitProcess(cmdline)
	if (cmd.Error != nil) {
		return
	}

	cmd.FileCreate("STDOUT")
	cmd.FileCreate("STDERR")

	cmd.Process.Start()

	// TODO; 'goroutine' both of these
	// and make your own wait that will make sure
	// the process is then done and run process.Wait()
	go cmd.Capture(cmd.STDERR)
	cmd.Capture(cmd.STDOUT)

	// wait until the process exists
	// https://golang.org/pkg/os/exec/#Cmd.Wait
	// What should happen here, before calling Wait()
	// is checks to make sure the READERS() on STDOUT and STDERR are done
	err := cmd.Process.Wait()

	// time.Sleep(2 * time.Second) // putting this here doesn't help STDOUT flush()

	if (err != nil) {
		cmd.Fail = true
		cmd.Error = err
		log.Println("process.Wait() END err =", err.Error())
	} else {
		log.Println("process.Wait() END")
	}
	return
}

// nonblocking read until file errors
func (cmd *Shell) Capture(f *File) {
	// log.Debugln("nbrREADER() START")

	if (cmd.Buffer == nil) {
		cmd.Buffer = new(bytes.Buffer)
	}
	if (cmd.Buffer == nil) {
		// log.Debugln("f.Buffer == nil")
		// log.Debugln("SHOULD DIE HERE")
		f.Dead = false
		cmd.Error = fmt.Errorf("could not make buffer")
		cmd.Done = true
	}

	f.Dead = false

	// loop that keeps trying to read from f
	for (f.Dead == false) {
		time.Sleep(time.Duration(msecDelay) * time.Millisecond)   // only check the buffer 500 times a second

		// set to false so it keeps retrying reads
		f.Empty = false

		// tight loop that reads 1024 bytes at a time until buffer is empty
		// 1024 is set in f.BufferSize
		for (f.Empty == false) {
			f.Empty = cmd.ReadToBuffer(f)
		}
	}
}

// returns true  if filehandle buffer is empty
func (cmd *Shell) ReadToBuffer(f *File) bool {
	// log.Debugln("ReadToBuffer() START")
	nbr := f.Fnbreader
	oneByte := make([]byte, 1024)
	if (nbr == nil) {
		// log.Debugln("ReadToBuffer() ERROR nbr is nil")
		f.Dead = true
		return true
	}
	count, err := nbr.Read(oneByte)
	f.TotalCount += count

	if (err != nil) {
		// log.Debugln("ReadToBuffer() file has closed with", err)
		// log.Debugln("ReadToBuffer() count = ", count, "err = ", err)
		f.Dead = true
		return true
	}
	if (count == 0) {
		// log.Debugln("ReadToBuffer() START count == 0 return true")
		return true
	}
	// log.Debugln("ReadToBuffer() count = ", count)
	// tmp := Chomp(oneByte)
	// log.Debugln("ReadToBuffer() tmp = ", tmp)
	io.WriteString(cmd.Buffer, strings.Trim(string(oneByte), "\x00"))
	return false
}
