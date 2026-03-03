package gmail

import (
	"backend/internal/config"
	"fmt"
	"net/smtp"
)

func SendPasswordResetEmail(email, resetLink string) error {
	htmlBody := fmt.Sprintf(`
<html>
<body>
    <h2>Восстановление пароля</h2>
    <p>Для восстановления пароля перейдите по ссылке:</p>
    <p><code>%s</code></p>
    <p><strong>Ссылка действительна в течение 10 минут.</strong></p>
    <hr>
    <p>Если вы не запрашивали восстановление пароля, проигнорируйте это письмо.</p>
</body>
</html>`, resetLink)

	return SendGmail(email, "newPass", htmlBody)
}

func SendVerificationCode(email, code string) error {
	htmlBody := fmt.Sprintf(`
<html>
<body>
    <h2>Код подтверждения</h2>
    <p>Ваш код подтверждения: <strong>%s</strong></p>
    <p>Введите этот код в приложении для завершения регистрации.</p>
    <p><strong>Код действителен в течение 5 минут.</strong></p>
</body>
</html>`, code)

	return SendGmail(email, "verifyCode", htmlBody)
}

func SendGmail(to, theme, body string) error {
	cfg := config.LoadConfig()

	pass := cfg.Google.GooglePass
	from := cfg.Google.From

	subject := ""
	switch theme {
	case "verifyCode":
		subject = "Код подтверждения"
	case "newPass":
		subject = "Ссылка для восстановления пароля"
	default:
		subject = theme
	}

	msg := fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: quoted-printable\r\n" +
		"\r\n" +
		body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		return err
	}

	return nil
}
