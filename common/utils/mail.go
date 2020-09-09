package utils

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func SendMail(To, Cc, Bcc, att []string, host string, port int, user, password, subject, body, mailtype string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", user, "通知消息") // 发件人
	ToAddresses := make([]string, len(To))
	for i, recipient := range To {
		ToAddresses[i] = m.FormatAddress(recipient, "")
	}
	CcAddresses := make([]string, len(Cc))
	for i, recipient := range Cc {
		CcAddresses[i] = m.FormatAddress(recipient, "")
	}
	BccAddresses := make([]string, len(Bcc))
	for i, recipient := range Bcc {
		BccAddresses[i] = m.FormatAddress(recipient, "")
	}
	m.SetHeader("To", ToAddresses...) // 收件人
	m.SetHeader("Cc", CcAddresses...) //抄送
	m.SetHeader("Bcc", BccAddresses...)
	m.SetHeader("Subject", subject) // 主题
	m.SetBody(mailtype, body)       // 可以放html..
	for _, i := range att {
		m.Attach(i)
	}
	d := gomail.NewDialer(host, port, user, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	return err
}
