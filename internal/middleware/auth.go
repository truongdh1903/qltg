package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"time-tracker/internal/config"
)

const UserContextKey = "current_user_id"

// RequireAuth kiểm tra session cookie hợp lệ.
// Nếu không hợp lệ → redirect về /login.
func RequireAuth(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Sprint 1 — implement session lookup
		// 1. Đọc cookie cfg.SessionCookieName
		// 2. Hash raw token → SHA-256
		// 3. Query user_sessions: WHERE session_token_hash = ? AND expires_at > NOW()
		// 4. Nếu không tìm thấy → redirect /login
		// 5. Set c.Set(UserContextKey, session.UserID)
		// 6. Update last_active_at
		// 7. c.Next()

		// Placeholder: luôn redirect (sẽ implement Sprint 1)
		token, err := c.Cookie(cfg.SessionCookieName)
		if err != nil || token == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequestID thêm unique request ID vào mỗi request (dùng cho logging)
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement với google/uuid
		c.Next()
	}
}
