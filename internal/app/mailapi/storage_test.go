package mailapi

import (
	"fmt"
	"testing"
	"time"
)

var (
	usernames = []string{
		"beck.cierra@wb.ru",
		"karlenko.anton@wb.ru",
		"newton.alejandra@wb.ru",
		"patrick.sharon@wb.ru",
	}
)

// TestNewLink ...
func TestNewLink(t *testing.T) {
	storage := NewStorage()
	for _, uname := range usernames {
		link := storage.Add(uname)
		go func() {
			ticker := time.NewTicker(time.Second * 1)
			for {
				<-ticker.C
				storage.Remove(link)
			}
		}()
	}
	timer := time.NewTimer(time.Second * 2)
	<-timer.C
	for _, uname := range usernames {
		link, found := storage.GetLink(uname)
		if found {
			t.Fatal(fmt.Errorf("expected expired link, got link %v", link))
		}
	}
}
