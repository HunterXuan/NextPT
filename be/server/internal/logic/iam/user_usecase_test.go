package iam

import (
	"context"
	"testing"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/iamin"
	"server/internal/service"

	"github.com/goflyfox/gtoken/v2/gtoken"
	"golang.org/x/crypto/bcrypt"
)

func TestUpdateProfileUsesActorWithoutLoadingUser(t *testing.T) {
	origUserDomain := service.IamUserDomain()
	defer service.RegisterIamUserDomain(origUserDomain)

	fakeUserDomain := &fakeIamUserDomain{}
	service.RegisterIamUserDomain(fakeUserDomain)

	err := NewIamUserUsecase().UpdateProfile(context.Background(), &model.Actor{Id: 42}, iamin.UserProfileUpdateInp{
		Avatar:    "https://nextpt.local/avatar.png",
		Info:      "hello",
		Signature: "sig",
	})
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}
	if fakeUserDomain.getUserByIdCalls != 0 {
		t.Fatalf("UpdateProfile() called GetUserById %d times, want 0", fakeUserDomain.getUserByIdCalls)
	}
	if fakeUserDomain.updateProfileCalls != 1 {
		t.Fatalf("UpdateProfile() called UpdateUserProfile %d times, want 1", fakeUserDomain.updateProfileCalls)
	}
	if fakeUserDomain.profileUserId != 42 {
		t.Fatalf("UpdateProfile() userId = %d, want 42", fakeUserDomain.profileUserId)
	}
}

func TestChangePasswordValidationAndTokenRemoval(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("old-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		oldPassword    string
		newPassword    string
		wantErr        bool
		wantUpdate     bool
		wantTokenClear bool
	}{
		{
			name:        "old password invalid",
			oldPassword: "wrong-password",
			newPassword: "new-password",
			wantErr:     true,
		},
		{
			name:        "new password same as old",
			oldPassword: "old-password",
			newPassword: "old-password",
			wantErr:     true,
		},
		{
			name:           "success updates hash and removes token",
			oldPassword:    "old-password",
			newPassword:    "new-password",
			wantUpdate:     true,
			wantTokenClear: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origUserDomain := service.IamUserDomain()
			origSessionDomain := service.IamSessionDomain()
			defer service.RegisterIamUserDomain(origUserDomain)
			defer service.RegisterIamSessionDomain(origSessionDomain)

			fakeUserDomain := &fakeIamUserDomain{passwordHash: string(hash)}
			fakeSessionDomain := &fakeIamSessionDomain{}
			service.RegisterIamUserDomain(fakeUserDomain)
			service.RegisterIamSessionDomain(fakeSessionDomain)

			err := NewIamUserUsecase().ChangePassword(context.Background(), &model.Actor{Id: 42}, iamin.UserPasswordChangeInp{
				OldPassword: tt.oldPassword,
				NewPassword: tt.newPassword,
			})
			if tt.wantErr && err == nil {
				t.Fatal("ChangePassword() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("ChangePassword() error = %v", err)
			}
			if got := fakeUserDomain.updatePasswordHashCalls > 0; got != tt.wantUpdate {
				t.Fatalf("UpdatePasswordHash called = %v, want %v", got, tt.wantUpdate)
			}
			if got := fakeSessionDomain.removeTokenCalls > 0; got != tt.wantTokenClear {
				t.Fatalf("RemoveToken called = %v, want %v", got, tt.wantTokenClear)
			}
			if tt.wantUpdate {
				if err := bcrypt.CompareHashAndPassword([]byte(fakeUserDomain.updatedPasswordHash), []byte(tt.newPassword)); err != nil {
					t.Fatalf("updated password hash does not match new password: %v", err)
				}
				if fakeSessionDomain.removedUserKey != "42" {
					t.Fatalf("RemoveToken userKey = %q, want 42", fakeSessionDomain.removedUserKey)
				}
			}
		})
	}
}

func TestEnsureCanAuthenticateRejectsPendingUser(t *testing.T) {
	err := NewIamUserUsecase().EnsureCanAuthenticate(context.Background(), &entity.IamUser{
		Id:     42,
		Status: consts.IamUserStatusPending,
	})
	if err == nil {
		t.Fatal("EnsureCanAuthenticate() error = nil, want pending user rejection")
	}
}

func TestTemporaryTokenIsRandomAndDigestible(t *testing.T) {
	usecase := NewIamUserUsecase()
	first, err := usecase.newTemporaryToken()
	if err != nil {
		t.Fatal(err)
	}
	second, err := usecase.newTemporaryToken()
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("newTemporaryToken() returned duplicate tokens")
	}
	if len(first) != 43 {
		t.Fatalf("newTemporaryToken() length = %d, want 43", len(first))
	}
	if digest := usecase.temporaryTokenDigest(first); len(digest) != 64 {
		t.Fatalf("temporaryTokenDigest() length = %d, want 64", len(digest))
	}
}

type fakeIamUserDomain struct {
	service.IIamUserDomain

	passwordHash            string
	updatedPasswordHash     string
	getUserByIdCalls        int
	updateProfileCalls      int
	updatePasswordHashCalls int
	profileUserId           uint64
}

func (f *fakeIamUserDomain) GetUserById(ctx context.Context, id uint64) (*entity.IamUser, error) {
	f.getUserByIdCalls++
	return nil, nil
}

func (f *fakeIamUserDomain) GetUserPasswordHash(ctx context.Context, userId uint64) (string, error) {
	return f.passwordHash, nil
}

func (f *fakeIamUserDomain) UpdateUserProfile(ctx context.Context, userId uint64, avatar string, info string, signature string) error {
	f.updateProfileCalls++
	f.profileUserId = userId
	return nil
}

func (f *fakeIamUserDomain) UpdatePasswordHash(ctx context.Context, userId uint64, passwordHash string) error {
	f.updatePasswordHashCalls++
	f.updatedPasswordHash = passwordHash
	return nil
}

type fakeIamSessionDomain struct {
	service.IIamSessionDomain

	removeTokenCalls int
	removedUserKey   string
}

func (f *fakeIamSessionDomain) GetGFToken() gtoken.Token {
	return nil
}

func (f *fakeIamSessionDomain) GetGFMiddleware() gtoken.Middleware {
	return gtoken.Middleware{}
}

func (f *fakeIamSessionDomain) GenerateToken(ctx context.Context, userKey string, data any) (string, error) {
	return "", nil
}

func (f *fakeIamSessionDomain) RemoveToken(ctx context.Context, userKey string) error {
	f.removeTokenCalls++
	f.removedUserKey = userKey
	return nil
}
