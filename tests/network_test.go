package network_test

import (
	"fmt"
	"ipaddresspackage/network"
	"log"
	"testing"
)

func TestAccounting(t *testing.T) {
	net, err := network.NewNetwork("172.16.0.0/16")
	if err != nil {
		log.Panic(err)
	}

	ip1, _ := net.SetUsedIP()
	ip2, _ := net.SetUsedIP()
	ip3, err := net.SetUsedIP()
	fmt.Println("1", ip1)
	fmt.Println("2", ip2)
	fmt.Println("3", ip3, err)
	net.ReleaseIP(ip3)
	net.ReleaseIP(ip1)
	ip4, _ := net.SetUsedIP()
	ip5, _ := net.SetUsedIP()
	fmt.Println(ip4, " ", ip5)
}

var ()
