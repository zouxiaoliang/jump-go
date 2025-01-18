package main

type Forwardings struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

type Config struct {
	Tunnel      string        `json:"tunnel"`
	Forwardings []Forwardings `json:"forwardings"`
}
