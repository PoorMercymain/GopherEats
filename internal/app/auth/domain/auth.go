package domain

import "context"

type AuthService interface {
	RegisterUser(ctx context.Context, email string, password string) error
}

type AuthRepository interface {
	SaveUserAuthData(ctx context.Context, email string, passwordHash string) error
}
