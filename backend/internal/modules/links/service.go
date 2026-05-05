package links

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"tracklink/internal/shared"
)

var ErrValidation = errors.New("validation failed")
var ErrCodeGenerationExhausted = errors.New("code generation attempts exhausted")
var ErrAliasAlreadyExists = errors.New("custom alias already exists")
var ErrLinkNotFound = errors.New("link not found")
var ErrStatusChangeNotAllowed = errors.New("status change not allowed")

const (
	codeLength          = 6
	maxCodeGenAttempts  = 10
	minCustomAliasLen   = 3
	maxCustomAliasLen   = 64
	defaultListPage     = 1
	defaultListPageSize = 20
	maxListPageSize     = 100
)

var customAliasPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
var allowedLinkStatuses = map[string]struct{}{
	StatusActive:   {},
	StatusInactive: {},
	StatusBlocked:  {},
	StatusDeleted:  {},
}
var allowedUserUpdateStatuses = map[string]struct{}{
	StatusActive:   {},
	StatusInactive: {},
}

type Service struct {
	repo     Repository
	generate func(length int) (string, error)
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:     repo,
		generate: generateCode,
	}
}

func (s *Service) Create(ctx context.Context, ownerID string, req CreateLinkRequest) (Link, map[string]string, error) {
	targetURL := strings.TrimSpace(req.TargetURL)
	fields := map[string]string{}
	customAlias := normalizeAlias(req.CustomAlias)
	if targetURL == "" {
		fields["targetUrl"] = "Target URL is required"
	} else if !isValidTargetURL(targetURL) {
		fields["targetUrl"] = "Target URL must be a valid absolute URL with http or https scheme"
	}
	if customAlias != nil {
		if len(*customAlias) < minCustomAliasLen || len(*customAlias) > maxCustomAliasLen {
			fields["customAlias"] = "Custom alias must be between 3 and 64 characters"
		} else if !customAliasPattern.MatchString(*customAlias) {
			fields["customAlias"] = "Custom alias can contain only letters, digits, _ and -"
		}
	}
	if ownerID == "" {
		fields["ownerId"] = "Owner ID is required"
	}
	if len(fields) > 0 {
		return Link{}, fields, ErrValidation
	}
	if customAlias != nil {
		exists, err := s.repo.ExistsByCustomAlias(ctx, *customAlias)
		if err != nil {
			return Link{}, nil, err
		}
		if exists {
			return Link{}, nil, ErrAliasAlreadyExists
		}
	}

	code, err := s.generateUniqueCode(ctx)
	if err != nil {
		return Link{}, nil, err
	}

	link := Link{
		OwnerID:     ownerID,
		Code:        code,
		CustomAlias: customAlias,
		TargetURL:   targetURL,
		Status:      StatusActive,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, &link); err != nil {
		if errors.Is(err, ErrAliasAlreadyExists) {
			return Link{}, nil, ErrAliasAlreadyExists
		}
		return Link{}, nil, err
	}

	return link, nil, nil
}

func (s *Service) List(ctx context.Context, ownerID string, query ListLinksQuery) ([]Link, PaginationResponse, map[string]string, error) {
	fields := map[string]string{}

	if strings.TrimSpace(ownerID) == "" {
		fields["ownerId"] = "Owner ID is required"
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
		fields["pageSize"] = "Page size must be less than or equal to 100"
	}

	status := strings.TrimSpace(query.Status)
	if status != "" {
		if _, ok := allowedLinkStatuses[status]; !ok {
			fields["status"] = "Status must be one of: active, inactive, blocked, deleted"
		}
	}

	if len(fields) > 0 {
		return nil, PaginationResponse{}, fields, ErrValidation
	}

	filter := ListLinksFilter{
		OwnerID:  ownerID,
		Page:     page,
		PageSize: pageSize,
		Q:        strings.TrimSpace(query.Q),
		Status:   status,
	}

	items, totalItems, err := s.repo.ListByOwner(ctx, filter)
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

func (s *Service) UpdateStatus(ctx context.Context, ownerID, linkID string, req UpdateLinkStatusRequest) (Link, map[string]string, error) {
	fields := map[string]string{}

	if strings.TrimSpace(ownerID) == "" {
		fields["ownerId"] = "Owner ID is required"
	}
	normalizedLinkID := strings.TrimSpace(linkID)
	if normalizedLinkID == "" {
		fields["linkId"] = "Link ID is required"
	} else if !shared.IsUUID(normalizedLinkID) {
		fields["linkId"] = "Link ID must be a valid UUID"
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		fields["status"] = "Status is required"
	} else if _, ok := allowedUserUpdateStatuses[status]; !ok {
		fields["status"] = "Status must be one of: active, inactive"
	}

	if len(fields) > 0 {
		return Link{}, fields, ErrValidation
	}

	link, err := s.repo.GetByIDAndOwner(ctx, normalizedLinkID, ownerID)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return Link{}, nil, ErrLinkNotFound
		}
		return Link{}, nil, err
	}

	if link.Status == StatusDeleted || link.Status == StatusBlocked {
		return Link{}, nil, ErrStatusChangeNotAllowed
	}

	updated, err := s.repo.UpdateStatus(ctx, normalizedLinkID, ownerID, status)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return Link{}, nil, ErrLinkNotFound
		}
		return Link{}, nil, err
	}

	return updated, nil, nil
}

func (s *Service) Delete(ctx context.Context, ownerID, linkID string) (map[string]string, error) {
	fields := map[string]string{}

	if strings.TrimSpace(ownerID) == "" {
		fields["ownerId"] = "Owner ID is required"
	}
	normalizedLinkID := strings.TrimSpace(linkID)
	if normalizedLinkID == "" {
		fields["linkId"] = "Link ID is required"
	} else if !shared.IsUUID(normalizedLinkID) {
		fields["linkId"] = "Link ID must be a valid UUID"
	}
	if len(fields) > 0 {
		return fields, ErrValidation
	}

	link, err := s.repo.GetByIDAndOwner(ctx, normalizedLinkID, ownerID)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}

	if link.Status == StatusDeleted || link.DeletedAt != nil {
		return nil, nil
	}

	if err := s.repo.SoftDelete(ctx, normalizedLinkID, ownerID); err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}

	return nil, nil
}

func generateCode(length int) (string, error) {
	raw := make([]byte, length)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	encoded := base64.RawURLEncoding.EncodeToString(raw)
	if len(encoded) < length {
		return "", errors.New("encoded code is shorter than expected")
	}

	return encoded[:length], nil
}

func (s *Service) generateUniqueCode(ctx context.Context) (string, error) {
	for attempt := 0; attempt < maxCodeGenAttempts; attempt++ {
		code, err := s.generate(codeLength)
		if err != nil {
			return "", fmt.Errorf("generate code: %w", err)
		}

		exists, err := s.repo.ExistsByCode(ctx, code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}

	return "", ErrCodeGenerationExhausted
}

func isValidTargetURL(raw string) bool {
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	if strings.TrimSpace(parsed.Host) == "" {
		return false
	}

	return true
}

func normalizeAlias(alias *string) *string {
	if alias == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*alias)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
