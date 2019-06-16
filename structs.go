package shell

import "io"
import "os/exec"
import "bufio"
import "bytes"
import "github.com/svent/go-nbreader"

var FileMap	map[string]*File

var readBufferSize int

type File struct {
	Name		string
	// BufferSize	int
	// Buffer		*bytes.Buffer
	// Fbytes		[]byte
	TotalCount	int
	Empty		bool
	Dead		bool

	Fio		io.ReadCloser		// := process.StdoutPipe()
	Fbufio		*bufio.Reader		// := bufio.NewReader(pOUT)
	Fnbreader	*nbreader.NBReader	// := nbreader.NewNBReader(readOUT, 1024)
}

type Shell struct {
	Cmdline		string
	Process		*exec.Cmd
	Done		bool
	Quiet		bool
	Fail		bool
	Error		error
	Buffer		*bytes.Buffer

	// which names are really better here?
	// for now I init them both to test out
	// how the code looks and feels
	STDOUT		*File
	STDERR		*File
	Stdout		*File
	Stderr		*File
}

// default values for Shell
func New() *Shell {
	var tmp Shell

	tmp.Done = false
	tmp.Fail = false
	tmp.Quiet = quiet

	return &tmp
}
