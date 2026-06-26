package site

import (
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func TestEncodeConfigValueWrapsValAndPreservesInputType(t *testing.T) {
	s := NewSiteConfigDomain()
	tests := []struct {
		name  string
		value string
		want  any
	}{
		{
			name:  "string",
			value: "NextPT",
			want:  "NextPT",
		},
		{
			name:  "int",
			value: "1800",
			want:  float64(1800),
		},
		{
			name:  "bool",
			value: "true",
			want:  true,
		},
		{
			name:  "quoted string",
			value: `"1800"`,
			want:  "1800",
		},
		{
			name:  "object",
			value: `{"enabled":true}`,
			want: map[string]any{
				"enabled": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, err := gjson.LoadJson([]byte(s.encodeConfigValue(tt.value)))
			if err != nil {
				t.Fatalf("LoadJson error = %v", err)
			}

			got := j.Get(siteConfigValueField).Val()
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("encoded val = %#v, want %#v", got, tt.want)
			}
		})
	}
}
