package config

import (
	"os"
	"strings"
)

var (
	cmdArgs map[string]string
)

func init() {
	cmdArgs = make(map[string]string)
	var name string
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-") {
			name = strings.TrimLeft(arg, "-")
		}
		if name != "" {
			cmdArgs[name] = arg
			name = ""
		}
	}
}
