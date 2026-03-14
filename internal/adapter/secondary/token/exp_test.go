//nolint:testpackage // tests cover unexported expOK helper with all supported claim types.
package token

import (
	"encoding/json"
	"testing"
	"time"
)

func TestExpOK(t *testing.T) {
	future := time.Now().Add(time.Minute).Unix()
	past := time.Now().Add(-time.Minute).Unix()

	tests := []struct {
		name string
		v    any
		want bool
	}{
		{name: "nil", v: nil, want: false},
		{name: "float64 future", v: float64(future), want: true},
		{name: "float64 past", v: float64(past), want: false},
		{name: "int64 future", v: future, want: true},
		{name: "int64 past", v: past, want: false},
		{name: "int future", v: int(future), want: true},
		{name: "int past", v: int(past), want: false},
		{name: "json number future", v: json.Number("4102444800"), want: true},
		{name: "json number invalid", v: json.Number("abc"), want: false},
		{name: "string future", v: "4102444800", want: true},
		{name: "string invalid", v: "abc", want: false},
		{name: "unsupported type", v: true, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expOK(tt.v)
			if got != tt.want {
				t.Fatalf("expOK(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}
