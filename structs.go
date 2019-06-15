package shell

import "io"
import "bufio"
import "bytes"
import "github.com/svent/go-nbreader"

var FileMap	map[string]*File

var readBufferSize int

type File struct {
	Name		string
	BufferSize	int
	FbytesBuffer	bytes.Buffer
	Fbytes		[]byte

	Fio		io.ReadCloser		// := process.StdoutPipe()
	Fbufio		*bufio.Reader		// := bufio.NewReader(pOUT)
	Fnbreader	*nbreader.NBReader	// := nbreader.NewNBReader(readOUT, 1024)
}

func FileCreate(f io.ReadCloser) *File {
	var newfile File

	newfile.Fio = f
	newfile.Fbufio = bufio.NewReader(f)
	newfile.Fnbreader = nbreader.NewNBReader(newfile.Fbufio, 1024)

	return &newfile
}
