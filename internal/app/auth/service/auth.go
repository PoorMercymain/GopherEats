package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/pquerna/otp/totp"

	"github.com/PoorMercymain/GopherEats/internal/app/auth/domain"
	"github.com/PoorMercymain/GopherEats/internal/app/auth/errors"
)

var _ domain.AuthService = (*auth)(nil)

type auth struct {
	repo domain.AuthRepository
}

func New(repo domain.AuthRepository) *auth {
	return &auth{repo: repo}
}

func (s *auth) RegisterUser(ctx context.Context, email string, password string, address string, secretKey string, isAdmin bool) error {
	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	strHash := string(byteHash)

	return s.repo.SaveUserData(ctx, email, strHash, address, secretKey, isAdmin)
}

func (s *auth) CheckAuthData(ctx context.Context, email string, password string) error {
	hash, err := s.repo.GetPasswordHash(ctx, email)
	if err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return errors.ErrorWrongPassword
	}

	return nil
}

func (s *auth) GetPasswordHash(ctx context.Context, email string) (string, error) {
	return s.repo.GetPasswordHash(ctx, email)
}

func (s *auth) UpdatePassword(ctx context.Context, email string, oldPassword string, newPassword string) error {
	hash, err := s.repo.GetPasswordHash(ctx, email)
	if err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(oldPassword)) != nil {
		return errors.ErrorWrongPassword
	}

	byteHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	strHash := string(byteHash)

	err = s.repo.UpdatePassword(ctx, email, strHash)
	if err != nil {
		return err
	}

	return nil
}

func (s *auth) LoginWithOTP(ctx context.Context, email string, otpCode string) error {
	secretKey, err := s.repo.GetSecretKey(ctx, email)
	if err != nil {
		return err
	}

	if !totp.Validate(otpCode, secretKey) {
		return errors.ErrorWrongOTP
	}

	return nil
}

func (s *auth) UpdateAddress(ctx context.Context, email string, password string, newAddress string) error {
	hash, err := s.repo.GetPasswordHash(ctx, email)
	if err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return errors.ErrorWrongPassword
	}

	err = s.repo.UpdateAddress(ctx, email, newAddress)
	if err != nil {
		return err
	}

	return nil
}

func (s *auth) CheckIfUserIsAdmin(ctx context.Context, email string) (bool, error) {
	return s.repo.GetIsAdmin(ctx, email)
}

func (s *auth) GetAddress(ctx context.Context, email string) (string, error) {
	return s.repo.GetAddress(ctx, email)
}

func (s *auth) UpdateWarning(ctx context.Context, email string, warning string) error {
	return s.repo.UpdateWarning(ctx, email, warning)
}

func (s *auth) GetWarning(ctx context.Context, email string) (string, error) {
	return s.repo.GetWarning(ctx, email)
}
