package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"reflect"
)

/*
Конфигурирование модуля:  адрес сети и маска.
Система сама определяем сколько в ней есть адресов и их диапазон с учетом адреса сети, первый адрес - адрес шлюза, и широковещательный адрес, последний адрес сети.
*/
var (
	ErrIPADressIsNotIncludedInNetwork = errors.New("ip address is not included in network")
	ErrStorageIsEmpty                 = errors.New("storage is empty")
	ErrIPIsNotFound                   = errors.New("ip address is not found")
	ErrNoFreeIPAddress                = errors.New("no free ip address")
	ErrNetMaskIsNotCorrect            = errors.New("net mask is not correct")
	octets                            = []byte{0, 128, 192, 224, 240, 248, 242, 254, 255}
)

type Network struct {
	network net.IPNet
	IPdb    map[string]net.IPNet
}

func NewNetwork(network string) (Network, error) {
	_, ipv4Net, err := net.ParseCIDR(network)
	return Network{
		network: *ipv4Net,
		IPdb:    make(map[string]net.IPNet),
	}, err
}

func SetUsedIP(ip string) {
	ipnet := net.ParseIP("192.168.0.5")
	fmt.Print(ipnet)
}

func (netw Network) SetUsedIP() (string, error) {

	broadcast := make(net.IP, len(netw.network.IP.To4()))
	binary.BigEndian.PutUint32(broadcast, binary.BigEndian.Uint32(netw.network.IP.To4())|^binary.BigEndian.Uint32(net.IP(netw.network.Mask).To4()))

	if len(netw.IPdb) == 0 {
		// network Example: 192.168.0.0
		netw.IPdb[netw.network.IP.String()] = netw.network
		// gateway Example: 192.168.0.1
		netw.IPdb[nextIP(netw.network.IP, 1).String()] = netw.network
		// broadcast Example: 192.168.255.255
		netw.IPdb[broadcast.String()] = netw.network
	}

	for k, v := range netw.IPdb {
		fmt.Print(k, v)
	}

	for key := range netw.IPdb {
		nextIP := nextIP(netw.network.IP, 1)
		storageIP := net.ParseIP(key)
		if !reflect.DeepEqual(nextIP, storageIP) {
			if netw.network.Contains(nextIP) {
				netw.IPdb[nextIP.String()] = netw.network
				return nextIP.String(), nil
			} else {
				return "", ErrIPADressIsNotIncludedInNetwork
			}
		}
	}
	return "", nil
}

func (netw Network) ReleaseIP(ip string) error {

	if len(netw.IPdb) == 0 {
		return ErrStorageIsEmpty
	}

	for key := range netw.IPdb {
		storageIP := net.ParseIP(key)

		if !reflect.DeepEqual(ip, storageIP) {
			delete(netw.IPdb, ip)
			return nil
		}
	}
	return ErrIPIsNotFound
}

func nextIP(ip net.IP, inc uint) net.IP {
	ip = ip.To4()
	octets := uint(ip[0])<<24 + uint(ip[1])<<16 + uint(ip[2])<<8 + uint(ip[3])
	octets += inc
	octets4 := byte(octets & 0xFF)
	octets3 := byte((octets >> 8) & 0xFF)
	octets2 := byte((octets >> 16) & 0xFF)
	octets1 := byte((octets >> 24) & 0xFF)
	return net.IPv4(octets1, octets2, octets3, octets4)
}

// func NewNetwork(firestIPAddressOctet, secondIPAddressOctet, thierdIPAddressOctet, fourthIPAddressOctet byte,
// 	firestMaskOctet, secondMaskOctet, thierdMaskOctet, fourthMaskOctet byte) (net.IPNet, error) {

// 	if err := checkMask(firestMaskOctet, secondMaskOctet, thierdMaskOctet, fourthMaskOctet); err != nil {
// 		return net.IPNet{
// 			IP:   nil,
// 			Mask: nil,
// 		}, err
// 	}

// 	ip := net.IPv4(firestIPAddressOctet, secondIPAddressOctet, thierdIPAddressOctet, fourthIPAddressOctet)
// 	mask := net.IPv4Mask(firestMaskOctet, secondMaskOctet, thierdMaskOctet, fourthMaskOctet)

// 	if err := checkIPAddress(ip, mask); err != nil {
// 		return net.IPNet{
// 			IP:   nil,
// 			Mask: nil,
// 		}, err
// 	}

// 	return net.IPNet{
// 		IP:   ip,
// 		Mask: mask,
// 	}, nil
// }

// func checkMask(firestMaskOctet, secondMaskOctet, thierdMaskOctet, fourthMaskOctet byte) error {
// 	for _, val := range octets {
// 		if firestMaskOctet == val && secondMaskOctet == 0 && thierdMaskOctet == 0 && fourthMaskOctet == 0 {
// 			return nil
// 		}
// 		if firestMaskOctet == 255 && secondMaskOctet == val && thierdMaskOctet == 0 && fourthMaskOctet == 0 {
// 			return nil
// 		}
// 		if firestMaskOctet == 255 && secondMaskOctet == 255 && thierdMaskOctet == val && fourthMaskOctet == 0 {
// 			return nil
// 		}
// 		if firestMaskOctet == 255 && secondMaskOctet == 255 && thierdMaskOctet == 255 && fourthMaskOctet == val {
// 			return nil
// 		}
// 	}
// 	return ErrNetMaskIsNotCorrect
// }

// func checkIPAddress(ip net.IP, mask net.IPMask) error {
// 	broadcast := make(net.IP, len(ip.To4()))
// 	binary.BigEndian.PutUint32(broadcast, binary.BigEndian.Uint32(ip.To4())|^binary.BigEndian.Uint32(net.IP(mask).To4()))

// 	if bytes.Compare(ip.To4(), broadcast.To4()) == 0 {
// 		return ErrNoFreeIPAddress
// 	}

// 	return nil
// }
