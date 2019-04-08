package shell

import "fmt"
import "log"
import "strings"
import "time"
import "os"
import "os/exec"
import "bufio"
import "github.com/davecgh/go-spew/spew"
import "github.com/svent/go-nbreader"

func Script(cmds string) int {
	// split on new lines (while we are at it, handle stupid windows text files
	lines := strings.Split(strings.Replace(cmds, "\r\n", "\n", -1), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line) // this is like 'chomp' in perl
		fmt.Println("LINE:", line)
		time.Sleep(1)
		Run(line)
	}
	return 0
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
	stdout, _ := process.StdoutPipe()
	stderr, _ := process.StderrPipe()
	process.Start()

	f := bufio.NewWriter(os.Stdout)

	newreader := bufio.NewReader(stdout)
	nbr := nbreader.NewNBReader(newreader, 1024)

	newerrreader := bufio.NewReader(stderr)
	nbrerr := nbreader.NewNBReader(newerrreader, 1024)

	for {
		time.Sleep(2 * time.Millisecond)   // only check the buffer 500 times a second
		// log.Println("sleep done")

		oneByte := make([]byte, 1024)
		count, err := nbr.Read(oneByte)

		if (err != nil) {
			// log.Println("Read() count = ", count, "err = ", err)
			oneByte = make([]byte, 1024)
			count, err = nbr.Read(oneByte)
			f.Write([]byte(string(oneByte)))
			f.Flush()
		}
		f.Write([]byte(string(oneByte)))
		f.Flush()

		oneByte = make([]byte, 1024)
		count, err = nbrerr.Read(oneByte)

		if (err != nil) {
			oneByte = make([]byte, 1024)
			count, err = nbrerr.Read(oneByte)
			f.Write([]byte(string(oneByte)))
			f.Flush()

			log.Println("Read() count = ", count, "err = ", err)
			spew.Dump(process.Process)
			spew.Dump(process.ProcessState)
			err := process.Wait()
			if err != nil {
				spew.Dump(err.(*exec.ExitError))
				spew.Dump(process.ProcessState)
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
		spew.Dump(err.(*exec.ExitError))
		spew.Dump(process.ProcessState)
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
