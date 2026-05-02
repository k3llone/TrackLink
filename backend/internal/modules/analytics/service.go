package analytics

import (
	"context"
	"errors"
	"strings"
	"time"

	"tracklink/internal/modules/links"
)

const (
	defaultRecentLinksLimit  = 5
	defaultRecentClicksLimit = 20
	maxRecentClicksLimit     = 100
	defaultAnalyticsPeriod   = 7 * 24 * time.Hour
	defaultAnalyticsGroupBy  = GroupByDay
	GroupByHour              = "hour"
	GroupByDay               = "day"
)

var ErrValidation = errors.New("validation failed")
var ErrLinkNotFound = errors.New("link not found")

type TimeSeriesBucket struct {
	PeriodStart time.Time
	Clicks      int64
}

type ServiceRepository interface {
	CountTotalLinks(ctx context.Context, ownerID string) (int64, error)
	CountActiveLinks(ctx context.Context, ownerID string) (int64, error)
	SumTotalClicks(ctx context.Context, ownerID string) (int64, error)
	CountClicksSince(ctx context.Context, ownerID string, since time.Time) (int64, error)
	ListRecentLinks(ctx context.Context, ownerID string, limit int) ([]links.Link, error)
	GetLinkByIDAndOwner(ctx context.Context, linkID, ownerID string) (links.Link, error)
	CountLinkClicks(ctx context.Context, linkID string, from, to time.Time) (int64, error)
	CountLinkClicksSince(ctx context.Context, linkID string, since time.Time) (int64, error)
	LastLinkClickedAt(ctx context.Context, linkID string, from, to time.Time) (*time.Time, error)
	ListLinkClickSeries(ctx context.Context, linkID string, from, to time.Time, groupBy string) ([]TimeSeriesBucket, error)
	ListRecentClicks(ctx context.Context, linkID string, limit int) ([]ClickEvent, error)
}

type Service struct {
	repo      ServiceRepository
	publicURL string
	now       func() time.Time
}

func NewService(repo ServiceRepository, publicURL string) *Service {
	return &Service{
		repo:      repo,
		publicURL: strings.TrimRight(strings.TrimSpace(publicURL), "/"),
		now:       time.Now,
	}
}

func (s *Service) LoadDashboard(ctx context.Context, userID string) (DashboardResponse, map[string]string, error) {
	fields := map[string]string{}
	ownerID := strings.TrimSpace(userID)
	if ownerID == "" {
		fields["ownerId"] = "Owner ID is required"
	}
	if s.repo == nil {
		fields["repository"] = "Repository is required"
	}
	if len(fields) > 0 {
		return DashboardResponse{}, fields, ErrValidation
	}

	totalLinks, err := s.repo.CountTotalLinks(ctx, ownerID)
	if err != nil {
		return DashboardResponse{}, nil, err
	}
	activeLinks, err := s.repo.CountActiveLinks(ctx, ownerID)
	if err != nil {
		return DashboardResponse{}, nil, err
	}
	totalClicks, err := s.repo.SumTotalClicks(ctx, ownerID)
	if err != nil {
		return DashboardResponse{}, nil, err
	}
	clicksLast24h, err := s.repo.CountClicksSince(ctx, ownerID, s.now().UTC().Add(-24*time.Hour))
	if err != nil {
		return DashboardResponse{}, nil, err
	}
	recentLinks, err := s.repo.ListRecentLinks(ctx, ownerID, defaultRecentLinksLimit)
	if err != nil {
		return DashboardResponse{}, nil, err
	}

	recent := make([]links.LinkResponse, 0, len(recentLinks))
	for _, link := range recentLinks {
		recent = append(recent, links.MapLinkToResponse(link, s.publicURL))
	}

	return DashboardResponse{
		TotalLinks:   totalLinks,
		ActiveLinks:  activeLinks,
		TotalClicks:  totalClicks,
		ClicksLast24: clicksLast24h,
		RecentLinks:  recent,
	}, nil, nil
}

func (s *Service) LoadLinkAnalytics(ctx context.Context, userID, linkID string, query LinkAnalyticsQuery) (LinkAnalyticsResponse, map[string]string, error) {
	fields := map[string]string{}
	ownerID := strings.TrimSpace(userID)
	normalizedLinkID := strings.TrimSpace(linkID)
	if ownerID == "" {
		fields["ownerId"] = "Owner ID is required"
	}
	if normalizedLinkID == "" {
		fields["linkId"] = "Link ID is required"
	}
	if s.repo == nil {
		fields["repository"] = "Repository is required"
	}

	groupBy := strings.TrimSpace(query.GroupBy)
	if groupBy == "" {
		groupBy = defaultAnalyticsGroupBy
	}
	if groupBy != GroupByHour && groupBy != GroupByDay {
		fields["groupBy"] = "Group by must be one of: hour, day"
	}

	now := s.now().UTC()
	to := query.To.UTC()
	if to.IsZero() {
		to = now
	}
	from := query.From.UTC()
	if from.IsZero() {
		from = to.Add(-defaultAnalyticsPeriod)
	}
	if from.After(to) {
		fields["from"] = "From must be before or equal to to"
	}

	if len(fields) > 0 {
		return LinkAnalyticsResponse{}, fields, ErrValidation
	}

	link, err := s.repo.GetLinkByIDAndOwner(ctx, normalizedLinkID, ownerID)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) || errors.Is(err, links.ErrLinkNotFound) {
			return LinkAnalyticsResponse{}, nil, ErrLinkNotFound
		}
		return LinkAnalyticsResponse{}, nil, err
	}

	totalClicks, err := s.repo.CountLinkClicks(ctx, link.ID, from, to)
	if err != nil {
		return LinkAnalyticsResponse{}, nil, err
	}
	clicksLast24h, err := s.repo.CountLinkClicksSince(ctx, link.ID, now.Add(-24*time.Hour))
	if err != nil {
		return LinkAnalyticsResponse{}, nil, err
	}
	lastClickedAt, err := s.repo.LastLinkClickedAt(ctx, link.ID, from, to)
	if err != nil {
		return LinkAnalyticsResponse{}, nil, err
	}
	seriesBuckets, err := s.repo.ListLinkClickSeries(ctx, link.ID, from, to, groupBy)
	if err != nil {
		return LinkAnalyticsResponse{}, nil, err
	}

	series := make([]TimeSeriesPoint, 0, len(seriesBuckets))
	for _, bucket := range seriesBuckets {
		series = append(series, TimeSeriesPoint{
			PeriodStart: bucket.PeriodStart.UTC().Format(time.RFC3339),
			Clicks:      bucket.Clicks,
		})
	}

	var lastClickedAtValue *string
	if lastClickedAt != nil {
		formatted := lastClickedAt.UTC().Format(time.RFC3339)
		lastClickedAtValue = &formatted
	}

	return LinkAnalyticsResponse{
		LinkID:        link.ID,
		TotalClicks:   totalClicks,
		ClicksLast24:  clicksLast24h,
		LastClickedAt: lastClickedAtValue,
		Series:        series,
	}, nil, nil
}

func (s *Service) LoadRecentClicks(ctx context.Context, userID, linkID string, query RecentClicksQuery) (RecentClicksResponse, map[string]string, error) {
	fields := map[string]string{}
	ownerID := strings.TrimSpace(userID)
	normalizedLinkID := strings.TrimSpace(linkID)
	if ownerID == "" {
		fields["ownerId"] = "Owner ID is required"
	}
	if normalizedLinkID == "" {
		fields["linkId"] = "Link ID is required"
	}
	if s.repo == nil {
		fields["repository"] = "Repository is required"
	}

	limit := query.Limit
	if limit == 0 && !query.limitProvided {
		limit = defaultRecentClicksLimit
	}
	if limit < 1 {
		fields["limit"] = "Limit must be greater than or equal to 1"
	} else if limit > maxRecentClicksLimit {
		limit = maxRecentClicksLimit
	}

	if len(fields) > 0 {
		return RecentClicksResponse{}, fields, ErrValidation
	}

	link, err := s.repo.GetLinkByIDAndOwner(ctx, normalizedLinkID, ownerID)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) || errors.Is(err, links.ErrLinkNotFound) {
			return RecentClicksResponse{}, nil, ErrLinkNotFound
		}
		return RecentClicksResponse{}, nil, err
	}

	events, err := s.repo.ListRecentClicks(ctx, link.ID, limit)
	if err != nil {
		return RecentClicksResponse{}, nil, err
	}

	items := make([]ClickEventResponse, 0, len(events))
	for _, event := range events {
		items = append(items, ClickEventResponse{
			ID:        event.ID,
			LinkID:    event.LinkID,
			ClickedAt: event.ClickedAt.UTC().Format(time.RFC3339),
			Referrer:  event.Referrer,
			UserAgent: event.UserAgent,
		})
	}

	return RecentClicksResponse{Items: items}, nil, nil
}
