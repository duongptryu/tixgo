package mail

import "github.com/duongptryu/gox/notification/mail"

type EventSendMail struct {
	ToMail   []mail.EmailAddress `json:"to_mail"`
	CC       []mail.EmailAddress `json:"cc"`
	BCC      []mail.EmailAddress `json:"bcc"`
	Subject  string              `json:"subject"`
	TextBody string              `json:"text_body"`
	HTMLBody string              `json:"html_body"`
	Priority mail.Priority       `json:"priority"`
}
