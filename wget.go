package shell

/* 
	This simply parses the command line arguments using the default golang
	package called 'flag'. This can be used as a simple template to parse
	command line arguments in other programs.

	It puts everything in a 'config' Protobuf which I think is a good
	wrapper around the 'flags' package and doesn't need a whole mess of
	global variables
*/

import "io"
import "os"
import "bytes"
import "strings"
import "net/http"

/*
import "git.wit.com/wit/shell"
import "github.com/davecgh/go-spew/spew"
*/

func Wget(url string) (*bytes.Buffer) {
	buf := new(bytes.Buffer)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		handleError(err, -1)
		return nil
	}
	defer resp.Body.Close()

	buf.ReadFrom(resp.Body)
	return buf
}

func WgetToFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		handleError(err, -1)
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		handleError(err, -1)
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func Write(filepath string, data string) bool {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return false
	}
	defer out.Close()

	// Write the body to file
	// _, err = io.Copy(out, resp.Body)
	count, err := io.Copy(out, strings.NewReader(data))
	if err != nil {
		handleError(err, -1)
		return false
	}
	handleError(nil, int(count))
	return true
}
