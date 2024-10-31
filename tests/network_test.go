package network_test

import (
	"fmt"
	"ipaddresspackage/network"
	"testing"
)

func TestAccounting(t *testing.T) {
	net, err := network.NewNetwork(`192.168.0.0`, `255.255.0.0`)
	if err != nil {
		fmt.Print(err)
	}
	//fmt.Println(net)
	//fmt.Println(cidr.AddressCount(&net))

	network.NewNetwrokMapping(net)
}

var ()
