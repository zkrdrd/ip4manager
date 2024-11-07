package network

import (
	"encoding/binary"
	"errors"
	"net"
	"reflect"
)

var (
	ErrIPADressIsNotIncludedInNetwork = errors.New("ip address is not included in network")
	ErrStorageIsEmpty                 = errors.New("storage is empty")
	ErrIPAddressIsNotFound            = errors.New("ip address is not found")
)

type DB struct {
	IPdb map[string]net.IPNet
}

func NewDB() DB {
	return DB{
		IPdb: make(map[string]net.IPNet),
	}
}

func (db DB) SetUsedIP(netw net.IPNet) (string, error) {

	broadcast := make(net.IP, len(netw.IP.To4()))
	binary.BigEndian.PutUint32(broadcast, binary.BigEndian.Uint32(netw.IP.To4())|^binary.BigEndian.Uint32(net.IP(netw.Mask).To4()))

	if len(db.IPdb) == 0 {
		// network Example: 192.168.0.0
		db.IPdb[netw.IP.String()] = netw
		// gateway Example: 192.168.0.1
		db.IPdb[nextIP(netw.IP, 1).String()] = netw
		// broadcast Example: 192.168.255.255
		db.IPdb[broadcast.String()] = netw
	}

	for key := range db.IPdb {
		nextIP := nextIP(netw.IP, 1)
		storageIP := net.ParseIP(key)
		if !reflect.DeepEqual(nextIP, storageIP) {
			if netw.Contains(nextIP) {
				db.IPdb[nextIP.String()] = netw
				return nextIP.String(), nil
			} else {
				return "", ErrIPADressIsNotIncludedInNetwork
			}
		}
	}
	return "", nil
}

func (db DB) ReleaseIPAddress(ip string) error {

	if len(db.IPdb) == 0 {
		return ErrStorageIsEmpty
	}

	for key := range db.IPdb {
		storageIP := net.ParseIP(key)

		if !reflect.DeepEqual(ip, storageIP) {
			delete(db.IPdb, ip)
			return nil
		}
	}
	return ErrIPAddressIsNotFound
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
