package EMail

import (
	"crypto/tls"
	"fmt"
	gomail "gopkg.in/mail.v2"
)

var emailFrom, emailTo, password string
var message *gomail.Message
var dialer *gomail.Dialer

func Init(emailFrom_ string, emailTo_ string, password_ string){
	emailFrom = emailFrom_
	emailTo = emailTo_
	password = password_

	message = gomail.NewMessage()

	// Set E-Mail sender
	message.SetHeader("From", emailFrom)

	// Set E-Mail receivers
	message.SetHeader("To", emailTo)

	// Settings for SMTP server
	dialer = gomail.NewDialer("smtp.gmail.com", 587, emailFrom, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false}
}

func Send(messageText string, messageSubject string){
	// Set E-Mail subject
	message.SetHeader("Subject", "Instagram bot. " + messageSubject)

	// Set E-Mail body. You can set plain text or html with text/html
	message.SetBody("text/plain", messageText)

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println(err)
	}
	return
}