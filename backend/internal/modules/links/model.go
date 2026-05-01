package links

import "time"

const (
	StatusActive = "active"
)

type Link struct {
	ID           string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerID      string     `gorm:"column:owner_id;type:uuid;not null;index"`
	Code         string     `gorm:"column:code;type:text;not null;uniqueIndex:idx_links_code"`
	CustomAlias  *string    `gorm:"column:custom_alias;type:text"`
	TargetURL    string     `gorm:"column:target_url;type:text;not null"`
	Status       string     `gorm:"column:status;type:text;not null;default:'active'"`
	TotalClicks  int64      `gorm:"column:total_clicks;not null;default:0"`
	LastClickedAt *time.Time `gorm:"column:last_clicked_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;not null;default:now()"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (Link) TableName() string {
	return "links"
}
