package ipaddress

import "net/netip"

type IPAddressHandler interface {
	GetCurrent() (netip.Addr, error)
}
