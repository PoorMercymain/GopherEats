package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/PoorMercymain/GopherEats/internal/app/auth/domain"
	"github.com/PoorMercymain/GopherEats/internal/app/auth/errors"
)

type auth struct {
	repo domain.AuthRepository
}

func New(repo domain.AuthRepository) *auth {
	return &auth{repo: repo}
}

func (s *auth) RegisterUser(ctx context.Context, email string, password string) error {
	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	strHash := string(byteHash)

	return s.repo.SaveUserAuthData(ctx, email, strHash)
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
