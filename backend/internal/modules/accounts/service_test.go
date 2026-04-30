package accounts

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type fakeRepo struct {
	createFn func(ctx context.Context, user *User) error
}

func (f fakeRepo) Create(ctx context.Context, user *User) error {
	return f.createFn(ctx, user)
}

func TestServiceRegisterSuccess(t *testing.T) {
	repo := fakeRepo{
		createFn: func(_ context.Context, user *User) error {
			user.ID = "6d5ac8d0-34c5-40eb-a1c7-b0ecf2f95819"
			user.CreatedAt = time.Date(2026, 4, 30, 15, 0, 0, 0, time.UTC)
			return nil
		},
	}

	service := NewService(repo)
	user, fields, err := service.Register(context.Background(), RegisterRequest{
		Email:    "  User@Example.COM ",
		Password: "StrongPass123",
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil validation fields, got %v", fields)
	}
	if user.Email != "user@example.com" {
		t.Fatalf("email normalization failed: %s", user.Email)
	}
	if user.Role != RoleCustomer {
		t.Fatalf("expected role %q, got %q", RoleCustomer, user.Role)
	}
	if user.PasswordHash == "" {
		t.Fatal("password hash must not be empty")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("StrongPass123")) != nil {
		t.Fatal("password hash does not match the source password")
	}
}

func TestServiceRegisterValidation(t *testing.T) {
	repo := fakeRepo{
		createFn: func(_ context.Context, _ *User) error {
			t.Fatal("repository must not be called on validation error")
			return nil
		},
	}

	service := NewService(repo)
	_, fields, err := service.Register(context.Background(), RegisterRequest{
		Email:    "wrong-email",
		Password: "123",
	})

	if err == nil {
		t.Fatal("expected validation error")
	}
	if fields == nil {
		t.Fatal("expected validation fields")
	}
	if _, ok := fields["email"]; !ok {
		t.Fatal("expected email validation error")
	}
	if _, ok := fields["password"]; !ok {
		t.Fatal("expected password validation error")
	}
}
