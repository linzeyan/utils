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

/* HasNullByteInFile checks the file has the ASCII 0 or not. */
func HasNullByteInFile(filePath string) bool {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return true
	}
	for _, b := range f {
		if b == 0 {
			return true
		}
	}
	return false
}

/* HasNullByteInReader checks the reader has the ASCII 0 or not. */
func HasNullByteInReader(r io.Reader) bool {
	data, err := io.ReadAll(r)
	if err != nil {
		return true
	}
	for _, b := range data {
		if b == 0 {
			return true
		}
	}
	return false
}

/* IsDomain checks if i is a domain. */
func IsDomain(i any) bool {
	const elements = "~!@#$%^&*()_+`={}|[]\\:\"<>?,/"
	if val, ok := i.(string); ok {
		if strings.ContainsAny(val, elements) {
			return false
		}
		slice := strings.Split(val, ".")
		l := len(slice)
		if l > 1 {
			n, err := strconv.Atoi(slice[l-1])
			if err != nil {
				return true
			}
			s := strconv.Itoa(n)
			return slice[l-1] != s
		}
	}
	return false
}

/* IsPathExist checks if f is a valid path. */
func IsPathExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

/* IsIP checks if i is an IP address. */
func IsIP(i string) bool {
	ip, err := netip.ParseAddr(i)
	if err != nil {
		return false
	}
	return ip.IsValid()
}

/* IsCIDR checks if i is a valid CIDR. */
func IsCIDR(i string) bool {
	ip, err := netip.ParsePrefix(i)
	if err != nil {
		return false
	}
	return ip.IsValid()
}

/* IsIPv4 checks if i is an ipv4 address. */
func IsIPv4(i string) bool {
	ip, err := netip.ParseAddr(i)
	if err != nil {
		return false
	}
	return ip.Is4()
}

/* IsIPv4CIDR checks if i is a valid IPv4 CIDR. */
func IsIPv4CIDR(i string) bool {
	ip, err := netip.ParsePrefix(i)
	if err != nil {
		return false
	}
	return ip.IsValid() && ip.Addr().Is4()
}

/* IsIPv6 checks if i is an ipv6 address. */
func IsIPv6(i string) bool {
	ip, err := netip.ParseAddr(i)
	if err != nil {
		return false
	}
	return ip.Is6()
}

/* IsIPv6CIDR checks if i is a valid IPv6 CIDR. */
func IsIPv6CIDR(i string) bool {
	ip, err := netip.ParsePrefix(i)
	if err != nil {
		return false
	}
	return ip.IsValid() && ip.Addr().Is6()
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

/* IsDarwin checks if host is windows. */
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
