package main

type Forwarding struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

type Config struct {
	Version     uint         `json:"version"`
	Tunnel      string       `json:"tunnel"`
	Key         string       `json:"key"`
	Forwardings []Forwarding `json:"forwardings"`
}
