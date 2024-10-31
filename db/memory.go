package db

import (
	"net"
)

type DB struct {
	IPdb map[*net.IPNet]bool
}

func NewDB() DB {
	return DB{
		IPdb: make(map[*net.IPNet]bool),
	}
}
