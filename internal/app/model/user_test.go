package model

import (
	"path/filepath"
	"testing"
)

var users = []*User{
	{Username: "beck.cierra@wb.ru"},
	{Username: "karlenko.anton@wb.ru"},
	{Username: "newton.alejandra@wb.ru"},
	{Username: "patrick.sharon@wb.ru"},
}

func TestZip(t *testing.T) {
	for _, u := range users {
		zip, err := u.Zip(filepath.Join("/Users/artem/go/src/mail-api/test/data/mailbox", u.Username), u.Username, "/Users/artem/go/src/mail-api/static/readme.pdf")
		if err != nil {
			t.Fatal(err)
		}
		if len(zip) == 0 {
			t.Fatal()
		}
	}
}
