package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sort"
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
	UsedIPStorage map[string]net.IPNet
	FreeIPStorage map[string]net.IPNet
}

// Передается строка формата "192.168.0.0/16"
// Создается 2 карты
// 1. UsedIPStorage - для хранения используемых IP
// 2. FreeIPStorage - для хранения освобожденных IP
func NewNetwork(network string) (NetworkControl, error) {
	_, ipv4Net, err := net.ParseCIDR(network)
	return NetworkControl{
		network:       *ipv4Net,
		UsedIPStorage: make(map[string]net.IPNet),
		FreeIPStorage: make(map[string]net.IPNet),
	}, err
}

// Метод аренды IP адрессов
// Возвращает IP в строковом формате
func (netw NetworkControl) SetUsedIP() (string, error) {

	broadcast := make(net.IP, len(netw.network.IP.To4()))
	binary.BigEndian.PutUint32(broadcast, binary.BigEndian.Uint32(netw.network.IP.To4())|^binary.BigEndian.Uint32(net.IP(netw.network.Mask).To4()))

	if len(netw.UsedIPStorage) == 0 {
		// broadcast Example: 192.168.255.255
		netw.UsedIPStorage[broadcast.String()] = netw.network
		// network Example: 192.168.0.0
		netw.UsedIPStorage[netw.network.IP.String()] = netw.network
		// gateway Example: 192.168.0.1
		netw.UsedIPStorage[nextIP(netw.network.IP, 1).String()] = netw.network
	}

	pl := make([]string, len(netw.UsedIPStorage))
	//sort.Sort(sort.Reverse(sort.StringSlice(pl)))
	sort.Strings(pl)

	if len(netw.FreeIPStorage) != 0 {
		for key := range netw.FreeIPStorage {
			delete(netw.FreeIPStorage, key)
			netw.UsedIPStorage[key] = netw.network
			fmt.Println(key)
			return key, nil
		}
	}

	for k := range netw.UsedIPStorage {
		fmt.Println("strg", k)
	}

	for key := range netw.UsedIPStorage {
		nextIP := nextIP(net.ParseIP(key), 1)
		if _, ok := netw.UsedIPStorage[nextIP.String()]; !ok {
			if netw.network.Contains(nextIP) {
				netw.UsedIPStorage[nextIP.String()] = netw.network
				return nextIP.String(), nil
			} else {
				return "", ErrIPADressIsNotIncludedInNetwork
			}
		}

	}

	return "", nil
}

// Метод осбождения IP адрессов из под аренды
func (netw NetworkControl) ReleaseIP(ip string) error {

	if len(netw.UsedIPStorage) == 0 {
		return ErrStorageIsEmpty
	}

	for key := range netw.UsedIPStorage {
		storageIP := net.ParseIP(key)

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
