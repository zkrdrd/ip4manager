package network_test

import (
	"ipaddresspackage/network"
	"log"
	"testing"
)

func TestAccounting(t *testing.T) {
	net, err := network.NewNetwork("172.16.0.0/16")
	if err != nil {
		log.Panic(err)
	}

	net.SetUsedIP()

}

var ()
