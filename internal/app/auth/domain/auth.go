// Package domain contains interfaces for auth service.
package domain

import "context"

type AuthService interface {
	RegisterUser(ctx context.Context, email string, password string, address string, secretKey string, isAdmin bool) error
	CheckAuthData(ctx context.Context, email string, password string) error
	GetPasswordHash(ctx context.Context, email string) (string, error)
	UpdatePassword(ctx context.Context, email string, oldPassword string, newPassword string) error
	LoginWithOTP(ctx context.Context, email string, otpCode string) error
	UpdateAddress(ctx context.Context, email string, password string, newAddress string) error
	CheckIfUserIsAdmin(ctx context.Context, email string) (bool, error)
	GetAddress(ctx context.Context, email string) (string, error)
	UpdateWarning(ctx context.Context, email string, warning string) error
	GetWarning(ctx context.Context, email string) (string, error)
}

type AuthRepository interface {
	SaveUserData(ctx context.Context, email string, passwordHash string, address string, secretKey string, isAdmin bool) error
	GetPasswordHash(ctx context.Context, email string) (string, error)
	UpdatePassword(ctx context.Context, email string, newPasswordHash string) error
	GetSecretKey(ctx context.Context, email string) (string, error)
	UpdateAddress(ctx context.Context, email string, newAddress string) error
	GetIsAdmin(ctx context.Context, email string) (bool, error)
	GetAddress(ctx context.Context, email string) (string, error)
	UpdateWarning(ctx context.Context, email string, warning string) error
	GetWarning(ctx context.Context, email string) (string, error)
}
