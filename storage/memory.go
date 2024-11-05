package storage

import "net"

type DB struct {
	IPdb map[string]net.IPNet
}

func NewDB() DB {
	return DB{
		IPdb: make(map[string]net.IPNet),
	}
}
