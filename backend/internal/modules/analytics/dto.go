package analytics

import "tracklink/internal/modules/links"

type DashboardResponse struct {
	TotalLinks   int64                `json:"totalLinks"`
	ActiveLinks  int64                `json:"activeLinks"`
	TotalClicks  int64                `json:"totalClicks"`
	ClicksLast24 int64                `json:"clicksLast24h"`
	RecentLinks  []links.LinkResponse `json:"recentLinks"`
}
