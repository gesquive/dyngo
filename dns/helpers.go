package dns

import (
	"net"
)

// IsValidIPv4Addr returns true if the given address is a valid IPv4 address
func IsValidIPv4Addr(addr string) bool {
	ip := net.ParseIP(addr)
	return ip.To4() != nil
}

// IsValidIPv6Addr returns true if the given address is a valid IPv6 address
func IsValidIPv6Addr(addr string) bool {
	ip := net.ParseIP(addr)
	return ip.To16() != nil
}
