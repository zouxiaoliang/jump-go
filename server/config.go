package main

type Config struct {
	Version   uint     `json:"version"`
	Tunnel    string   `json:"tunnel"`
	Key       string   `json:"key"`
	Blacklist []string `json:"blacklist"`
	Whitelist []string `json:"whitelist"`
}
