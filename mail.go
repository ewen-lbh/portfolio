package main

import (
	"fmt"
	"os"

	"gopkg.in/mail.v2"
)

func SendMailToSelf(from, subject, body string) error {
	// fmt.Println("setting auth")
	// auth := smtp.PlainAuth("", "hey@ewen.works", os.Getenv("MAIL_PASSWORD"), "mail.ewen.works:465")
	// fmt.Printf("sending mail with %#v\n", auth)
	// return smtp.SendMail("mail.ewen.works:465", auth, "contact-form@ewen.works", []string{"contact@ewen.works"}, []byte(body))
	fmt.Printf("[%s] %s -> %s\n", subject, body, from)
	envelope := mail.NewMessage()
	envelope.SetHeader("From", "contact-form@ewen.works")
	envelope.SetHeader("To", "contact@ewen.works")
	envelope.SetHeader("Reply-To", from)
	envelope.SetHeader("Subject", subject)
	envelope.SetBody("text/plain", body)
	dialer := mail.NewDialer("mail.ewen.works", 465, "hey@ewen.works", os.Getenv("MAIL_PASSWORD"))
	dialer.StartTLSPolicy = mail.MandatoryStartTLS

	return dialer.DialAndSend(envelope)
}
