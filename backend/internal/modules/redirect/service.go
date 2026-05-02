package redirect

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Resolve(ctx context.Context, code string) (Link, error) {
	normalizedCode := strings.TrimSpace(code)
	if normalizedCode == "" {
		return Link{}, ErrLinkNotFound
	}

	link, err := s.repo.FindByCodeOrAlias(ctx, normalizedCode)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return Link{}, ErrLinkNotFound
		}
		return Link{}, err
	}

	return link, nil
}
