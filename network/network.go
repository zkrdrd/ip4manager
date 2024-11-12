package network

import (
	"encoding/binary"
	"errors"
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

type NetworkControl struct {
	network       net.IPNet
	UsedIPStorage map[[4]byte]net.IPNet
	FreeIPStorage map[[4]byte]net.IPNet
}

// Передается строка формата "192.168.0.0/16"
// Создается 2 карты
// 1. UsedIPStorage - для хранения используемых IP
// 2. FreeIPStorage - для хранения освобожденных IP
func NewNetwork(network string) (NetworkControl, error) {
	_, ipv4Net, err := net.ParseCIDR(network)
	return NetworkControl{
		network:       *ipv4Net,
		UsedIPStorage: make(map[[4]byte]net.IPNet),
		FreeIPStorage: make(map[[4]byte]net.IPNet),
	}, err
}

// Метод аренды IP адрессов
// Возвращает IP в строковом формате
func (netw NetworkControl) SetUsedIP() (string, error) {

	if len(netw.UsedIPStorage) == 0 {
		// gateway Example: 192.168.0.1
		netw.UsedIPStorage[[4]byte(nextIP(netw.network.IP, 1).To4())] = netw.network
	}

	if len(netw.FreeIPStorage) != 0 {
		return freeIPStorage(netw), nil
	}

	return usedIPStorage(netw)

}

// Метод осбождения IP адрессов из под аренды
func (netw NetworkControl) ReleaseIP(ip string) error {

	if len(netw.UsedIPStorage) == 0 {
		return ErrStorageIsEmpty
	}

	for storageIP := range netw.UsedIPStorage {
		ip := [4]byte(net.ParseIP(ip).To4())
		if !reflect.DeepEqual(ip, storageIP) {
			delete(netw.UsedIPStorage, ip)
			netw.FreeIPStorage[ip] = netw.network
			return nil
		}
	}

	return ErrIPIsNotFound
}

// Функция расчета следущего IP адресса
func nextIP(ip net.IP, inc uint) net.IP {
	ip = ip.To4()
	octets := uint(ip[0])<<24 + uint(ip[1])<<16 + uint(ip[2])<<8 + uint(ip[3])
	octets += inc
	octet4 := byte(octets & 0xFF)
	octet3 := byte((octets >> 8) & 0xFF)
	octet2 := byte((octets >> 16) & 0xFF)
	octet1 := byte((octets >> 24) & 0xFF)
	return net.IPv4(octet1, octet2, octet3, octet4)
}

func freeIPStorage(netw NetworkControl) string {
	maxIP := make(net.IP, len(netw.network.IP.To4()))
	binary.BigEndian.PutUint32(maxIP, binary.BigEndian.Uint32(netw.network.IP.To4())|^binary.BigEndian.Uint32(net.IP(netw.network.Mask).To4()))
	minKey := [4]byte(maxIP)

	for key := range netw.FreeIPStorage {
		if key[2] <= minKey[2] && key[3] < minKey[3] {
			minKey = key
		}
	}

	delete(netw.FreeIPStorage, minKey)
	netw.UsedIPStorage[minKey] = netw.network
	return net.IPv4(minKey[0], minKey[1], minKey[2], minKey[3]).String()
}

func usedIPStorage(netw NetworkControl) (string, error) {

	maxKey := [4]byte(netw.network.IP.To4())

	for key := range netw.UsedIPStorage {
		if key[2] >= maxKey[2] && key[3] > maxKey[3] {
			maxKey = key
		}
	}

	maxIP := net.IPv4(maxKey[0], maxKey[1], maxKey[2], maxKey[3])
	nextIP := [4]byte(nextIP(maxIP, 1).To4())

	if _, ok := netw.UsedIPStorage[nextIP]; !ok {
		if netw.network.Contains(maxIP) {
			netw.UsedIPStorage[nextIP] = netw.network
			return maxIP.String(), nil
		} else {
			return "", ErrIPADressIsNotIncludedInNetwork
		}
	}
	return "", nil
}
