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

func New(mongoCollection *mongo.Collection) *auth {
	return &auth{mongo: &mCollection{mongoCollection, &sync.RWMutex{}}}
}

func (r *auth) SaveUserData(ctx context.Context, email string, passwordHash string, address string, secretKey string, isAdmin bool) error {
	return r.mongo.InsertOne(ctx, bson.M{"email": email, "passwordHash": passwordHash, "address": address, "isAdmin": isAdmin, "secretKey": secretKey, "warning": ""})
}

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

func (r *auth) UpdatePassword(ctx context.Context, email string, newPasswordHash string) error {
	_, err := r.mongo.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"passwordHash": newPasswordHash}})

	return err
}

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

func (r *auth) UpdateAddress(ctx context.Context, email string, newAddress string) error {
	_, err := r.mongo.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"address": newAddress}})

	return err
}

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

func (r *auth) UpdateWarning(ctx context.Context, email string, warning string) error {
	_, err := r.mongo.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"warning": warning}})

	return err
}

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
