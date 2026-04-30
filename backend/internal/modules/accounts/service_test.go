package accounts

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type fakeRepo struct {
	createFn      func(ctx context.Context, user *User) error
	findByEmailFn func(ctx context.Context, email string) (User, error)
}

func (f fakeRepo) Create(ctx context.Context, user *User) error {
	if f.createFn == nil {
		return nil
	}
	return f.createFn(ctx, user)
}

func (f fakeRepo) FindByEmail(ctx context.Context, email string) (User, error) {
	if f.findByEmailFn == nil {
		return User{}, ErrUserNotFound
	}
	return f.findByEmailFn(ctx, email)
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

func TestServiceRegisterDuplicateEmail(t *testing.T) {
	repo := fakeRepo{
		createFn: func(_ context.Context, _ *User) error {
			return ErrEmailAlreadyExists
		},
	}

	service := NewService(repo)
	_, fields, err := service.Register(context.Background(), RegisterRequest{
		Email:    "new@example.com",
		Password: "StrongPass123",
	})

	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil validation fields, got %v", fields)
	}
}

func TestServiceLoginSuccess(t *testing.T) {
	passwordHash, err := hashPassword("StrongPass123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := fakeRepo{
		findByEmailFn: func(_ context.Context, email string) (User, error) {
			if email != "user@example.com" {
				t.Fatalf("expected normalized email, got %s", email)
			}
			return User{
				ID:           "f98b832d-6f5b-4bcf-9175-8e56f5e983f0",
				Email:        "user@example.com",
				PasswordHash: passwordHash,
				Role:         RoleCustomer,
				CreatedAt:    time.Date(2026, 4, 30, 19, 0, 0, 0, time.UTC),
			}, nil
		},
	}

	service := NewService(repo)
	user, fields, err := service.Login(context.Background(), LoginRequest{
		Email:    " User@Example.com ",
		Password: "StrongPass123",
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if user.Email != "user@example.com" {
		t.Fatalf("unexpected user email: %s", user.Email)
	}
}

func TestServiceLoginInvalidCredentialsForUnknownEmail(t *testing.T) {
	repo := fakeRepo{
		findByEmailFn: func(_ context.Context, _ string) (User, error) {
			return User{}, ErrUserNotFound
		},
	}

	service := NewService(repo)
	_, fields, err := service.Login(context.Background(), LoginRequest{
		Email:    "unknown@example.com",
		Password: "StrongPass123",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
}

func TestServiceLoginInvalidCredentialsForWrongPassword(t *testing.T) {
	passwordHash, err := hashPassword("StrongPass123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := fakeRepo{
		findByEmailFn: func(_ context.Context, _ string) (User, error) {
			return User{
				Email:        "user@example.com",
				PasswordHash: passwordHash,
			}, nil
		},
	}

	service := NewService(repo)
	_, fields, err := service.Login(context.Background(), LoginRequest{
		Email:    "user@example.com",
		Password: "WrongPass123",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
}

func TestServiceLoginValidation(t *testing.T) {
	repo := fakeRepo{
		findByEmailFn: func(_ context.Context, _ string) (User, error) {
			t.Fatal("repository must not be called on validation error")
			return User{}, nil
		},
	}

	service := NewService(repo)
	_, fields, err := service.Login(context.Background(), LoginRequest{
		Email:    "bad-email",
		Password: "",
	})

	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
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
