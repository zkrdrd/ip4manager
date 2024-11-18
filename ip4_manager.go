package ip4manager

import (
	"errors"
	"net"
	"reflect"
	"sync"
)

var (
	ErrStorageIsEmpty                 = errors.New("storage is empty")
	ErrNoFreeIPAddress                = errors.New("no free ip address")
	ErrIPAddressIsUsed                = errors.New("ip address is used")
	ErrIPIsNotFound                   = errors.New("ip address is not found")
	ErrNetMaskIsNotCorrect            = errors.New("net mask is not correct")
	ErrIPIsANetworkAddress            = errors.New("ip address is a network address")
	ErrIPADressIsNotIncludedInNetwork = errors.New("ip address is not included in network")
)

type IP4Manager struct {
	ip4manager    net.IPNet
	mx            *sync.RWMutex
	UsedIPStorage map[[4]byte]struct{}
	FreeIPStorage map[[4]byte]struct{}
}

// Передается строка формата "192.168.0.0/16"
func NewNetwork(ip4manager string) (IP4Manager, error) {
	_, ipv4Net, err := net.ParseCIDR(ip4manager)
	return IP4Manager{
		ip4manager:    *ipv4Net,
		mx:            &sync.RWMutex{},
		UsedIPStorage: make(map[[4]byte]struct{}),
		FreeIPStorage: make(map[[4]byte]struct{}),
	}, err
}

// Метод аренды IP адресса
func (ip4mng IP4Manager) GetFreeIP() (string, error) {

	ip4Byte := [4]byte(ip4mng.ip4manager.IP.To4())
	mask4Byte := [4]byte(net.IP(ip4mng.ip4manager.Mask).To4())
	var broadcast [4]byte

	for i := range len(ip4Byte) {
		broadcast[i] = ip4Byte[i] | ^mask4Byte[i]

	}

	if len(ip4mng.UsedIPStorage) == 0 {
		// gateway Example: 192.168.0.1
		ip4mng.mx.Lock()
		ip4mng.UsedIPStorage[[4]byte(nextIP(ip4mng.ip4manager.IP, 1).To4())] = struct{}{}
		ip4mng.mx.Unlock()
	}

	if len(ip4mng.FreeIPStorage) != 0 {
		return freeIPStorage(ip4mng, broadcast), nil
	}

	return getIPStorage(ip4mng, broadcast)

}

// Метод указания занятых IP
func (ip4mng IP4Manager) SetUsedIP(ip string) error {

	ipv4 := net.ParseIP(ip).To4()
	ip4Byte := [4]byte(ipv4)
	mask4Byte := [4]byte(net.IP(ip4mng.ip4manager.Mask).To4())
	var broadcast [4]byte

	for i := range len(ip4Byte) {
		broadcast[i] = ip4Byte[i] | ^mask4Byte[i]

	}

	if ipv4[2] == broadcast[2] && ipv4[3] == broadcast[3] || ipv4[2] == ip4mng.ip4manager.IP[2] && ipv4[3] == ip4mng.ip4manager.IP[3] {
		return ErrIPIsANetworkAddress
	}

	if !ip4mng.ip4manager.Contains(ipv4) {
		return ErrIPADressIsNotIncludedInNetwork
	}

	ip4mng.mx.RLock()
	if _, ok := ip4mng.UsedIPStorage[ip4Byte]; ok {
		ip4mng.mx.RUnlock()
		return ErrIPAddressIsUsed
	}
	ip4mng.mx.RUnlock()

	ip4mng.mx.Lock()
	ip4mng.UsedIPStorage[ip4Byte] = struct{}{}
	ip4mng.mx.Unlock()

	return nil
}

// Метод осбождения IP адрессов из под аренды
func (ip4mng IP4Manager) ReleaseIP(ip string) error {

	if len(ip4mng.UsedIPStorage) == 0 {
		return ErrStorageIsEmpty
	}

	ip4mng.mx.RLock()
	for storageIP := range ip4mng.UsedIPStorage {
		ip := [4]byte(net.ParseIP(ip).To4())
		if reflect.DeepEqual(ip, storageIP) {
			ip4mng.mx.RUnlock()
			ip4mng.mx.Lock()
			delete(ip4mng.UsedIPStorage, ip)
			ip4mng.FreeIPStorage[ip] = struct{}{}
			ip4mng.mx.Unlock()
			return nil
		}
	}
	ip4mng.mx.RUnlock()
	return ErrIPIsNotFound
}

// Используется storage освобожденных ip
func freeIPStorage(ip4mng IP4Manager, broadcast [4]byte) string {

	minKey := [4]byte(broadcast)
	ip4mng.mx.RLock()
	for key := range ip4mng.FreeIPStorage {
		if key[2] <= minKey[2] && key[3] < minKey[3] {
			minKey = key
		}
	}
	ip4mng.mx.RUnlock()

	ip4mng.mx.Lock()
	delete(ip4mng.FreeIPStorage, minKey)
	ip4mng.UsedIPStorage[minKey] = struct{}{}
	ip4mng.mx.Unlock()
	return net.IPv4(minKey[0], minKey[1], minKey[2], minKey[3]).String()
}

// используется storage занятых ip
func getIPStorage(ip4mng IP4Manager, broadcast [4]byte) (string, error) {

	maxKey := [4]byte(ip4mng.ip4manager.IP.To4())

	ip4mng.mx.RLock()
	for key := range ip4mng.UsedIPStorage {
		if key[2] >= maxKey[2] && key[3] > maxKey[3] {
			maxKey = key
		}
	}
	ip4mng.mx.RUnlock()

	maxIP := net.IPv4(maxKey[0], maxKey[1], maxKey[2], maxKey[3])
	netxIPv4 := nextIP(maxIP, 1).To4()
	nextIPbyte := [4]byte(netxIPv4)

	if nextIPbyte[2] == broadcast[2] && nextIPbyte[3] == broadcast[3] {
		return "", ErrNoFreeIPAddress
	}

	ip4mng.mx.RLock()
	if _, ok := ip4mng.UsedIPStorage[nextIPbyte]; !ok {
		ip4mng.mx.RUnlock()
		if ip4mng.ip4manager.Contains(netxIPv4) {
			ip4mng.mx.Lock()
			ip4mng.UsedIPStorage[nextIPbyte] = struct{}{}
			ip4mng.mx.Unlock()
			return netxIPv4.String(), nil
		} else {
			ip4mng.mx.Unlock()
			return "", ErrIPADressIsNotIncludedInNetwork
		}
	}
	ip4mng.mx.RUnlock()
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
