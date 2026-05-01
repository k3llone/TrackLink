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
)

var ErrValidation = errors.New("validation failed")
var ErrCodeGenerationExhausted = errors.New("code generation attempts exhausted")

const (
	codeLength          = 6
	maxCodeGenAttempts  = 10
	minCustomAliasLen   = 3
	maxCustomAliasLen   = 64
)

var customAliasPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type Service struct {
	repo      Repository
	generate  func(length int) (string, error)
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
