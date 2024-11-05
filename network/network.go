package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"ipaddresspackage/storage"
	"net"
	"reflect"
)

/*
Конфигурирование модуля:  адрес сети и маска.
Система сама определяем сколько в ней есть адресов и их диапазон с учетом адреса сети, первый адрес - адрес шлюза, и широковещательный адрес, последний адрес сети.
*/

var (
	ErrNetworkIsNotCorrect            = errors.New("network is not correct")
	ErrNetMaskIsNotCorrect            = errors.New("net mask is not correct")
	ErrIPADressIsNotIncludedInNetwork = errors.New("ip address is not included in network")
)

func NewNetwork(firestIPAddressOctet, secondIPAddressOctet, thierdIPAddressOctet, fourthIPAddressOctet byte,
	firestMaskOctet, secondMaskOctet, thierdMaskOctet, fourthMaskOctet byte) net.IPNet {

	ip := net.IPv4(firestIPAddressOctet, secondIPAddressOctet, thierdIPAddressOctet, fourthIPAddressOctet)
	mask := net.IPv4Mask(firestMaskOctet, secondMaskOctet, thierdMaskOctet, fourthMaskOctet)

	return net.IPNet{
		IP:   ip,
		Mask: mask,
	}
}

func NewNetwrokMapping(netw net.IPNet) {
	db := storage.NewDB()
	// fmt.Println(netw.Contains(net.ParseIP(`172.16.255.255`)))
	// // если карта пустая то идем от идем от 0 если там что то есть проверяем карту

	// broadcust := make(net.IP, len(netw.IP.To4()))
	// binary.BigEndian.PutUint32(broadcust, binary.BigEndian.Uint32(netw.IP.To4())|^binary.BigEndian.Uint32(net.IP(netw.Mask).To4()))

	// fmt.Println(broadcust)
	// fmt.Println(nextIP(broadcust, 1))
	selectIPAddress(netw, db)
	fmt.Println(db.IPdb)
}

func selectIPAddress(netw net.IPNet, db storage.DB) error {
	netwIP := netw.IP

	broadcust := make(net.IP, len(netw.IP.To4()))
	binary.BigEndian.PutUint32(broadcust, binary.BigEndian.Uint32(netw.IP.To4())|^binary.BigEndian.Uint32(net.IP(netw.Mask).To4()))

	if len(db.IPdb) == 0 {
		db.IPdb[netwIP.String()] = netw
		db.IPdb[broadcust.String()] = netw
	}

	for key := range db.IPdb {
		nextIP := nextIP(netwIP, 1)
		storageIP := net.ParseIP(key)
		if !reflect.DeepEqual(nextIP, storageIP) {
			// TODO:
			// add locker
			// return ip
			if netw.Contains(nextIP) {
				db.IPdb[nextIP.String()] = netw
			} else {
				return ErrIPADressIsNotIncludedInNetwork
			}
		}
		continue
	}
	return nil
}

func nextIP(ip net.IP, inc uint) net.IP {
	octets := uint(ip[0])<<24 + uint(ip[1])<<16 + uint(ip[2])<<8 + uint(ip[3])
	octets += inc
	octets4 := byte(octets & 0xFF)
	octets3 := byte((octets >> 8) & 0xFF)
	octets2 := byte((octets >> 16) & 0xFF)
	octets1 := byte((octets >> 24) & 0xFF)
	return net.IPv4(octets1, octets2, octets3, octets4)
}
