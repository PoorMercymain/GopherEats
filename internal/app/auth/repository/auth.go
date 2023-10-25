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

func (r *auth) SaveUserAuthData(ctx context.Context, email string, passwordHash string) error {
	return r.mongo.InsertOne(ctx, bson.M{"email": email, "passwordHash": passwordHash})
}

func (r *auth) GetPasswordHash(ctx context.Context, email string) (string, error) {
	hash, err := r.mongo.FindOne(ctx, bson.M{"email": email})
	if err != nil {
		return "", err
	}

	var user struct {
		Email        string `bson:"email"`
		PasswordHash string `bson:"passwordHash"`
	}

	err = hash.Decode(&user)
	if err != nil {
		return "", err
	}

	return user.PasswordHash, nil
}
