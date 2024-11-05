package network_test

import (
	"fmt"
	"ipaddresspackage/network"
	"testing"
)

func TestAccounting(t *testing.T) {
	net := network.NewNetwork(172, 16, 0, 0, 255, 255, 0, 0)
	fmt.Println(net)
	//fmt.Println(cidr.AddressCount(&net))

	network.NewNetwrokMapping(net)
}

var ()
