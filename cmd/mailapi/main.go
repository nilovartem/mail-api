package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nilovartem/mail-api/internal/app/mailapi"
	"github.com/sirupsen/logrus"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "configs/mailapi.json", "path to JSON file for server configuration")
}
func main() {
	flag.Parse()
	c, err := mailapi.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", c)
	s := mailapi.NewServer(c)
	if err := s.Start(); err != nil {
		logrus.Fatal(err)
	}
}
