package mw

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/sakithb/hcblk-server/internal/models"
)

func Authentication(sm *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, ok := sm.Get(r.Context(), "user").(models.User)
			if !ok {
				next.ServeHTTP(w, r)
			} else {
				ctx := context.WithValue(r.Context(), "user", u)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}
