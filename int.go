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
