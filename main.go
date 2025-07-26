package main

import (
	"fmt"
	"os"

	"github.com/bhorvath/ddclient/config"
	"github.com/bhorvath/ddclient/dns"
	"github.com/bhorvath/ddclient/ipaddress"
)

func main() {
	args := config.ParseArgs()
	cfgS := config.NewService(args)
	checkSave(args, cfgS)
	cfg := prepareConfigs(cfgS)

	ih := ipaddress.NewIpifyIPAddressHandler("https://api.ipify.org")
	dh, err := dns.NewPorkbunDNSHandler("https://api.porkbun.com", cfg)
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

func checkSave(args *config.Args, cfgS config.Service) {
	// If saving config then no other action is taken
	if args.Save {
		fmt.Println("Saving configuration to file")
		if err := cfgS.SaveConfig(); err != nil {
			fmt.Println("Error encountered while saving configuration: %v",
				err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
}

func prepareConfigs(cfgS config.Service) *config.App {
	cfg, err := cfgS.BuildConfig()
	if err != nil {
		fmt.Println("Error encountered while configuring application: %v",
			err.Error())
		os.Exit(1)
	}
	return cfg
}
