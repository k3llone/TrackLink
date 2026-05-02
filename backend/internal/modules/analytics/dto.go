package analytics

import (
	"time"

	"tracklink/internal/modules/links"
)

type DashboardResponse struct {
	TotalLinks   int64                `json:"totalLinks"`
	ActiveLinks  int64                `json:"activeLinks"`
	TotalClicks  int64                `json:"totalClicks"`
	ClicksLast24 int64                `json:"clicksLast24h"`
	RecentLinks  []links.LinkResponse `json:"recentLinks"`
}

type LinkAnalyticsQuery struct {
	From    time.Time
	To      time.Time
	GroupBy string
}

type LinkAnalyticsResponse struct {
	LinkID        string            `json:"linkId"`
	TotalClicks   int64             `json:"totalClicks"`
	ClicksLast24  int64             `json:"clicksLast24h"`
	LastClickedAt *string           `json:"lastClickedAt"`
	Series        []TimeSeriesPoint `json:"series"`
}

type TimeSeriesPoint struct {
	PeriodStart string `json:"periodStart"`
	Clicks      int64  `json:"clicks"`
}
