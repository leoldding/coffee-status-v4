package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSAllowsOnlyConfiguredDevelopmentOrigin(t *testing.T) {
	handler := CORS("http://localhost:5173")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name   string
		origin string
		want   string
	}{
		{name: "configured origin", origin: "http://localhost:5173", want: "http://localhost:5173"},
		{name: "other localhost port", origin: "http://localhost:3000"},
		{name: "former suffix origin", origin: "https://attacker.leoding.com"},
		{name: "production same origin request", origin: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Origin", tt.origin)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			if got := w.Header().Get("Access-Control-Allow-Origin"); got != tt.want {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, tt.want)
			}
		})
	}
}
