package ipaddress

import (
	"io"
	"net/http"
	"net/netip"
)

type IpifyIPAddressHandler struct {
	baseURL string
}

func NewIpifyIPAddressHandler() *IpifyIPAddressHandler {
	return &IpifyIPAddressHandler{baseURL: "https://api.ipify.org"}
}

func (h *IpifyIPAddressHandler) GetCurrent() (netip.Addr, error) {
	res, err := http.Get(h.baseURL)
	if err != nil {
		return netip.Addr{}, err
	}

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return netip.Addr{}, err
	}

	ip, err := netip.ParseAddr(string(resBody))

	if err != nil {
		return netip.Addr{}, err
	}

	return ip, nil
}