package domain

import "context"

type AuthService interface {
	RegisterUser(ctx context.Context, email string, password string) error
	CheckAuthData(ctx context.Context, email string, password string) error
	GetPasswordHash(ctx context.Context, email string) (string, error)
	UpdatePassword(ctx context.Context, email string, oldPassword string, newPassword string) error
}

type AuthRepository interface {
	SaveUserAuthData(ctx context.Context, email string, passwordHash string) error
	GetPasswordHash(ctx context.Context, email string) (string, error)
	UpdatePassword(ctx context.Context, email string, newPasswordHash string) error
}
