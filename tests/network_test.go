package network_test

import (
	"fmt"
	"ipaddresspackage/network"
	"log"
	"testing"
)

func TestAccounting(t *testing.T) {
	net, err := network.NewNetwork(networkString)
	if err != nil {
		log.Panic(err)
	}

	for _, expect := range GetIPParam {
		ip, err := net.GetFreeIP()
		if expect.GetIP != ip {
			t.Error(fmt.Errorf(`result field %v != %v`, expect.GetIP, ip))
		}
		if expect.Error != err {
			t.Error(fmt.Errorf(`result field %v != %v`, expect.Error, err))
		}
	}

	for _, expect := range SetIPParam {
		err := net.SetUsedIP(expect.SetIP)
		if expect.Error != err {
			t.Error(fmt.Errorf(`result field %v != %v`, expect.Error, err))
		}
	}

	for _, expect := range ReleaseIPParam {
		err := net.ReleaseIP(expect.ReleaseIP)
		if expect.Error != err {
			t.Error(fmt.Errorf(`result field %v != %v`, expect.Error, err))
		}
	}

	for _, expect := range GetIPParamAfterRelease {
		ip, err := net.GetFreeIP()
		if expect.GetIP != ip {
			t.Error(fmt.Errorf(`result field %v != %v`, expect.GetIP, ip))
		}
		if expect.Error != err {
			t.Error(fmt.Errorf(`result field %v != %v`, expect.Error, err))
		}
	}

}

var (
	networkString = "172.16.0.0/16"
	GetIPParam    = []struct {
		GetIP string
		Error error
	}{
		{
			GetIP: "172.16.0.2",
			Error: nil,
		},
		{
			GetIP: "172.16.0.3",
			Error: nil,
		},
		{
			GetIP: "172.16.0.4",
			Error: nil,
		},
	}

	SetIPParam = []struct {
		SetIP string
		Error error
	}{
		{
			SetIP: "172.16.0.5",
			Error: nil,
		},
		{
			SetIP: "172.16.0.3",
			Error: network.ErrIPAddressIsUsed,
		},
		{
			SetIP: "192.168.0.6",
			Error: network.ErrIPADressIsNotIncludedInNetwork,
		},
	}

	ReleaseIPParam = []struct {
		ReleaseIP string
		Error     error
	}{
		{
			ReleaseIP: "172.16.0.2",
			Error:     nil,
		},
		{
			ReleaseIP: "172.16.0.3",
			Error:     nil,
		},
		{
			ReleaseIP: "192.168.0.6",
			Error:     network.ErrIPIsNotFound,
		},
	}

	GetIPParamAfterRelease = []struct {
		GetIP string
		Error error
	}{
		{
			GetIP: "172.16.0.2",
			Error: nil,
		},
		{
			GetIP: "172.16.0.3",
			Error: nil,
		},
		{
			GetIP: "172.16.0.6",
			Error: nil,
		},
		{
			GetIP: "172.16.0.7",
			Error: nil,
		},
	}
)
