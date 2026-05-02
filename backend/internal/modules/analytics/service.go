package analytics

import (
	"context"
	"errors"
	"strings"
	"time"

	"tracklink/internal/modules/links"
)

const defaultRecentLinksLimit = 5

var ErrValidation = errors.New("validation failed")

type DashboardRepository interface {
	SumTotalClicks(ctx context.Context, ownerID string) (int64, error)
	CountClicksSince(ctx context.Context, ownerID string, since time.Time) (int64, error)
	ListRecentLinks(ctx context.Context, ownerID string, limit int) ([]links.Link, error)
}

type Service struct {
	repo      DashboardRepository
	publicURL string
	now       func() time.Time
}

func NewService(repo DashboardRepository, publicURL string) *Service {
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
		TotalLinks:   0,
		ActiveLinks:  0,
		TotalClicks:  totalClicks,
		ClicksLast24: clicksLast24h,
		RecentLinks:  recent,
	}, nil, nil
}
