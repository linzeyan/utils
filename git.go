package utils

import (
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

/*
Version returns version information string containing git tag,
git commit hash, build time, runtime and pkg path.
*/
func Version() string {
	var (
		revision   = "unknown"
		version    = "unknown"
		dirtyBuild = true
		commitTime time.Time
	)
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	if info.Main.Version != "" {
		version = info.Main.Version
	}
	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			revision = kv.Value
		case "vcs.time":
			commitTime, _ = time.Parse(time.RFC3339, kv.Value)
		case "vcs.modified":
			dirtyBuild = kv.Value == "true"
		}
	}

	buf := new(strings.Builder)
	buf.WriteString("Version: " + version + "\n")
	buf.WriteString("GitCommit: " + revision + "\n")
	buf.WriteString("Time: " + commitTime.UTC().Format(time.RFC3339) + "\n")
	buf.WriteString("Runtime: " + info.GoVersion + " " + runtime.GOOS + "/" + runtime.GOARCH + "\n")
	buf.WriteString("Path: " + info.Path)
	if dirtyBuild {
		buf.WriteString("\ndirty build!")
	}
	return buf.String()
}
