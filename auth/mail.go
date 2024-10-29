package main

import (
	"fmt"
	"html/template"
	"strconv"

	"github.com/wneessen/go-mail"
)

type MailInfo struct {
	ReceiverMail string
	BodyTemplate *template.Template
	BodyData     interface{}
	Subject      string
}

type ActivationInfo struct {
	ActivationLink template.URL
}

func sendActivationMain(email string, username string, code string) error {
	t := template.Must(template.ParseFiles("templates/activation.html"))

	activationInfo := ActivationInfo{
		ActivationLink: template.URL(fmt.Sprintf("https://%s:%s/auth/activate?username=%s&code=%s", hostName, hostPort, username, code)),
	}

	mailInfo := MailInfo{
		ReceiverMail: email,
		BodyTemplate: t,
		BodyData:     activationInfo,
		Subject:      EMail.ActivationSubject,
	}

	if err := sendMail(mailInfo); err != nil {
		return err
	}

	return nil
}

func sendMail(mainInfo MailInfo) error {
	m := mail.NewMsg()

	m.From(fromAddress)

	m.To(mainInfo.ReceiverMail)

	m.Subject(mainInfo.Subject)

	m.SetBodyHTMLTemplate(mainInfo.BodyTemplate, mainInfo.BodyData)

	port, err := strconv.Atoi(smtpPort)

	if err != nil {
		panic(err)
	}

	var client *mail.Client

	if devEnv {
		client, err = mail.NewClient(smtpHost, mail.WithPort(port), mail.WithTLSPortPolicy(mail.NoTLS),
			mail.WithSMTPAuth(mail.SMTPAuthNoAuth), mail.WithoutNoop(), mail.WithHELO("helo"))
		fmt.Println(client.ServerAddr())
	} else {
		client, err = mail.NewClient(smtpHost, mail.WithPort(port), mail.WithTLSPortPolicy(mail.TLSMandatory),
			mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(smtpUser), mail.WithPassword(smtpPassword),
		)
	}

	if err != nil {
		return err
	}

	if err := client.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
