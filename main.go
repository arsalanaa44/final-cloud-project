package main

import (
	"log"
	"net/smtp"
)

func main(){
// Set up authentication information.
auth := smtp.PlainAuth(
    "",
    "mehran.akmah@gmail.com", // replace with your email
    "oemcglbamvusrrht", // replace with your password
    "smtp.gmail.com",
)

// Connect to the server, authenticate, set the sender and recipient,
// and send the email all in one step.
to := []string{"mehran.mah44@gmail.com"} // replace with the recipient email
msg := []byte("To:mehran.mah44@gmail.com\r\n" +
    "Subject: Hello Gopher!\r\n" +
    "\r\n" +
    "This is the email body.\r\n")
err := smtp.SendMail(
    "smtp.gmail.com:587",
    auth,
    "mehran.akmah@gmail.com", // replace with your email
    to,
    msg,
)
if err != nil {
    log.Fatal(err)
}}