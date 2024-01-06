package shell

// initializes logging and command line options

import (
	"go.wit.com/log"
)

var INFO log.LogFlag
var RUN log.LogFlag
var SSH log.LogFlag

func init() {
	INFO.B = false
	INFO.Name = "INFO"
	INFO.Subsystem = "shell"
	INFO.Short = "shell"
	INFO.Desc = "general info"
	INFO.Register()

	RUN.B = false
	RUN.Name = "RUN"
	RUN.Subsystem = "shell"
	RUN.Short = "shell"
	RUN.Desc = "Run() info"
	RUN.Register()

	SSH.B = false
	SSH.Name = "SSH"
	SSH.Subsystem = "shell"
	SSH.Short = "shell"
	SSH.Desc = "ssh() info"
	SSH.Register()
}
