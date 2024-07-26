package main

import (
	"fmt"
	"os"

	"github.com/bhorvath/ddclient/dns"
	"github.com/bhorvath/ddclient/ipaddress"
)

func main() {
	ih := ipaddress.NewIpifyIPAddressHandler()
	dh, err := dns.NewPorkbunDNSHandler("https://api.porkbun.com")
	if err != nil {
		fmt.Println("Error setting up DNS handler:", err)
		os.Exit(1)
	}

	ip, err := ih.GetCurrent()
	if err != nil {
		fmt.Println("Error getting current IP address:", err)
		os.Exit(1)
	}
	fmt.Println("Current IP address:", ip)

	err = dh.Update(ip)
	if err != nil {
		fmt.Println("Error updating DNS entry:", err)
		os.Exit(1)
	}
}
