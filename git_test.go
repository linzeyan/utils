package utils

import (
	"runtime"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	var commitTime time.Time
	info, _ := debug.ReadBuildInfo()
	buf := new(strings.Builder)
	buf.WriteString("Version: unknown\n")
	buf.WriteString("GitCommit: unknown\n")
	buf.WriteString("Time: " + commitTime.UTC().Format(time.RFC3339) + "\n")
	buf.WriteString("Runtime: " + info.GoVersion + " " + runtime.GOOS + "/" + runtime.GOARCH + "\n")
	buf.WriteString("Path: " + "\ndirty build!")
	assert.Equal(t, buf.String(), Version())
}
