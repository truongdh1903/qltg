package scheduler

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Scheduler struct {
	cron *cron.Cron
	db   *gorm.DB
}

func New(db *gorm.DB) *Scheduler {
	// Dùng timezone UTC+7
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	c := cron.New(cron.WithLocation(loc))
	return &Scheduler{cron: c, db: db}
}

func (s *Scheduler) Start() {
	// Trash cleanup: chạy mỗi ngày lúc 2:00 AM UTC+7 (BR-38)
	s.cron.AddFunc("0 2 * * *", s.cleanupTrash)

	// Session cleanup: chạy mỗi ngày lúc 3:00 AM
	s.cron.AddFunc("0 3 * * *", s.cleanupExpiredSessions)

	s.cron.Start()
	log.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

// cleanupTrash xóa vĩnh viễn activities đã ở Trash > 30 ngày (BR-38)
func (s *Scheduler) cleanupTrash() {
	cutoff := time.Now().UTC().Add(-30 * 24 * time.Hour)
	result := s.db.Exec(
		"DELETE FROM activities WHERE deleted_at IS NOT NULL AND deleted_at < ?",
		cutoff,
	)
	if result.Error != nil {
		log.Printf("Trash cleanup error: %v", result.Error)
		return
	}
	log.Printf("Trash cleanup: deleted %d records", result.RowsAffected)
}

// cleanupExpiredSessions xóa sessions đã hết hạn
func (s *Scheduler) cleanupExpiredSessions() {
	result := s.db.Exec("DELETE FROM user_sessions WHERE expires_at < NOW()")
	if result.Error != nil {
		log.Printf("Session cleanup error: %v", result.Error)
		return
	}
	log.Printf("Session cleanup: deleted %d expired sessions", result.RowsAffected)
}
