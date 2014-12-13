package cmd

import (
	"fmt"
	"runtime"
)

var gitsha string

func init() {
	AddPlugin("Version", "(?i)^\\.version?$", MessageHandler(GetVersion), false, false)
}

func GetVersion(msg *Message) {
	version := runtime.Version()
	out := fmt.Sprintf("[Bogon] %s running on %s", gitsha, version)
	msg.Return(out)
}
