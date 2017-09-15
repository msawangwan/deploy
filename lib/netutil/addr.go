package netutil

import (
	"net"
	"strings"
)

// LocalIP parses an ip given an interface name (linux only)
func LocalIP(ifname string) (string, error) {
	intfs, e := net.Interfaces()

	if e != nil {
		return "none", e
	}

	for _, intf := range intfs {
		if strings.Contains(intf.Name, ifname) {
			addrs, e := intf.Addrs()

			if e != nil {
				return "none", e
			}

			for _, addr := range addrs {
				addrstr := addr.String()
				if !strings.Contains(addrstr, "[") {
					return strings.Split(addrstr, "/")[0], nil
				}
			}
		}
	}

	return "none", nil
}
