package network

import (
	"errors"
	"fmt"
	"ipaddresspackage/db"
	"math/big"
	"net"
	"strconv"

	"github.com/apparentlymart/go-cidr/cidr"
)

/*
Конфигурирование модуля:  адрес сети и маска.
Система сама определяем сколько в ней есть адресов и их диапазон с учетом адреса сети, первый адрес - адрес шлюза, и широковещательный адрес, последний адрес сети.
*/

var (
	ErrNetworkIsNotCorrect = errors.New("network is not correct")
	ErrNetMaskIsNotCorrect = errors.New("net mask is not correct")
)

func verifiedNetworkData(NetwrokAddresses string, NetwrokMask string) (net.IP, net.IPMask, error) {
	address := net.ParseIP(NetwrokAddresses).To4()
	if address == nil {
		return nil, nil, ErrNetworkIsNotCorrect
	}

	if len(NetwrokMask) >= 1 && len(NetwrokMask) <= 2 {
		mask, err := strconv.Atoi(NetwrokMask)
		if err != nil {
			return nil, nil, err
		}

		if mask <= 0 && mask >= 32 {
			return nil, nil, ErrNetMaskIsNotCorrect
		}
	} else if len(NetwrokMask) >= 7 && len(NetwrokMask) <= 15 {
		mask := net.ParseIP(NetwrokMask).To4()
		if mask == nil {
			return nil, nil, ErrNetMaskIsNotCorrect
		}
		netMask := net.IPMask(mask)
		return address, netMask, nil
	}
	return nil, nil, nil
}

func NewNetwork(NetwrokAddresses string, NetwrokMask string) (net.IPNet, error) {

	address, netMask, err := verifiedNetworkData(NetwrokAddresses, NetwrokMask)
	if err != nil {
		return net.IPNet{
			IP:   nil,
			Mask: nil,
		}, err
	}

	return net.IPNet{
		IP:   address,
		Mask: netMask,
	}, nil
}

func NewNetwrokMapping(net net.IPNet) {
	db := db.NewDB()

	first, second := cidr.AddressRange(&net)

	startIP, mask := ipToInt(first)
	finishIP, mask := ipToInt(second)

	for i := new(big.Int).Set(startIP); i.Cmp(finishIP) < 0; i.Add(i, big.NewInt(1)) {
		net.IP = intToIP(i, mask)
		db.IPdb[&net] = true
	}

	for key, val := range db.IPdb {
		fmt.Println(key, val)
	}

}

func ipToInt(ip net.IP) (*big.Int, int) {
	val := &big.Int{}
	val.SetBytes([]byte(ip))
	if len(ip) == net.IPv4len {
		return val, 32
	} else if len(ip) == net.IPv6len {
		return val, 128
	} else {
		panic(fmt.Errorf("Unsupported address length %d", len(ip)))
	}
}

func intToIP(ipInt *big.Int, bits int) net.IP {
	ipBytes := ipInt.Bytes()
	ret := make([]byte, bits/8)
	for i := 1; i <= len(ipBytes); i++ {
		ret[len(ret)-i] = ipBytes[len(ipBytes)-i]
	}
	return net.IP(ret)
}
