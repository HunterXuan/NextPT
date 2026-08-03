package sys

import (
	"context"
	"strings"
	"time"

	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mailgun/mailgun-go/v4"
)

const mailgunSendTimeout = 10 * time.Second

type sSysMailgun struct{}

func NewSysMailgun() *sSysMailgun {
	return &sSysMailgun{}
}

func init() {
	service.RegisterSysMailgun(NewSysMailgun())
}

func (m *sSysMailgun) SendTextMail(ctx context.Context, subject string, body string, recipient string) error {
	return m.send(ctx, subject, body, "", recipient)
}

func (m *sSysMailgun) SendHtmlMail(ctx context.Context, subject string, textBody string, htmlBody string, recipient string) error {
	return m.send(ctx, subject, textBody, htmlBody, recipient)
}

func (m *sSysMailgun) send(ctx context.Context, subject string, textBody string, htmlBody string, recipient string) error {
	domain := strings.TrimSpace(g.Cfg().MustGet(ctx, "mailgun.domain").String())
	key := strings.TrimSpace(g.Cfg().MustGet(ctx, "mailgun.key").String())
	sender := strings.TrimSpace(g.Cfg().MustGet(ctx, "mailgun.sender").String())
	recipient = strings.TrimSpace(recipient)

	if domain == "" {
		return gerror.New("mailgun.domain is required")
	}
	if key == "" {
		return gerror.New("mailgun.key is required")
	}
	if sender == "" {
		return gerror.New("mailgun.sender is required")
	}
	if recipient == "" {
		return gerror.New("mail recipient is required")
	}

	mg := mailgun.NewMailgun(domain, key)
	message := mailgun.NewMessage(sender, subject, textBody, recipient)
	if htmlBody != "" {
		message.SetHTML(htmlBody)
	}

	ctx, cancel := context.WithTimeout(ctx, mailgunSendTimeout)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	return err
}
