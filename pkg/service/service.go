package service

import (
	"net"
	"time"
)

func IsDatabaseUp() bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:3306", 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func DatabaseStart() error {
	return nil
}
