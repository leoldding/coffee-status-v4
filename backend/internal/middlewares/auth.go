package middlewares

import (
	"net/http"
)

func AuthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token_leoding")
		if err != nil {
			http.Error(w, "unauthorized: missing auth cookie", http.StatusUnauthorized)
			return
		}

		req, err := http.NewRequest("GET", "/api/v1/auth/check", nil)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		req.AddCookie(cookie)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
