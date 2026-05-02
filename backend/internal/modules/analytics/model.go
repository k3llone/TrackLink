package analytics

import "time"

type ClickEvent struct {
	ID        string     `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	LinkID    string     `gorm:"column:link_id;type:uuid;not null;index"`
	ClickedAt time.Time  `gorm:"column:clicked_at;not null"`
	Referrer  *string    `gorm:"column:referrer;type:text"`
	UserAgent *string    `gorm:"column:user_agent;type:text"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;default:now()"`
}

func (ClickEvent) TableName() string {
	return "click_events"
}
