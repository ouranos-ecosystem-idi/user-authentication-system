package authentication

import (
	"net"
)

// Cidr
// Summary: This is structure which defines the CIDR model.
type Cidr struct {
	Cidr string `json:"cidr"`
}

// Cidrs
// Summary: This is structure which defines the slice of Cidr.
type Cidrs []*Cidr

// Contains
// Summary: This is the function which checks whether the IP address exists in this struct slice.
// input: ip(string): IP address
// output: (bool) true if the IP address exists in this slice, false otherwise
func (ms Cidrs) Contains(ip string) bool {
	for _, cidr := range ms {
		_, subnet, _ := net.ParseCIDR(cidr.Cidr)
		if subnet.Contains(net.ParseIP(ip)) {
			return true
		}
	}
	return false
}
