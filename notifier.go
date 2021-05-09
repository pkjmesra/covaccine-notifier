package main

import (
	"net/smtp"
	"log"
)

func sendMail(id, pass, body string, criteria string) error {
	msg := "From: " + id + "\n" +
		"To: " + id + "\n" +
		"Subject: Vaccination slots for " + criteria +" are available\n\n" +
		"Vaccination slots for " + criteria + " are available at the following centers:\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", id, pass, "smtp.gmail.com"),
		id, []string{id}, []byte(msg))

	if err != nil {
		return err
	}
	log.Print(body)
	return nil
}
