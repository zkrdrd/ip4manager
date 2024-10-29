package network_test

import (
	"fmt"
	"ipaddresspackage/pkg/network"
	"testing"
)

func TestAccounting(t *testing.T) {
	net, err := network.NewNetwork(`192.168.0.0`, `255.255.0.0`)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(net)
	fmt.Println(network.AddressCount(net))
	first, second := network.AddressRange(net)
	fmt.Println(first, second)
}

var ()
