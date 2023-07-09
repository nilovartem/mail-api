package mailapi

import (
	"crypto/sha256"
	"encoding/json"
	"os"
	"time"
)

// Config ...
type Config struct {
	Mailbox     string
	TTL         time.Duration
	LogLevel    string
	BindAddress string
	Users       map[string][32]byte
}

// NewConfig ...
func NewConfig(configPath string) (*Config, error) {
	config := Config{}
	var phony struct {
		Mailbox     string            `json:"mailbox"`
		TTL         string            `json:"ttl"`
		LogLevel    string            `json:"log_level"`
		BindAddress string            `json:"bind_address"`
		Users       map[string]string `json:"users"`
	}
	contents, err := os.ReadFile(configPath)
	if err != nil {
		return &config, err
	}
	if err = json.Unmarshal(contents, &phony); err != nil {
		return &config, err
	}
	if config.TTL, err = time.ParseDuration(phony.TTL); err != nil {
		return &config, err
	}
	config.Mailbox = phony.Mailbox
	config.LogLevel = phony.LogLevel
	config.BindAddress = phony.BindAddress
	config.Users = make(map[string][32]byte)
	for key, value := range phony.Users {
		config.Users[key] = sha256.Sum256([]byte(value))
	}
	//TODO: check map empty or not
	return &config, nil
}
