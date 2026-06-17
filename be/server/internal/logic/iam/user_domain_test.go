package iam

import (
	"testing"

	"server/internal/dao"
)

func TestNormalizeAdminUserOrderWhitelist(t *testing.T) {
	s := NewIamUserDomain()
	columns := dao.IamUser.Columns()
	defaultOrder := columns.Id + " DESC"

	tests := []struct {
		name  string
		order string
		want  string
	}{
		{
			name:  "empty uses default",
			order: "",
			want:  defaultOrder,
		},
		{
			name:  "allowed field asc",
			order: "username asc",
			want:  columns.Username + " ASC",
		},
		{
			name:  "allowed snake field desc",
			order: "created_at desc",
			want:  columns.CreatedAt + " DESC",
		},
		{
			name:  "allowed camel field asc",
			order: "vipUntil asc",
			want:  columns.VipUntil + " ASC",
		},
		{
			name:  "unknown field uses default",
			order: "password_hash asc",
			want:  defaultOrder,
		},
		{
			name:  "removed enabled field uses default",
			order: "enabled asc",
			want:  defaultOrder,
		},
		{
			name:  "invalid direction uses default",
			order: "email sideways",
			want:  defaultOrder,
		},
		{
			name:  "extra tokens use default",
			order: "username desc, id asc",
			want:  defaultOrder,
		},
		{
			name:  "punctuated field uses default",
			order: "username;drop asc",
			want:  defaultOrder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.normalizeAdminUserOrder(tt.order)
			if got != tt.want {
				t.Fatalf("normalizeAdminUserOrder(%q) = %q, want %q", tt.order, got, tt.want)
			}
		})
	}
}
