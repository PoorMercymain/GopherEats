package repository

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/PoorMercymain/GopherEats/internal/app/auth/domain"
)

var _ domain.AuthRepository = (*auth)(nil)

type auth struct {
	mongo *mCollection
}

// New function returns pointer to new instance of auth mongo repository struct.
func New(mongoCollection *mongo.Collection) *auth {
	return &auth{mongo: &mCollection{mongoCollection, &sync.RWMutex{}}}
}

// SaveUserData saves user data to repo.
func (r *auth) SaveUserData(ctx context.Context, email string, passwordHash string, address string, secretKey string, isAdmin bool) error {
	return r.mongo.InsertOne(ctx, bson.M{"email": email, "passwordHash": passwordHash, "address": address, "isAdmin": isAdmin, "secretKey": secretKey, "warning": ""})
}

// GetPasswordHash returns password hash.
func (r *auth) GetPasswordHash(ctx context.Context, email string) (string, error) {
	user, err := r.mongo.FindOne(ctx, bson.M{"email": email})
	if err != nil {
		return "", err
	}

	var authData struct {
		Email        string `bson:"email"`
		PasswordHash string `bson:"passwordHash"`
	}

	err = user.Decode(&authData)
	if err != nil {
		return "", err
	}

	return authData.PasswordHash, nil
}

// UpdatePassword allows to update user password.
func (r *auth) UpdatePassword(ctx context.Context, email string, newPasswordHash string) error {
	_, err := r.mongo.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"passwordHash": newPasswordHash}})

	return err
}

// GetSecretKey returns secret key.
func (r *auth) GetSecretKey(ctx context.Context, email string) (string, error) {
	user, err := r.mongo.FindOne(ctx, bson.M{"email": email})
	if err != nil {
		return "", err
	}

	var key struct {
		SecretKey string `bson:"secretKey"`
	}

	err = user.Decode(&key)
	if err != nil {
		return "", err
	}

	return key.SecretKey, nil
}

// UpdateAddress allows to update client address.
func (r *auth) UpdateAddress(ctx context.Context, email string, newAddress string) error {
	_, err := r.mongo.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"address": newAddress}})

	return err
}

// GetIsAdmin returns whether user has administrator role.
func (r *auth) GetIsAdmin(ctx context.Context, email string) (bool, error) {
	user, err := r.mongo.FindOne(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	var isAdmin struct {
		IsAdmin bool `bson:"isAdmin"`
	}

	err = user.Decode(&isAdmin)
	if err != nil {
		return false, err
	}

	return isAdmin.IsAdmin, nil
}

// GetAddress returns user address.
func (r *auth) GetAddress(ctx context.Context, email string) (string, error) {
	user, err := r.mongo.FindOne(ctx, bson.M{"email": email})
	if err != nil {
		return "", err
	}

	var address struct {
		Address string `bson:"address"`
	}

	err = user.Decode(&address)
	if err != nil {
		return "", err
	}

	return address.Address, nil
}

// UpdateWarning updates warning.
func (r *auth) UpdateWarning(ctx context.Context, email string, warning string) error {
	_, err := r.mongo.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"warning": warning}})

	return err
}

// GetWarning returns warning for user with specified email.
func (r *auth) GetWarning(ctx context.Context, email string) (string, error) {
	user, err := r.mongo.FindOne(ctx, bson.M{"email": email})
	if err != nil {
		return "", err
	}

	var warning struct {
		Warning string `bson:"warning"`
	}

	err = user.Decode(&warning)
	if err != nil {
		return "", err
	}

	return warning.Warning, nil
}
