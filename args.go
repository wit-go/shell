package shell

// initializes logging and command line options

import (
	"go.wit.com/log"
)

var INFO log.LogFlag

func init() {
	INFO.B = false
	INFO.Name = "INFO"
	INFO.Subsystem = "shell"
	INFO.Short = "shell"
	INFO.Desc = "general info"
	INFO.Register()
}
