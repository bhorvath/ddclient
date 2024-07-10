package ipaddress

import "net/netip"

type MockIPAddressHandler struct{}

func NewMockIPAddressHandler() *MockIPAddressHandler {
	return &MockIPAddressHandler{}
}

func (h *MockIPAddressHandler) GetCurrent() (netip.Addr, error) {
	return netip.ParseAddr("127.0.0.1")
}
