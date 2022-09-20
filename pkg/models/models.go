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
	Action   string    `json:"action"`   //
	Checksum string    `json:"checksum"` //
	IP       *net.IP   `json:"ip"`       //
	Date     time.Time `json:"time"`     //
}

type Content struct {
	Payload string `json:"payload"`
	Hash    string `json:"hashsum"`
}

type Version struct {
	Date time.Time `json:"time"`
	Hash string    `json:"hashsum"`
}

type PGP struct {
	Date      time.Time `json:"time,omitempty"`
	Publickey string    `json:"publickey"`
	Confirmed bool      `json:"confirmed,omitempty"`
}

type Message struct {
	Text    string `json:"text"`
	Content string `json:"content"`
}
