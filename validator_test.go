package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestHasNullByteInFile(t *testing.T) {
	createDir(testDir)
	testCases := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"1", []byte{0, 1, 2, 3}, true},
		{"2", []byte{44, 55, 66, 77, 88}, false},
		{"3", []byte{111, 222}, false},
	}
	for _, testCase := range testCases {
		file := filepath.Join(testDir, "test.txt")
		err := os.WriteFile(file, testCase.data, os.ModePerm)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		b, err := HasNullByteInFile(file)
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.expected {
				assert.True(t, b)
				assert.NoError(t, err)
				return
			}
			assert.False(t, b)
		})
	}
	_ = os.RemoveAll(testDir)
}

func TestHasNullByteInReader(t *testing.T) {
	createDir(testDir)
	testCases := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"1", []byte{0, 1, 2, 3}, true},
		{"2", []byte{44, 55, 00, 77, 88}, true},
	}
	for _, testCase := range testCases {
		file := filepath.Join(testDir, "test.txt")
		err := os.WriteFile(file, testCase.data, os.ModePerm)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		t.Run(testCase.name, func(t *testing.T) {
			f, err := os.Open(file)
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}
			defer f.Close()
			b, err := HasNullByteInReader(f)
			if testCase.expected {
				assert.True(t, b)
				assert.NoError(t, err)
				return
			}
			assert.False(t, b)
		})
	}
	_ = os.RemoveAll(testDir)
}

func TestIsDomain(t *testing.T) {
	testCases := []struct {
		input    any
		expected bool
	}{
		{"1.1.1.1", false},
		{"example.com", true},
		{"Hello world", false},
		{11111, false},
		{"dns-admin.google.com", true},
		{"dns-admin.google.com.", true},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", testCase.input), func(t *testing.T) {
			if testCase.expected {
				assert.True(t, IsDomain(testCase.input))
				return
			}
			assert.False(t, IsDomain(testCase.input))
		})
	}
}

func TestIsPathExist(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1", false},
		{"vendor", true},
		{"command.json", false},
		{"root_test.go", false},
		{"/dev/null", !IsWindows()},
	}
	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			if testCase.expected {
				assert.True(t, IsPathExist(testCase.input))
				return
			}
			assert.False(t, IsPathExist(testCase.input))
		})
	}
}

func TestIsIP(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1", true},
		{"999.1.1.1", false},
		{"260.2.3.4", false},
		{"example.com", false},
		{"2404:6800:4008:c01::65", true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			if testCase.expected {
				assert.True(t, IsIP(testCase.input))
				return
			}
			assert.False(t, IsIP(testCase.input))
		})
	}
}

func TestIsCIDR(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1/32", true},
		{"999.1.1.1", false},
		{"260.2.3.4/24", false},
		{"example.com/8", false},
		{"2404:6800:4008:c01::65/32", true},
		{"fe80::aede:48ff:fe00:1122/64", true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			if testCase.expected {
				assert.True(t, IsCIDR(testCase.input))
				return
			}
			assert.False(t, IsCIDR(testCase.input))
		})
	}
}

func TestIsIPv4(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1", true},
		{"999.1.1.1", false},
		{"260.2.3.4", false},
		{"example.com", false},
		{"2404:6800:4008:c01::65", false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			if testCase.expected {
				assert.True(t, IsIPv4(testCase.input))
				return
			}
			assert.False(t, IsIPv4(testCase.input))
		})
	}
}

func TestIsIPv4CIDR(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1/24", true},
		{"999.1.1.1/8", false},
		{"260.2.3.4/12", false},
		{"1.2.3.4/10", true},
		{"example.com/22", false},
		{"2404:6800:4008:c01::65/64", false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			if testCase.expected {
				assert.True(t, IsIPv4CIDR(testCase.input))
				return
			}
			assert.False(t, IsIPv4CIDR(testCase.input))
		})
	}
}

func TestIsIPv6(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1", false},
		{"999.1.1.1", false},
		{"260.2.3.4", false},
		{"example.com", false},
		{"2404:6800:4008:c01::65", true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			if testCase.expected {
				assert.True(t, IsIPv6(testCase.input))
				return
			}
			assert.False(t, IsIPv6(testCase.input))
		})
	}
}

func TestIsIPv6CIDR(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1/3", false},
		{"999.1.1.1/4", false},
		{"260.2.3.4/5", false},
		{"example.com/7", false},
		{"2404:6800:4008:c01::65/10", true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			if testCase.expected {
				assert.True(t, IsIPv6CIDR(testCase.input))
				return
			}
			assert.False(t, IsIPv6CIDR(testCase.input))
		})
	}
}

func TestIsURL(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1", false},
		{"999.1.1.1", false},
		{"2404:6800:4008:c01::65", false},
		{"https://1.1.1.1", true},
		{"example.com", false},
		{"https://example.com", true},
		{"https://example.com/?", true},
		{"https://example.com/api/v1/add?user=1", true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			if testCase.expected {
				assert.True(t, IsURL(testCase.input))
				return
			}
			assert.False(t, IsURL(testCase.input))
		})
	}
}

func TestIsDarwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		assert.True(t, IsDarwin())
		return
	}
	assert.False(t, IsDarwin())
}

func TestIsWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		assert.True(t, IsWindows())
		return
	}
	assert.False(t, IsWindows())
}
