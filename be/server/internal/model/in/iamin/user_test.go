package iamin

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/util/gvalid"
)

func TestUserCreateUsername(t *testing.T) {
	for _, tt := range []struct {
		username string
		valid    bool
	}{
		{"abc", true}, {"User123", true}, {"123", true},
		{strings.Repeat("a", 20), true},
		{"", false}, {"ab", false}, {strings.Repeat("a", 21), false},
		{"user_name", false}, {"user-name", false}, {"user.name", false},
		{"a b", false}, {" abc", false}, {"abc ", false},
		{"abc\n", false}, {"abc\t", false}, {"\u7528\u6237\u540d", false},
		{"\uff41\uff42\uff43", false},
	} {
		t.Run(tt.username, func(t *testing.T) {
			in := UserCreateInp{Username: tt.username, Email: "user@example.com", Password: "password123"}
			err := gvalid.New().Data(in).Run(context.Background())
			if (err == nil) != tt.valid {
				t.Fatalf("username %q: valid=%v, error=%v", tt.username, tt.valid, err)
			}
		})
	}
}
