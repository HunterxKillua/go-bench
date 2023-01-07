package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/jordan-wright/email"
)

func SendEmail() {
	temp, err := ioutil.ReadFile("template/email.html")
	if err != nil {
		log.Fatal(err)
	}
	t := template.New("email")
	t, err = t.Parse(string(temp))
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]any{
		"name": "email模板",
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		log.Fatal(err)
	}
	e := email.NewEmail()
	e.From = fmt.Sprintf("王成杰 <%s>", address)
	e.To = []string{"1445290905@qq.com"}
	files := []string{"template/a.text", "template/avatar.jpg"}
	for index, filePath := range files {
		attachment, _ := e.AttachFile(filePath)
		attachment.Filename = fmt.Sprintf("附件%v", index)
	}
	e.Subject = "找回密码"
	e.HTML = buf.Bytes()
	tlsConfig := &tls.Config{
		ServerName: "smtp.gmail.com", //"smtp.office365.com",
	}
	err = e.SendWithStartTLS("smtp.gmail.com:587", Auth, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
}
