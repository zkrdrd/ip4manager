package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

/*
Конфигурирование модуля:  адрес сети и маска.
Система сама определяем сколько в ней есть адресов и их диапазон с учетом адреса сети, первый адрес - адрес шлюза, и широковещательный адрес, последний адрес сети.
*/
var (
	ErrNetworkIsNotCorrect = errors.New("network and broadcast is equal")
	ErrNetMaskIsNotCorrect = errors.New("net mask is not correct")
	octets                 = []byte{0, 128, 192, 224, 240, 248, 242, 254, 255}
)

func NewNetwork(firestIPAddressOctet, secondIPAddressOctet, thierdIPAddressOctet, fourthIPAddressOctet byte,
	firestMaskOctet, secondMaskOctet, thierdMaskOctet, fourthMaskOctet byte) (net.IPNet, error) {

	if err := checkMask(firestMaskOctet, secondMaskOctet, thierdMaskOctet, fourthMaskOctet); err != nil {
		return net.IPNet{
			IP:   nil,
			Mask: nil,
		}, err
	}

	ip := net.IPv4(firestIPAddressOctet, secondIPAddressOctet, thierdIPAddressOctet, fourthIPAddressOctet)
	mask := net.IPv4Mask(firestMaskOctet, secondMaskOctet, thierdMaskOctet, fourthMaskOctet)

	if err := checkIPAddress(ip, mask); err != nil {
		return net.IPNet{
			IP:   nil,
			Mask: nil,
		}, err
	}

	return net.IPNet{
		IP:   ip,
		Mask: mask,
	}, nil
}

func checkMask(firestMaskOctet, secondMaskOctet, thierdMaskOctet, fourthMaskOctet byte) error {
	for _, val := range octets {
		if firestMaskOctet == val && secondMaskOctet == 0 && thierdMaskOctet == 0 && fourthMaskOctet == 0 {
			return nil
		}
		if firestMaskOctet == 255 && secondMaskOctet == val && thierdMaskOctet == 0 && fourthMaskOctet == 0 {
			return nil
		}
		if firestMaskOctet == 255 && secondMaskOctet == 255 && thierdMaskOctet == val && fourthMaskOctet == 0 {
			return nil
		}
		if firestMaskOctet == 255 && secondMaskOctet == 255 && thierdMaskOctet == 255 && fourthMaskOctet == val {
			return nil
		}
	}
	return ErrNetMaskIsNotCorrect
}

func checkIPAddress(ip net.IP, mask net.IPMask) error {
	broadcast := make(net.IP, len(ip.To4()))
	binary.BigEndian.PutUint32(broadcast, binary.BigEndian.Uint32(ip.To4())|^binary.BigEndian.Uint32(net.IP(mask).To4()))

	if bytes.Compare(ip.To4(), broadcast.To4()) == 0 {
		return ErrNetworkIsNotCorrect
	}

	return nil
}
