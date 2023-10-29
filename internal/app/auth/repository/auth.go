package repository

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type auth struct {
	mongo *mCollection
}

func New(mongoCollection *mongo.Collection) *auth {
	return &auth{mongo: &mCollection{mongoCollection, &sync.RWMutex{}}}
}

func (r *auth) SaveUserData(ctx context.Context, email string, passwordHash string, address string, secretKey string, isAdmin bool) error {
	return r.mongo.InsertOne(ctx, bson.M{"email": email, "passwordHash": passwordHash, "address": address, "isAdmin": isAdmin, "secretKey": secretKey})
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
