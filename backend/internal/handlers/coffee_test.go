package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeStatus(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
		err  bool
	}{
		{name: "accepts yes", body: `{"value":"yes"}`, want: "yes"},
		{name: "accepts otw", body: `{"value":"otw"}`, want: "otw"},
		{name: "accepts no", body: `{"value":"no"}`, want: "no"},
		{name: "rejects unknown field", body: `{"value":"yes","extra":true}`, err: true},
		{name: "rejects trailing JSON", body: `{"value":"yes"} {}`, err: true},
		{name: "rejects invalid status", body: `{"value":"maybe"}`, err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			got, err := decodeStatus(r)
			if (err != nil) != tt.err {
				t.Fatalf("decodeStatus() error = %v, want error %t", err, tt.err)
			}
			if got.Value != tt.want {
				t.Errorf("decodeStatus() value = %q, want %q", got.Value, tt.want)
			}
		})
	}
}

func TestDecodeStatusRejectsLargeBody(t *testing.T) {
	body := `{"value":"yes"}` + strings.Repeat(" ", maxRequestBodyBytes)
	r := httptest.NewRequest("POST", "/", strings.NewReader(body))

	if _, err := decodeStatus(r); err == nil {
		t.Fatal("decodeStatus() accepted a body larger than 16 KiB")
	}
}
