package iam

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"server/internal/consts"
	_ "server/internal/logic/sys"
	"server/internal/model"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/iamin"
	"server/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

func TestInvitationMailTranslations(t *testing.T) {
	path, err := filepath.Abs("../../../manifest/i18n")
	if err != nil {
		t.Fatal(err)
	}
	translator := gi18n.New(gi18n.Options{Path: path})
	for _, language := range []string{"zh-CN", "zh-TW", "en-US"} {
		t.Run(language, func(t *testing.T) {
			if _, err := gjson.LoadPath(filepath.Join(path, language, "iam.yaml"), gjson.Options{}); err != nil {
				t.Fatalf("invalid IAM language file: %v", err)
			}
			ctx := gi18n.WithLanguage(context.Background(), language)
			for _, field := range []struct {
				key  string
				args []any
			}{
				{key: "subject", args: []any{"NextPT"}},
				{key: "greeting"},
				{key: "intro", args: []any{"NextPT", "friend@example.com"}},
				{key: "action"},
				{key: "expiry", args: []any{"2026-10-15 12:00:00 CST"}},
				{key: "permanent"},
				{key: "ignore"},
			} {
				key := "iam.invite.mail_" + field.key
				content := translator.GetContent(ctx, key)
				if content == "" {
					t.Fatalf("missing translation: %s", key)
				}
				if body := fmt.Sprintf(content, field.args...); strings.Contains(body, "%!") || strings.Contains(body, "iam.invite.mail_") {
					t.Fatalf("invalid rendered mail field %s: %s", key, body)
				}
			}
		})
	}
}

func TestMarkInviteSent(t *testing.T) {
	for _, tt := range []struct {
		name        string
		invite      *entity.IamInvite
		updateError error
		wantUpdate  bool
		wantError   bool
	}{
		{name: "missing", wantError: true},
		{name: "wrong owner", invite: &entity.IamInvite{Id: 7, InviterId: 99}, wantError: true},
		{name: "wrong id", invite: &entity.IamInvite{Id: 8, InviterId: 42}, wantError: true},
		{name: "already sent", invite: &entity.IamInvite{Id: 7, InviterId: 42, Status: consts.IamInviteStatusSent}, wantError: true},
		{name: "expired", invite: &entity.IamInvite{Id: 7, InviterId: 42, ExpireAt: gtime.Now().AddDate(0, 0, -1)}, wantError: true},
		{name: "update failure", invite: &entity.IamInvite{Id: 7, InviterId: 42, Hash: "invite-code"}, updateError: errors.New("db failed"), wantUpdate: true, wantError: true},
		{name: "success", invite: &entity.IamInvite{Id: 7, InviterId: 42, Hash: "invite-code"}, wantUpdate: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			originalDomain := service.IamInviteDomain()
			originalMail := service.SysMailgun()
			defer service.RegisterIamInviteDomain(originalDomain)
			defer service.RegisterSysMailgun(originalMail)
			mailer := &fakeInviteMailer{}
			domain := &fakeInviteDomain{mailer: mailer, err: tt.updateError}
			service.RegisterIamInviteDomain(domain)
			service.RegisterSysMailgun(mailer)
			err := NewIamInviteUsecase().markInviteSent(context.Background(), &model.Actor{Id: 42}, iamin.InviteSendInp{Id: 7, Email: " friend@example.com "}, tt.invite)
			if (err != nil) != tt.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tt.wantError)
			}
			if mailer.called || domain.called != tt.wantUpdate {
				t.Fatalf("mail = %v, update = %v", mailer.called, domain.called)
			}
			if tt.wantUpdate && (domain.data.Status != consts.IamInviteStatusSent || domain.data.InviteeEmail != "friend@example.com") {
				t.Fatalf("unexpected update: %#v", domain.data)
			}
		})
	}
}

func TestInvitationMailFailureDoesNotUpdateInvite(t *testing.T) {
	originalDomain := service.IamInviteDomain()
	originalMail := service.SysMailgun()
	defer service.RegisterIamInviteDomain(originalDomain)
	defer service.RegisterSysMailgun(originalMail)
	for _, mailError := range []error{nil, errors.New("mail failed")} {
		mailer := &fakeInviteMailer{err: mailError}
		domain := &fakeInviteDomain{mailer: mailer}
		service.RegisterIamInviteDomain(domain)
		service.RegisterSysMailgun(mailer)
		NewIamInviteUsecase().sendInvitationMail(context.Background(), "friend@example.com", &entity.IamInvite{Id: 7, Hash: "invite-code", Status: consts.IamInviteStatusSent})
		if !mailer.called || domain.called {
			t.Fatal("mail must be sent without updating or rolling back invite")
		}
		if !strings.Contains(mailer.text, "/register?invite=invite-code") || !strings.Contains(mailer.html, "/register?invite=invite-code") || mailer.recipient != "friend@example.com" {
			t.Fatal("mail missing registration link or correct recipient")
		}
	}
}

type fakeInviteMailer struct {
	service.ISysMailgun
	err                   error
	called                bool
	text, html, recipient string
}

func (f *fakeInviteMailer) SendHtmlMail(ctx context.Context, subject, text, html, recipient string) error {
	f.called = true
	f.text, f.html, f.recipient = text, html, recipient
	return f.err
}

type fakeInviteDomain struct {
	service.IIamInviteDomain
	mailer *fakeInviteMailer
	err    error
	called bool
	data   do.IamInvite
}

func (f *fakeInviteDomain) UpdateInvite(ctx context.Context, id uint64, data do.IamInvite) error {
	if f.mailer.called {
		return errors.New("mail sent before invite update")
	}
	f.called, f.data = true, data
	return f.err
}
