package admin

import (
	"context"
	"errors"
	"strings"
	"time"

	"tracklink/internal/modules/links"
)

var ErrValidation = errors.New("validation failed")
var ErrLinkNotFound = errors.New("link not found")
var ErrStatusChangeNotAllowed = errors.New("status change not allowed")

const (
	defaultListPage     = 1
	defaultListPageSize = 20
	maxListPageSize     = 100
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, adminUserID string, query ListLinksQuery) ([]links.Link, PaginationResponse, map[string]string, error) {
	fields := map[string]string{}

	if strings.TrimSpace(adminUserID) == "" {
		fields["adminUserId"] = "Admin user ID is required"
	}

	page := query.Page
	if page == 0 {
		page = defaultListPage
	}
	if page < 1 {
		fields["page"] = "Page must be greater than or equal to 1"
	}

	pageSize := query.PageSize
	if pageSize == 0 {
		pageSize = defaultListPageSize
	}
	if pageSize < 1 {
		fields["pageSize"] = "Page size must be greater than or equal to 1"
	} else if pageSize > maxListPageSize {
		pageSize = maxListPageSize
	}

	if len(fields) > 0 {
		return nil, PaginationResponse{}, fields, ErrValidation
	}

	items, totalItems, err := s.repo.List(ctx, ListLinksFilter{
		Page:     page,
		PageSize: pageSize,
		Q:        strings.TrimSpace(query.Q),
	})
	if err != nil {
		return nil, PaginationResponse{}, nil, err
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = int((totalItems + int64(pageSize) - 1) / int64(pageSize))
	}

	return items, PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil, nil
}

func (s *Service) Block(ctx context.Context, adminUserID, linkID string, req AdminBlockLinkRequest) (links.Link, map[string]string, error) {
	fields := map[string]string{}

	if strings.TrimSpace(adminUserID) == "" {
		fields["adminUserId"] = "Admin user ID is required"
	}
	if strings.TrimSpace(linkID) == "" {
		fields["linkId"] = "Link ID is required"
	}
	if len(fields) > 0 {
		return links.Link{}, fields, ErrValidation
	}

	_ = normalizeReason(req.Reason)

	link, err := s.repo.GetByID(ctx, linkID)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return links.Link{}, nil, ErrLinkNotFound
		}
		return links.Link{}, nil, err
	}

	if link.Status == links.StatusDeleted || link.DeletedAt != nil {
		return links.Link{}, nil, ErrLinkNotFound
	}

	if link.Status == links.StatusBlocked {
		return link, nil, nil
	}

	updated, err := s.repo.UpdateStatus(ctx, linkID, links.StatusBlocked)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return links.Link{}, nil, ErrLinkNotFound
		}
		return links.Link{}, nil, err
	}

	return updated, nil, nil
}

func (s *Service) Unblock(ctx context.Context, adminUserID, linkID string) (links.Link, map[string]string, error) {
	fields := map[string]string{}

	if strings.TrimSpace(adminUserID) == "" {
		fields["adminUserId"] = "Admin user ID is required"
	}
	if strings.TrimSpace(linkID) == "" {
		fields["linkId"] = "Link ID is required"
	}
	if len(fields) > 0 {
		return links.Link{}, fields, ErrValidation
	}

	link, err := s.repo.GetByID(ctx, linkID)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return links.Link{}, nil, ErrLinkNotFound
		}
		return links.Link{}, nil, err
	}

	if link.Status == links.StatusDeleted || link.DeletedAt != nil {
		return links.Link{}, nil, ErrLinkNotFound
	}

	if link.Status != links.StatusBlocked {
		return link, nil, nil
	}

	updated, err := s.repo.UpdateStatus(ctx, linkID, links.StatusActive)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return links.Link{}, nil, ErrLinkNotFound
		}
		return links.Link{}, nil, err
	}

	return updated, nil, nil
}

func (s *Service) Deactivate(ctx context.Context, adminUserID, linkID string) (links.Link, map[string]string, error) {
	fields := map[string]string{}

	if strings.TrimSpace(adminUserID) == "" {
		fields["adminUserId"] = "Admin user ID is required"
	}
	if strings.TrimSpace(linkID) == "" {
		fields["linkId"] = "Link ID is required"
	}
	if len(fields) > 0 {
		return links.Link{}, fields, ErrValidation
	}

	link, err := s.repo.GetByID(ctx, linkID)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return links.Link{}, nil, ErrLinkNotFound
		}
		return links.Link{}, nil, err
	}

	if link.Status == links.StatusDeleted || link.DeletedAt != nil {
		return links.Link{}, nil, ErrLinkNotFound
	}

	if link.Status == links.StatusBlocked {
		return links.Link{}, nil, ErrStatusChangeNotAllowed
	}

	if link.Status == links.StatusInactive {
		return link, nil, nil
	}

	updated, err := s.repo.UpdateStatus(ctx, linkID, links.StatusInactive)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return links.Link{}, nil, ErrLinkNotFound
		}
		return links.Link{}, nil, err
	}

	return updated, nil, nil
}

func normalizeReason(reason *string) *string {
	if reason == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*reason)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func mapAdminLink(link links.Link, publicURL string) AdminLink {
	return AdminLink(links.MapLinkToResponse(link, publicURL))
}

func mapPagination(pagination PaginationResponse) PaginationResponse {
	return PaginationResponse{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalItems: pagination.TotalItems,
		TotalPages: pagination.TotalPages,
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
