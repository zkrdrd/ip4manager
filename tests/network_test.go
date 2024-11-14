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

	ip1, _ := net.GetFreeIP()
	ip2, _ := net.GetFreeIP()
	ip3, err := net.GetFreeIP()
	fmt.Println("1", ip1)
	fmt.Println("2", ip2)
	fmt.Println("3", ip3, err)
	fmt.Println(net.SetUsedIP(ip1))
	fmt.Println(net.SetUsedIP("192.168.0.6"))
	net.ReleaseIP(ip3)
	net.ReleaseIP(ip1)
	ip4, _ := net.GetFreeIP()
	ip5, _ := net.GetFreeIP()
	fmt.Println(ip4, " ", ip5)
}

var ()
