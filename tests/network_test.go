package network_test

import (
	"ipaddresspackage/network"
	"ipaddresspackage/storage/memory"
	"log"
	"testing"
)

func TestAccounting(t *testing.T) {
	net, err := network.NewNetwork(172, 16, 0, 0, 128, 0, 0, 0)
	if err != nil {
		log.Panic(err)
	}
	// fmt.Println(net)
	//fmt.Println(cidr.AddressCount(&net))

	networkMapping := memory.NewDB()
	networkMapping.TakeIPAddress(net)

}

var ()
