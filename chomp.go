package shell

/* 
	perl 'chomp'

	send it anything, always get back a string
*/

import "log"
import "fmt"
import "reflect"
import "strings"
import "bytes"

// import "github.com/davecgh/go-spew/spew"

func chompBytesBuffer(buf *bytes.Buffer) string {
	var bytesSplice []byte
	bytesSplice = buf.Bytes()

	return Chomp(string(bytesSplice))
}

//
// TODO: obviously this is stupidly wrong
// TODO: fix this to trim fucking everything
// really world? 8 fucking years of this language
// and I'm fucking writing this? jesus. how the
// hell is everyone else doing this? Why isn't
// this already in the strings package?
//
func perlChomp(s string) string {
	// lots of stuff in go moves around the whole block of whatever it is so lots of things are padded with NULL values
	s = strings.Trim(s, "\x00") // removes NULL (needed!)

	// TODO: christ. make some fucking regex that takes out every NULL, \t, ' ", \n, and \r
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\n")
	s = strings.TrimSuffix(s, "\r")
	s = strings.TrimSuffix(s, "\n")

	s = strings.TrimSpace(s)		// this is like 'chomp' in perl
	s = strings.TrimSuffix(s, "\n")		// this is like 'chomp' in perl
	return s
}

// TODO: fix this to chomp \n \r NULL \t and ' '
func Chomp(a interface{}) string {
	// switch reflect.TypeOf(a) {
	switch t := a.(type) {
		case string:
			var s string
			s = a.(string)
			return perlChomp(s)
		case []uint8:
			log.Printf("shell.Chomp() FOUND []uint8")
			var tmp []uint8
			tmp = a.([]uint8)

			s := string(tmp)
			return perlChomp(s)
		case uint64:
			log.Printf("shell.Chomp() FOUND []uint64")
			s := fmt.Sprintf("%d", a.(uint64))
			return perlChomp(s)
		case int64:
			log.Printf("shell.Chomp() FOUND []int64")
			s := fmt.Sprintf("%d", a.(int64))
			return perlChomp(s)
		case *bytes.Buffer:
			log.Printf("shell.Chomp() FOUND *bytes.Buffer")
			var tmp *bytes.Buffer
			tmp = a.(*bytes.Buffer)

			var bytesSplice []byte
			bytesSplice = tmp.Bytes()
			return Chomp(string(bytesSplice))
		default:
			tmp := fmt.Sprint("shell.Chomp() NO HANDLER FOR TYPE: %T", a)
			handleError(fmt.Errorf(tmp), -1)
			log.Printf("shell.Chomp() NEED TO MAKE CONVERTER FOR type =", reflect.TypeOf(t))
	}
	tmp := "shell.Chomp() THIS SHOULD NEVER HAPPEN"
	handleError(fmt.Errorf(tmp), -1)
	return ""
}
