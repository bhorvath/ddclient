package main

import (
	"fmt"

	"github.com/bhorvath/ddclient/ipaddress"
)

func main() {
	ih := ipaddress.NewMockIPAddressHandler()

	ip, err := ih.GetCurrent()
	if err != nil {
		fmt.Println("Error getting current IP address", err)
	}
	fmt.Println("current IP address:", ip)
}
