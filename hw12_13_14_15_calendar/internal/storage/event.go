package storage

import "time"

type Event struct {
	ID           uint64        `db:"event_id"`
	Title        string        `db:"title"`
	StartDate    time.Time     `db:"start_date"`
	EndDate      time.Time     `db:"end_date"`
	Description  string        `db:"description"`
	UserID       uint64        `db:"user_id"`
	NotifyBefore time.Duration `db:"notify_before"`
}
