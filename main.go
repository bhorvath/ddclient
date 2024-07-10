package main

import (
	"fmt"
	"os"

	"github.com/bhorvath/ddclient/ipaddress"
)

func main() {
	ih := ipaddress.NewIpifyIPAddressHandler()

	ip, err := ih.GetCurrent()
	if err != nil {
		fmt.Println("Error getting current IP address", err)
		os.Exit(1)
	}
	fmt.Println("Current IP address:", ip)
}
