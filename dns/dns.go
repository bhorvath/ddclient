package dns

import "net/netip"

type DNSHandler interface {
	Update(netip.Addr) error
}
