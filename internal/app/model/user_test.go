package model

import (
	"path/filepath"
	"testing"

	"github.com/nilovartem/mail-api/internal/app/config"
)

var users = []*User{
	{Username: "beck.cierra@wb.ru"},
	{Username: "karlenko.anton@wb.ru"},
	{Username: "newton.alejandra@wb.ru"},
	{Username: "patrick.sharon@wb.ru"},
}

// TestZip ...
func TestZip(t *testing.T) {
	config, err := config.NewConfig("/Users/artem/go/src/mail-api/configs/mailapi.json", "/Users/artem/go/src/mail-api/static/readme.pdf")
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range users {
		zip, err := u.Zip(filepath.Join(config.Mailbox, u.Username), u.Username, config.PDF)
		if err != nil {
			t.Fatal(err)
		}
		if len(zip) == 0 {
			t.Fatal()
		}
	}
}
