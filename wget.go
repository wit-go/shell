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
import "fmt"
import "log"
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

	log.Printf("res.StatusCode: %d\n", resp.StatusCode)
	if (resp.StatusCode != 200) {
		handleError(fmt.Errorf(fmt.Sprint("%d", resp.StatusCode)), -1)
		return nil
	}

	buf.ReadFrom(resp.Body)
	return buf
}

func WgetToFile(filepath string, url string) error {
	log.Println("WgetToFile() filepath =", filepath)
	log.Println("WgetToFile() URL =", url)
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

// write out a file. Always be nice and end with '\n'
// if you are here and want to complain about ending in '\n'
// then you probably aren't going to like lots of things in this
// package. I will quote the evilwm man page:
//
// BUGS: The author's idea of friendly may differ to that of many other people.
//
func Write(filepath string, data string) bool {
	data = Chomp(data) + "\n"
	// Create the file
	out, err := os.Create(Path(filepath))
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
