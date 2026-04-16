package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"time-tracker/internal/config"
	"time-tracker/internal/middleware"
)

// RegisterRoutes đăng ký toàn bộ routes của ứng dụng
func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// ── Public routes (không cần auth) ──────────────────────
	public := r.Group("/")
	{
		public.GET("/login", func(c *gin.Context) { c.HTML(200, "auth/login.html", nil) })
		public.POST("/login", func(c *gin.Context) { /* TODO: Sprint 1 */ })
		public.GET("/register", func(c *gin.Context) { c.HTML(200, "auth/register.html", nil) })
		public.POST("/register", func(c *gin.Context) { /* TODO: Sprint 1 */ })
		public.GET("/forgot-password", func(c *gin.Context) { c.HTML(200, "auth/forgot-password.html", nil) })
		public.POST("/forgot-password", func(c *gin.Context) { /* TODO: Sprint 1 */ })
		public.GET("/reset-password", func(c *gin.Context) { /* TODO: Sprint 1 */ })
		public.POST("/reset-password", func(c *gin.Context) { /* TODO: Sprint 1 */ })
		public.GET("/verify-email", func(c *gin.Context) { /* TODO: Sprint 1 */ })
	}

	// ── Protected routes (cần auth) ─────────────────────────
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth(db, cfg))
	{
		// Dashboard
		protected.GET("/", func(c *gin.Context) { /* TODO: Sprint 2 */ })

		// Activities
		protected.POST("/activities", func(c *gin.Context) { /* TODO: Sprint 2 */ })
		protected.PUT("/activities/:id", func(c *gin.Context) { /* TODO: Sprint 2 */ })
		protected.DELETE("/activities/:id", func(c *gin.Context) { /* TODO: Sprint 2 */ })

		// Categories
		protected.GET("/categories", func(c *gin.Context) { /* TODO: Sprint 3 */ })
		protected.POST("/categories", func(c *gin.Context) { /* TODO: Sprint 3 */ })
		protected.PUT("/categories/:id", func(c *gin.Context) { /* TODO: Sprint 3 */ })
		protected.DELETE("/categories/:id", func(c *gin.Context) { /* TODO: Sprint 3 */ })

		// Reports
		protected.GET("/reports/daily", func(c *gin.Context) { /* TODO: Sprint 4 */ })
		protected.GET("/reports/weekly", func(c *gin.Context) { /* TODO: Sprint 4 */ })
		protected.GET("/reports/monthly", func(c *gin.Context) { /* TODO: Sprint 4 */ })

		// Trash
		protected.GET("/trash", func(c *gin.Context) { /* TODO: Sprint 5 */ })
		protected.POST("/trash/:id/restore", func(c *gin.Context) { /* TODO: Sprint 5 */ })
		protected.DELETE("/trash/:id", func(c *gin.Context) { /* TODO: Sprint 5 */ })

		// Export
		protected.GET("/export", func(c *gin.Context) { /* TODO: Sprint 5 */ })
		protected.POST("/export", func(c *gin.Context) { /* TODO: Sprint 5 */ })

		// Profile & Settings
		protected.GET("/profile", func(c *gin.Context) { /* TODO: Sprint 5 */ })
		protected.PUT("/profile", func(c *gin.Context) { /* TODO: Sprint 5 */ })
		protected.POST("/auth/logout", func(c *gin.Context) { /* TODO: Sprint 1 */ })
		protected.POST("/auth/change-password", func(c *gin.Context) { /* TODO: Sprint 1 */ })
	}
}
