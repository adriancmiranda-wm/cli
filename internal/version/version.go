package version

import (
	"runtime/debug"
)

var Version = "unknown"

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		panic("failed to read build info, < go version 1.22")
	}

	mainVersion := info.Main.Version
	if mainVersion == "" || mainVersion == "(devel)" {
		return
	}

	Version = mainVersion
}
