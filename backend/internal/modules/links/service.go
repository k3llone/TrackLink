package links

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

var ErrValidation = errors.New("validation failed")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, ownerID string, req CreateLinkRequest) (Link, map[string]string, error) {
	targetURL := strings.TrimSpace(req.TargetURL)
	fields := map[string]string{}
	if targetURL == "" {
		fields["targetUrl"] = "Target URL is required"
	} else if !isValidTargetURL(targetURL) {
		fields["targetUrl"] = "Target URL must be a valid absolute URL with http or https scheme"
	}
	if ownerID == "" {
		fields["ownerId"] = "Owner ID is required"
	}
	if len(fields) > 0 {
		return Link{}, fields, ErrValidation
	}

	code, err := generateCode(6)
	if err != nil {
		return Link{}, nil, fmt.Errorf("generate code: %w", err)
	}

	link := Link{
		OwnerID:     ownerID,
		Code:        code,
		CustomAlias: req.CustomAlias,
		TargetURL:   targetURL,
		Status:      StatusActive,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, &link); err != nil {
		return Link{}, nil, err
	}

	return link, nil, nil
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
