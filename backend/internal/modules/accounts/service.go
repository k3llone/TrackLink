package accounts

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength = 8
)

var ErrValidation = errors.New("validation failed")
var ErrConflict = errors.New("resource conflict")
var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (User, map[string]string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password

	fields := validateRegisterInput(email, strings.TrimSpace(password))
	if len(fields) > 0 {
		return User{}, fields, ErrValidation
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return User{}, nil, fmt.Errorf("hash password: %w", err)
	}

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         RoleCustomer,
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return User{}, nil, ErrConflict
		}
		return User{}, nil, err
	}

	return user, nil, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (User, map[string]string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password

	fields := validateLoginInput(email, strings.TrimSpace(password))
	if len(fields) > 0 {
		return User{}, fields, ErrValidation
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return User{}, nil, ErrInvalidCredentials
		}
		return User{}, nil, err
	}

	if err := comparePassword(user.PasswordHash, password); err != nil {
		return User{}, nil, ErrInvalidCredentials
	}

	return user, nil, nil
}

func (s *Service) CurrentUser(ctx context.Context, userID string) (User, error) {
	id := strings.TrimSpace(userID)
	if id == "" {
		return User{}, ErrUserNotFound
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func validateRegisterInput(email, password string) map[string]string {
	fields := make(map[string]string)

	if email == "" {
		fields["email"] = "Email is required"
	} else if _, err := mail.ParseAddress(email); err != nil {
		fields["email"] = "Email is invalid"
	}

	if password == "" {
		fields["password"] = "Password is required"
	} else if len(password) < MinPasswordLength {
		fields["password"] = fmt.Sprintf("Password must be at least %d characters long", MinPasswordLength)
	}

	if len(fields) == 0 {
		return nil
	}
	return fields
}

func validateLoginInput(email, password string) map[string]string {
	fields := make(map[string]string)

	if email == "" {
		fields["email"] = "Email is required"
	} else if _, err := mail.ParseAddress(email); err != nil {
		fields["email"] = "Email is invalid"
	}

	if password == "" {
		fields["password"] = "Password is required"
	}

	if len(fields) == 0 {
		return nil
	}
	return fields
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
