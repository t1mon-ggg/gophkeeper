package models

import (
	"net"
	"time"
)

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	PublicKey string `json:"publickey"`
}

// Action -
type Action struct {
	Action   string    //
	Checksum string    //
	Sign     string    //
	ip       *net.IP   //
	date     time.Time //
}
