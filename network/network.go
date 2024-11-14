package network

import (
	"encoding/binary"
	"errors"
	"net"
	"reflect"
	"sync"
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
	ErrIPAddressIsUsed                = errors.New("ip address is used")
	octets                            = []byte{0, 128, 192, 224, 240, 248, 242, 254, 255}
)

type NetworkControl struct {
	network       net.IPNet
	mx            *sync.RWMutex
	UsedIPStorage map[[4]byte]struct{}
	FreeIPStorage map[[4]byte]struct{}
}

// Передается строка формата "192.168.0.0/16"
// Создается 2 карты
// 1. UsedIPStorage - для хранения используемых IP
// 2. FreeIPStorage - для хранения освобожденных IP
func NewNetwork(network string) (NetworkControl, error) {
	_, ipv4Net, err := net.ParseCIDR(network)
	return NetworkControl{
		network:       *ipv4Net,
		mx:            &sync.RWMutex{},
		UsedIPStorage: make(map[[4]byte]struct{}),
		FreeIPStorage: make(map[[4]byte]struct{}),
	}, err
}

// Метод аренды IP адрессов
// Возвращает IP в строковом формате
func (netw NetworkControl) GetFreeIP() (string, error) {

	broadcast := make(net.IP, len(netw.network.IP.To4()))
	binary.BigEndian.PutUint32(broadcast, binary.BigEndian.Uint32(netw.network.IP.To4())|^binary.BigEndian.Uint32(net.IP(netw.network.Mask).To4()))

	if len(netw.UsedIPStorage) == 0 {
		// gateway Example: 192.168.0.1
		netw.mx.Lock()
		netw.UsedIPStorage[[4]byte(nextIP(netw.network.IP, 1).To4())] = struct{}{}
		netw.mx.Unlock()
	}

	if len(netw.FreeIPStorage) != 0 {
		return freeIPStorage(netw, broadcast), nil
	}

	return getIPStorage(netw, broadcast)

}

func (netw NetworkControl) SetUsedIP(ip string) error {

	ipv4 := net.ParseIP(ip).To4()
	ipByte := [4]byte(ipv4)

	if !netw.network.Contains(ipv4) {
		return ErrIPADressIsNotIncludedInNetwork
	}
	netw.mx.RLock()
	if _, ok := netw.UsedIPStorage[ipByte]; ok {
		netw.mx.RUnlock()
		return ErrIPAddressIsUsed
	}
	netw.mx.RUnlock()

	netw.mx.Lock()
	netw.UsedIPStorage[ipByte] = struct{}{}
	netw.mx.Unlock()

	return nil
}

// Метод осбождения IP адрессов из под аренды
func (netw NetworkControl) ReleaseIP(ip string) error {

	if len(netw.UsedIPStorage) == 0 {
		return ErrStorageIsEmpty
	}

	netw.mx.RLock()
	for storageIP := range netw.UsedIPStorage {
		ip := [4]byte(net.ParseIP(ip).To4())
		if !reflect.DeepEqual(ip, storageIP) {
			netw.mx.RUnlock()
			netw.mx.Lock()
			delete(netw.UsedIPStorage, ip)
			netw.FreeIPStorage[ip] = struct{}{}
			netw.mx.Unlock()
			return nil
		}
	}
	return ErrIPIsNotFound
}

// Используется storage освобожденных ip
func freeIPStorage(netw NetworkControl, broadcast net.IP) string {

	minKey := [4]byte(broadcast)
	netw.mx.RLock()
	for key := range netw.FreeIPStorage {
		if key[2] <= minKey[2] && key[3] < minKey[3] {
			minKey = key
		}
	}
	netw.mx.RUnlock()

	netw.mx.Lock()
	delete(netw.FreeIPStorage, minKey)
	netw.UsedIPStorage[minKey] = struct{}{}
	netw.mx.Unlock()
	return net.IPv4(minKey[0], minKey[1], minKey[2], minKey[3]).String()
}

// используется storage занятых ip
func getIPStorage(netw NetworkControl, broadcast net.IP) (string, error) {

	maxKey := [4]byte(netw.network.IP.To4())

	netw.mx.RLock()
	for key := range netw.UsedIPStorage {
		if key[2] >= maxKey[2] && key[3] > maxKey[3] {
			maxKey = key
		}
	}
	netw.mx.RUnlock()

	maxIP := net.IPv4(maxKey[0], maxKey[1], maxKey[2], maxKey[3])
	nextIP := [4]byte(nextIP(maxIP, 1).To4())

	if nextIP[2] == broadcast[2] && nextIP[3] == broadcast[3] {
		return "", ErrNoFreeIPAddress
	}

	netw.mx.RLock()
	if _, ok := netw.UsedIPStorage[nextIP]; !ok {
		netw.mx.RUnlock()
		if netw.network.Contains(maxIP) {
			netw.mx.Lock()
			netw.UsedIPStorage[nextIP] = struct{}{}
			netw.mx.Unlock()
			return maxIP.String(), nil
		} else {
			return "", ErrIPADressIsNotIncludedInNetwork
		}
	}
	netw.mx.RUnlock()
	return "", nil
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
