package mailapi

import (
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
}

// NewConfig ...
func NewConfig(configPath string) (*Config, error) {
	config := Config{}
	var phony struct {
		Mailbox     string `json:"mailbox"`
		TTL         string `json:"ttl"`
		LogLevel    string `json:"log_level"`
		BindAddress string `json:"bind_address"`
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
	return &config, nil
}
