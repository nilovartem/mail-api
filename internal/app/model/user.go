package model

import (
	"time"

	"github.com/google/uuid"
)

// User ...
type User struct {
	Mail string
	Link string
}

func (u *User) NewLink(ttl time.Duration) { //start & return timer & return User
	u.Link = uuid.New().String()
	//start timer
	go func() {
		ticker := time.NewTicker(ttl)
		for {
			<-ticker.C
			//link expired, "remove" it from list
			u.Link = ""
			/*if idx := slices.IndexFunc(s.Users, func(user *User) bool { return user.link == link }); idx != -1 {
				s.Users[idx].link = ""
			}*/
		}
	}()
}
func Validate() {}
func Zip()      {}

/*
func (s *Server) ZipHandler(mail string) ([]byte, error) {

	filename := filepath.Join(s.config.Mailbox, mail)
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	f, err := writer.Create(filename + "zip")
	if err != nil {
		return nil, err
	}
	_, err = f.Write([]byte(data))
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	//io.Copy(w, buf)
	return buf.Bytes(), nil
}
*/
