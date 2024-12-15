package mail

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	gomail "gopkg.in/mail.v2"
)

func ForgotPasswordEmailHtml(email string, verifyCode string, expireIn int) string {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Verify Email</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				margin: 0;
				padding: 0;
			}
			.container {
				background-color: #ffffff;
				max-width: 600px;
				margin: 20px auto;
				padding: 20px;
				border-radius: 8px;
				box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
				text-align: center;
			}
			h1 {
				color: #333333;
				font-size: 24px;
				text-align: center;
			}
			p {
				color: #666666;
				font-size: 16px;
				line-height: 1.5;
				text-align: center;
			}
			a.button {
				display: inline-block;
				margin: 20px auto;
				padding: 10px 20px;
				color: #ffffff;
				background-color: #007bff;
				text-decoration: none;
				border-radius: 5px;
				font-size: 16px;
				text-align: center;
			}
			.link {
				word-wrap: break-word;
			}
			.footer {
				color: #999999;
				font-size: 12px;
				text-align: center;
				margin-top: 20px;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Reset Your Password</h1>
			<p>We received a request to reset your password. If you didn't make the request, you can ignore this email.</p>
			<p>If you did make the request, you can reset your password using the following code:</p>
			<p class="link"><strong>` + verifyCode + `</strong></p>
			<p>This code will expire in ` + fmt.Sprintf("%d", expireIn) + ` minutes.</p>
			<a href="#" class="button">Reset Password</a>
			<p>If you're having trouble clicking the "Reset Password" button, copy and paste the URL below into your web browser:</p>
			<p class="link">http://localhost:3000/reset-password</p>
		</div>
	</body>
	</html>
	`
	return html
}

func SendMail(to string, subject string, body string) error {
	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", "bukharney@email.com")
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	// Set email body
	message.SetBody("text/html", body)

	// Create a new SMTP client
	dialer := MailConfig()

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		return err
	} else {
		fmt.Println("Email sent successfully!")
		return nil
	}
}

func MustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return v
}

func MailConfig() *gomail.Dialer {
	dialer := gomail.NewDialer(
		MustGetenv("SMTP_HOST"),
		587,
		MustGetenv("SMTP_USERNAME"),
		MustGetenv("SMTP_PASSWORD"),
	)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return dialer
}

func SentForgotPasswordMail(to string, token string) error {
	err := SendMail(to, "Reset Your Password", ForgotPasswordEmailHtml(to, token, 15))
	if err != nil {
		return err
	}

	log.Printf("Forgot password email sent to %s", to)

	return nil
}
