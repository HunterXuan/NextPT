package iam

import (
	"context"
	"testing"
)

func TestCheckPermissionWithListDenyOverridesRoleAndUserAllow(t *testing.T) {
	s := NewIamPermissionDomain()

	tests := []struct {
		name      string
		rolePerms []string
		userAcls  []string
		permKey   string
		want      bool
	}{
		{
			name:      "role exact allow",
			rolePerms: []string{"read:forum/topic:*"},
			permKey:   "read:forum/topic:123",
			want:      true,
		},
		{
			name:     "user exact allow",
			userAcls: []string{"download:catalog/torrent:*"},
			permKey:  "download:catalog/torrent:1",
			want:     true,
		},
		{
			name:      "user deny overrides role allow",
			rolePerms: []string{"download:catalog/torrent:*"},
			userAcls:  []string{"-download:catalog/torrent:*"},
			permKey:   "download:catalog/torrent:1",
			want:      false,
		},
		{
			name:     "user deny overrides earlier user allow",
			userAcls: []string{"download:catalog/torrent:*", "-download:catalog/torrent:*"},
			permKey:  "download:catalog/torrent:1",
			want:     false,
		},
		{
			name:      "global deny overrides global role allow",
			rolePerms: []string{"*"},
			userAcls:  []string{"-*"},
			permKey:   "read:forum/topic:1",
			want:      false,
		},
		{
			name:    "missing permission",
			permKey: "update:forum/topic:1",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CheckPermissionWithList(context.Background(), tt.rolePerms, tt.userAcls, tt.permKey)
			if err != nil {
				t.Fatalf("CheckPermissionWithList() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("CheckPermissionWithList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPermissionWithListForumBannedAcls(t *testing.T) {
	s := NewIamPermissionDomain()
	forumDenyAcls := []string{
		"-read:forum/topic:*",
		"-create:forum/topic:*",
		"-update:forum/topic:*",
		"-read:forum/reply:*",
		"-create:forum/reply:*",
		"-update:forum/reply:*",
	}

	tests := []struct {
		name      string
		rolePerms []string
		permKey   string
		want      bool
	}{
		{
			name:      "blocks topic read",
			rolePerms: []string{"read:forum/topic:*"},
			permKey:   "read:forum/topic:1",
			want:      false,
		},
		{
			name:      "blocks topic create",
			rolePerms: []string{"create:forum/topic:*"},
			permKey:   "create:forum/topic:*",
			want:      false,
		},
		{
			name:      "blocks topic update",
			rolePerms: []string{"update:forum/topic:*"},
			permKey:   "update:forum/topic:1",
			want:      false,
		},
		{
			name:      "blocks reply read",
			rolePerms: []string{"read:forum/reply:*"},
			permKey:   "read:forum/reply:1",
			want:      false,
		},
		{
			name:      "blocks reply create",
			rolePerms: []string{"create:forum/reply:*"},
			permKey:   "create:forum/reply:*",
			want:      false,
		},
		{
			name:      "blocks reply update",
			rolePerms: []string{"update:forum/reply:*"},
			permKey:   "update:forum/reply:1",
			want:      false,
		},
		{
			name:      "does not block admin forum permission",
			rolePerms: []string{"admin:forum/topic:*"},
			permKey:   "admin:forum/topic:1",
			want:      true,
		},
		{
			name:      "does not block catalog permission",
			rolePerms: []string{"read:catalog/torrent:*"},
			permKey:   "read:catalog/torrent:1",
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CheckPermissionWithList(context.Background(), tt.rolePerms, forumDenyAcls, tt.permKey)
			if err != nil {
				t.Fatalf("CheckPermissionWithList() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("CheckPermissionWithList() = %v, want %v", got, tt.want)
			}
		})
	}
}
