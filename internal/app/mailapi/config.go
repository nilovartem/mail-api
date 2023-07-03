package mailapi

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// phonyConfig ...
type phonyConfig struct {
	Mailbox string `json:"mailbox"`
	TTL     string `json:"ttl"`
}

// Config ...
type Config struct {
	Mailbox string
	TTL     time.Duration
}

// NewConfig ...
func NewConfig(configPath string) (*Config, error) {
	c := Config{}
	p := phonyConfig{}
	contents, err := os.ReadFile(configPath)
	fmt.Println(string(contents))
	if err != nil {
		logrus.Debugln("[FAIL] Can't read config file, err = ", err)
		return &c, err
	}
	if err = json.Unmarshal(contents, &p); err != nil {
		logrus.Debugln("[FAIL] Can't unmarshal config file, err = ", err)
		return &c, err
	}
	c.Mailbox = p.Mailbox
	if c.TTL, err = time.ParseDuration(p.TTL); err != nil {
		logrus.Debugln("[FAIL] Can't parse TTL duration, err = ", err)
		return &c, err
	}
	return &c, nil
}
