package utils

import (
	"io"
	"net/netip"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
)

/* HasNullByteInFile checks the byte slice has the ASCII 0 or not. */
func HasNullByte(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return true
		}
	}
	return false
}

/* HasNullByteInFile checks the file has the ASCII 0. */
func HasNullByteInFile(filePath string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		logPanic(err)
	}
	return HasNullByte(data)
}

/* HasNullByteInReader checks the reader has the ASCII 0. */
func HasNullByteInReader(r io.Reader) bool {
	data, err := io.ReadAll(r)
	if err != nil {
		logPanic(err)
	}
	return HasNullByte(data)
}

/* IsDomain checks if i is a valid domain. */
func IsDomain(i any) bool {
	const elements = "~!@#$%^&*()_+`={}|[]\\:\"<>?,/"
	val, ok := i.(string)
	if !ok {
		return false
	}
	if strings.ContainsAny(val, elements) {
		return false
	}

	raw := strings.Split(val, ".")
	rawLength := len(raw)
	if rawLength < 2 {
		return false
	}
	rawLast := raw[rawLength-1]

	var slice []string
	if rawLast == "" {
		slice = raw[:rawLength-1]
	} else {
		slice = raw
	}
	length := len(slice)
	if length < 2 {
		return false
	}
	last := slice[length-1]
	_, err := strconv.Atoi(last)
	return err != nil
}

/* IsPathExist checks if f is a valid path. */
func IsPathExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

/* IsIP checks if i is an IP address. */
func IsIP(i string) bool {
	ip, err := netip.ParseAddr(i)
	return err == nil && ip.IsValid()
}

/* IsCIDR checks if i is a valid CIDR. */
func IsCIDR(i string) bool {
	ip, err := netip.ParsePrefix(i)
	return err == nil && ip.IsValid()
}

/* IsIPv4 checks if i is an ipv4 address. */
func IsIPv4(i string) bool {
	ip, err := netip.ParseAddr(i)
	return err == nil && ip.Is4()
}

/* IsIPv4CIDR checks if i is a valid IPv4 CIDR. */
func IsIPv4CIDR(i string) bool {
	ip, err := netip.ParsePrefix(i)
	return err == nil && ip.IsValid() && ip.Addr().Is4()
}

/* IsIPv6 checks if i is an ipv6 address. */
func IsIPv6(i string) bool {
	ip, err := netip.ParseAddr(i)
	return err == nil && ip.Is6()
}

/* IsIPv6CIDR checks if i is a valid IPv6 CIDR. */
func IsIPv6CIDR(i string) bool {
	ip, err := netip.ParsePrefix(i)
	return err == nil && ip.IsValid() && ip.Addr().Is6()
}

/* IsURL checks if u is a valid url. */
func IsURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}

/* IsDarwin checks if host is darwin. */
func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

/* IsLinux checks if host is linux. */
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

/* IsDarwin checks if host is windows. */
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
