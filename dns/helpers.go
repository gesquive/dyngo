package dns

import (
	"net"
	"strings"
)

// SplitDomainRecord splits a record into a domain and record name
func SplitDomainRecord(domainRecord string) (domain string, record string) {
	domainParts := strings.Split(domainRecord, ".")
	if len(domainParts) > 2 {
		// sub.domain.net => domain.net
		domain = strings.Join(domainParts[len(domainParts)-2:], ".")
		record = strings.Join(domainParts[:len(domainParts)-2], ".")
	} else {
		domain = domainRecord
		record = "@"
	}
	return
}

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
