package main

import (
	"flag"

	"github.com/nilovartem/mail-api/internal/app/config"
	"github.com/nilovartem/mail-api/internal/app/mailapi"
	"github.com/sirupsen/logrus"
)

var (
	configPath = "configs/mailapi.json"
	pdfPath    = "static/readme.pdf"
)

func init() {
	flag.StringVar(&configPath, "config", configPath, "path to JSON file for server configuration")
	flag.StringVar(&pdfPath, "pdf", pdfPath, "path to PDF file")
}
func main() {
	flag.Parse()
	c, err := config.NewConfig(configPath, pdfPath)
	if err != nil {
		logrus.Fatal(err)
	}
	s := mailapi.NewServer(c)
	if err := s.Start(); err != nil {
		logrus.Fatal(err)
	}
}
