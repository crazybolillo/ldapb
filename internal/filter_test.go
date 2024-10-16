package internal

import (
	"github.com/crazybolillo/eryth/pkg/model"
	"testing"
)

func TestParseFilter(t *testing.T) {
	cases := []struct {
		name   string
		filter string
		want   model.ContactPageFilter
	}{
		{
			"onlyCn",
			"(cn=Churchill)",
			model.ContactPageFilter{
				Name: "Churchill",
			},
		},
		{
			"or",
			"(|(cn=1001*)(telephoneNumber=1001*))",
			model.ContactPageFilter{
				Name:     "1001*",
				Phone:    "1001*",
				Operator: "or",
			},
		},
		{
			"onlyPhone",
			"(telephoneNumber=*7500)",
			model.ContactPageFilter{
				Phone: "*7500",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFilter(tt.filter)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
