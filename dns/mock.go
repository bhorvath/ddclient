package dns

import (
	"fmt"
	"net/netip"
)

type MockDNSHandler struct{}

func NewMockDNSHandler() *MockDNSHandler {
	return &MockDNSHandler{}
}

func (h *MockDNSHandler) Update(ip netip.Addr) error {
	fmt.Println("Updated mock DNS")

	return nil
}
