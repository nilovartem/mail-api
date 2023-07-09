package config

import (
	"crypto/sha256"
	"encoding/json"
	"os"
	"time"
)

// Config contains all info from JSON file...
type Config struct {
	Mailbox     string
	TTL         time.Duration
	LogLevel    string
	BindAddress string
	Users       map[string][32]byte
	PDF         string
}

// NewConfig creates config and adjusts values from JSON
func NewConfig(configPath string, pdfPath string) (*Config, error) {
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
	config.PDF = pdfPath
	return &config, nil
}
