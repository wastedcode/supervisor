package supervisor

import (
	"net/http"
	"net/netip"
	"strings"
)

var headerTrueClientIP = http.CanonicalHeaderKey("True-Client-IP")
var headerXRealIP = http.CanonicalHeaderKey("X-Real-IP")
var headerXForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")

// GetIPFromHeaders returns the IP address from the request headers
// The request headers can be injected and ips can be spoofed
// They are also injected by proxies
func GetIPFromHeaders(r *http.Request) netip.Addr {
	ip := r.Header.Get(headerXForwardedFor)

	if ip == "" {
		ip = r.Header.Get(headerXRealIP)
		if ip == "" {
			ip = r.Header.Get(headerTrueClientIP)
		}
	}

	if ip == "" {
		return netip.IPv4Unspecified()
	}

	ips := strings.Split(ip, ",")
	ipaddr, err := netip.ParseAddr(ips[0])
	if err != nil {
		return netip.IPv4Unspecified()
	}
	return ipaddr
}
