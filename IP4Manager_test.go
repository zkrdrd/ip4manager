package ip4manager_test

import (
	"fmt"
	"ip4manager"
	"testing"
)

func TestNetwork(t *testing.T) {
	net, _ := ip4manager.NewNetwork(networkString)

	for _, expect := range GetIP {
		ip, err := net.GetFreeIP()
		if expect.GetIP != ip {
			t.Error(fmt.Errorf(`result field %v != %v`, expect.GetIP, ip))
		}
		if expect.Error != err {
			t.Error(fmt.Errorf(`result field %v != %v`, expect.Error, err))
		}
	}

	for _, expect := range SetUsedIP {
		err := net.SetUsedIP(expect.SetIP)
		if expect.Error != err {
			t.Error(fmt.Errorf(`result field %v != %v`, expect.Error, err))
		}
	}

	for _, expect := range ReleaseIP {
		err := net.ReleaseIP(expect.ReleaseIP)
		if expect.Error != err {
			t.Error(fmt.Errorf(`result field %v != %v`, expect.Error, err))
		}
	}

	for _, expect := range GetIPAfterRelease {
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

	GetIP = []struct {
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

	SetUsedIP = []struct {
		SetIP string
		Error error
	}{
		{
			SetIP: "172.16.0.5",
			Error: nil,
		},
		{
			SetIP: "172.16.0.0",
			Error: ip4manager.ErrIPIsANetworkAddress,
		},
		{
			SetIP: "172.16.255.255",
			Error: ip4manager.ErrIPIsANetworkAddress,
		},
		{
			SetIP: "172.16.0.3",
			Error: ip4manager.ErrIPAddressIsUsed,
		},
		{
			SetIP: "192.168.0.6",
			Error: ip4manager.ErrIPADressIsNotIncludedInNetwork,
		},
	}

	ReleaseIP = []struct {
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
			Error:     ip4manager.ErrIPIsNotFound,
		},
	}

	GetIPAfterRelease = []struct {
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
