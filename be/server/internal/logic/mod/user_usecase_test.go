package mod

import (
	"strings"
	"testing"

	"server/internal/consts"
)

func TestGetModDenyPermissionsForumBannedCoversForumDomain(t *testing.T) {
	s := NewModUserUsecase()

	got := s.getModDenyPermissions(consts.ModUserTypeForumBanned)
	want := []string{
		"-" + consts.IamPermissionForumTopicRead,
		"-" + consts.IamPermissionForumTopicCreate,
		"-" + consts.IamPermissionForumTopicUpdate,
		"-" + consts.IamPermissionForumReplyRead,
		"-" + consts.IamPermissionForumReplyCreate,
		"-" + consts.IamPermissionForumReplyUpdate,
	}

	if len(got) != len(want) {
		t.Fatalf("len(getModDenyPermissions()) = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("getModDenyPermissions()[%d] = %q, want %q", i, got[i], want[i])
		}
		if strings.HasPrefix(got[i], "-admin:forum/") {
			t.Fatalf("forum_banned should not deny admin permission: %q", got[i])
		}
	}
}
