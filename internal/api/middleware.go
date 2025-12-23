package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ContextKey string

const (
	TenantIDKey ContextKey = "tenant_id"
	UserIDKey   ContextKey = "user_id"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// skip auth for webhooks
		if r.URL.Path == "/webhooks/plaid" || r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		tenantIDStr := r.Header.Get("X-Tenant-ID")
		if tenantIDStr == "" {
			tenantIDStr = "00000000-0000-0000-0000-000000000001"
		}

		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
		ctx = context.WithValue(ctx, UserIDKey, tenantID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetTenantID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(TenantIDKey).(uuid.UUID)
	return id, ok
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}
