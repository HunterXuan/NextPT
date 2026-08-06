package iam

import (
	"context"
	"testing"

	"server/internal/model"
	"server/internal/service"
)

func TestSessionIdIsTokenDelimiterSafe(t *testing.T) {
	sessionId, err := (&sIamSessionDomain{}).newSessionId()
	if err != nil {
		t.Fatalf("newSessionId() error = %v", err)
	}
	if len(sessionId) != 32 {
		t.Fatalf("newSessionId() length = %d, want 32", len(sessionId))
	}
	for _, char := range sessionId {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			t.Fatalf("newSessionId() contains non-hex character %q", char)
		}
	}
}

func TestSessionListMarksCurrentSession(t *testing.T) {
	original := service.IamSessionDomain()
	defer service.RegisterIamSessionDomain(original)

	domain := &fakeSessionListDomain{
		sessions: []model.IamSession{
			{Id: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", UserId: 42, CreatedAt: 1, LastSeenAt: 2},
			{Id: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", UserId: 42, CreatedAt: 2, LastSeenAt: 3},
		},
	}
	service.RegisterIamSessionDomain(domain)

	out, err := NewIamSessionUsecase().List(context.Background(), &model.Actor{Id: 42}, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(out.List) != 2 {
		t.Fatalf("List() length = %d, want 2", len(out.List))
	}
	if out.List[0].Current || !out.List[1].Current {
		t.Fatalf("List() current flags = [%v, %v], want [false, true]", out.List[0].Current, out.List[1].Current)
	}
}

func TestSessionDeleteByIdChecksOwner(t *testing.T) {
	tests := []struct {
		name       string
		session    *model.IamSession
		wantErr    bool
		wantRemove bool
	}{
		{
			name:       "own session",
			session:    &model.IamSession{Id: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", UserId: 42},
			wantRemove: true,
		},
		{
			name:    "other user session",
			session: &model.IamSession{Id: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", UserId: 7},
			wantErr: true,
		},
		{
			name:    "missing session",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := service.IamSessionDomain()
			defer service.RegisterIamSessionDomain(original)

			domain := &fakeSessionListDomain{session: tt.session}
			service.RegisterIamSessionDomain(domain)
			err := NewIamSessionUsecase().DeleteById(context.Background(), &model.Actor{Id: 42}, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
			if tt.wantErr && err == nil {
				t.Fatal("DeleteById() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("DeleteById() error = %v", err)
			}
			if domain.removeCalls != boolInt(tt.wantRemove) {
				t.Fatalf("Remove() calls = %d, want %d", domain.removeCalls, boolInt(tt.wantRemove))
			}
		})
	}
}

type fakeSessionListDomain struct {
	service.IIamSessionDomain
	sessions    []model.IamSession
	session     *model.IamSession
	removeCalls int
}

func (f *fakeSessionListDomain) ListByUser(ctx context.Context, userId uint64) ([]model.IamSession, error) {
	return f.sessions, nil
}

func (f *fakeSessionListDomain) Get(ctx context.Context, sessionId string) (*model.IamSession, error) {
	return f.session, nil
}

func (f *fakeSessionListDomain) Remove(ctx context.Context, sessionId string) error {
	f.removeCalls++
	return nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
