package shell

/* 
	send it anything, always get back an int
*/

// import "log"
// import "reflect"
// import "strings"
// import "bytes"
import "strconv"

func Int(s string) int {
	s = Chomp(s)
	i, err := strconv.Atoi(s)
	if (err != nil) {
		handleError(err, -1)
		return 0
	}
	return i
}

func Int64(s string) int64 {
	s = Chomp(s)
	i, err := strconv.Atoi(s)
	if (err != nil) {
		handleError(err, -1)
		return 0
	}
	return int64(i)
}
