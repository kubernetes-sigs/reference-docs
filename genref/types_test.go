package main

import (
	"testing"

	"k8s.io/gengo/types"
)

func TestFieldName(t *testing.T) {
	tests := []struct {
		name       string
		member     apiMember
		wantName   string
		wantInline bool
	}{
		{
			name: "single option omitempty",
			member: apiMember{types.Member{
				Name: "LastUpdateTime",
				Tags: `json:"lastUpdateTime,omitempty"`,
			}},
			wantName:   "lastUpdateTime",
			wantInline: false,
		},
		{
			name: "multiple options omitempty,omitzero",
			member: apiMember{types.Member{
				Name: "LastUpdateTime",
				Tags: `json:"lastUpdateTime,omitempty,omitzero"`,
			}},
			wantName:   "lastUpdateTime",
			wantInline: false,
		},
		{
			name: "single option omitzero",
			member: apiMember{types.Member{
				Name: "LastUpdateTime",
				Tags: `json:"lastUpdateTime,omitzero"`,
			}},
			wantName:   "lastUpdateTime",
			wantInline: false,
		},
		{
			name: "multiple options omitzero,omitempty",
			member: apiMember{types.Member{
				Name: "LastUpdateTime",
				Tags: `json:"lastUpdateTime,omitzero,omitempty"`,
			}},
			wantName:   "lastUpdateTime",
			wantInline: false,
		},
		{
			name: "no options",
			member: apiMember{types.Member{
				Name: "LastUpdateTime",
				Tags: `json:"lastUpdateTime"`,
			}},
			wantName:   "lastUpdateTime",
			wantInline: false,
		},
		{
			name: "empty tag falls back to member name",
			member: apiMember{types.Member{
				Name: "LastUpdateTime",
				Tags: ``,
			}},
			wantName:   "LastUpdateTime",
			wantInline: false,
		},
		{
			name: "inline tag without field name",
			member: apiMember{types.Member{
				Name: "TypeMeta",
				Tags: `json:",inline"`,
			}},
			wantName:   "TypeMeta",
			wantInline: true,
		},
		{
			name: "inline with omitempty",
			member: apiMember{types.Member{
				Name: "TypeMeta",
				Tags: `json:",inline,omitempty"`,
			}},
			wantName:   "TypeMeta",
			wantInline: true,
		},
		{
			name: "inline with field name and multiple options",
			member: apiMember{types.Member{
				Name: "TypeMeta",
				Tags: `json:"meta,inline,omitempty"`,
			}},
			wantName:   "meta",
			wantInline: true,
		},
		{
			name: "ignored field dash",
			member: apiMember{types.Member{
				Name: "IgnoredField",
				Tags: `json:"-"`,
			}},
			wantName:   "-",
			wantInline: false,
		},
		{
			name: "field named dash -,-",
			member: apiMember{types.Member{
				Name: "DashField",
				Tags: `json:"-,"`,
			}},
			wantName:   "-",
			wantInline: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.member.FieldName(); got != tc.wantName {
				t.Errorf("FieldName() = %q, want %q", got, tc.wantName)
			}
			if got := tc.member.IsInline(); got != tc.wantInline {
				t.Errorf("IsInline() = %v, want %v", got, tc.wantInline)
			}
		})
	}
}
