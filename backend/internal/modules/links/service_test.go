package links

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	createFn func(ctx context.Context, link *Link) error
}

func (f fakeRepository) Create(ctx context.Context, link *Link) error {
	if f.createFn == nil {
		return nil
	}

	return f.createFn(ctx, link)
}

func TestServiceCreateSuccess(t *testing.T) {
	repo := fakeRepository{
		createFn: func(_ context.Context, link *Link) error {
			link.ID = "9e1f66cf-4f5c-421f-97fa-bf4a9f34ef7a"
			return nil
		},
	}

	service := NewService(repo)
	link, fields, err := service.Create(context.Background(), "c7364dce-f6fd-4ec4-b7bc-f2a95f2f9de8", CreateLinkRequest{
		TargetURL: "https://example.com/landing",
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if link.ID == "" {
		t.Fatal("expected link ID to be set")
	}
	if link.Status != StatusActive {
		t.Fatalf("expected status %q, got %q", StatusActive, link.Status)
	}
	if link.Code == "" {
		t.Fatal("expected code to be generated")
	}
}

func TestServiceCreateValidation(t *testing.T) {
	repo := fakeRepository{
		createFn: func(_ context.Context, _ *Link) error {
			t.Fatal("repo must not be called on validation errors")
			return nil
		},
	}

	service := NewService(repo)
	_, fields, err := service.Create(context.Background(), "", CreateLinkRequest{
		TargetURL: " ",
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
	if fields == nil {
		t.Fatal("expected validation fields")
	}
	if _, ok := fields["targetUrl"]; !ok {
		t.Fatal("expected targetUrl field validation error")
	}
	if _, ok := fields["ownerId"]; !ok {
		t.Fatal("expected ownerId field validation error")
	}
}

func TestServiceCreateTargetURLValidation(t *testing.T) {
	testCases := []struct {
		name      string
		targetURL string
		wantError bool
	}{
		{name: "relative URL", targetURL: "/landing", wantError: true},
		{name: "missing host", targetURL: "https://", wantError: true},
		{name: "unsupported scheme", targetURL: "ftp://example.com", wantError: true},
		{name: "valid https URL", targetURL: "https://example.com/path", wantError: false},
	}

	repo := fakeRepository{
		createFn: func(_ context.Context, _ *Link) error {
			return nil
		},
	}
	service := NewService(repo)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, fields, err := service.Create(context.Background(), "c7364dce-f6fd-4ec4-b7bc-f2a95f2f9de8", CreateLinkRequest{
				TargetURL: tc.targetURL,
			})

			if tc.wantError {
				if !errors.Is(err, ErrValidation) {
					t.Fatalf("expected ErrValidation, got %v", err)
				}
				if _, ok := fields["targetUrl"]; !ok {
					t.Fatalf("expected targetUrl field error, got %v", fields)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if fields != nil {
				t.Fatalf("expected nil fields, got %v", fields)
			}
		})
	}
}
