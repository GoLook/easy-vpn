package main

import (
	"fmt"
	"log"
	"time"

	"github.com/JamesClonk/easy-vpn/provider"
)

func getAllVMs(p provider.API) []provider.VM {
	machines, err := p.GetAllVMs()
	if err != nil {
		log.Println("Could not retrieve list of virtual machines")
		log.Fatal(err)
	}
	return machines
}

func getEasyVpnVM(p provider.API, sshkeyId string) (vm provider.VM) {
	cfg := p.GetConfig()
	os := cfg.Providers[cfg.Provider].OS
	size := cfg.Providers[cfg.Provider].Size
	region := cfg.Providers[cfg.Provider].Region

	// check to see if easy-vpn vm already exists
	vmExists := false
	for _, machine := range getAllVMs(p) {
		if machine.Name == EASYVPN_IDENTIFIER {
			vm = machine
			vmExists = true
			break
		}
	}

	// if it already exists, make sure its up and running
	if vmExists {
		fmt.Println("Virtual machine already exists")
		statusOfVM(p, &vm)
	} else { // otherwise, create a new vm and start it
		fmt.Println("Create new virtual machine")

		_, err := p.CreateVM(EASYVPN_IDENTIFIER, os, size, region, sshkeyId)
		if err != nil {
			log.Println("Could not create new virtual machine")
			log.Fatal(err)
		}

		waitForNewVM(p, &vm)
		statusOfVM(p, &vm)
	}
	fmt.Println()

	return
}

func waitForNewVM(p provider.API, vm *provider.VM) {
	ticker := ticker()

	// TODO: maybe have some maximum waiting/polling time for doing a timeout
POLL:
	for {
		for _, machine := range getAllVMs(p) {
			if machine.Name == EASYVPN_IDENTIFIER {
				vm.Id = machine.Id
				vm.Name = machine.Id
				vm.Status = machine.Status
				vm.IP = machine.IP
				vm.OS = machine.OS
				vm.Region = machine.Region

				ticker.Stop()
				break POLL
			}
		}
		time.Sleep(15 * time.Second)
	}

	return
}

func statusOfVM(p provider.API, vm *provider.VM) {
	ticker := ticker()

	// TODO: maybe have some maximum waiting/polling time for doing a timeout
POLL:
	for {
		for _, machine := range getAllVMs(p) {
			if machine.Name == EASYVPN_IDENTIFIER && machine.Status == "active" {
				vm.Status = machine.Status

				ticker.Stop()
				break POLL
			}
		}
		time.Sleep(15 * time.Second)
	}
}

func ticker() *time.Ticker {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Print(".")
		}
	}()
	return ticker
}
