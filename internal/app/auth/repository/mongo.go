package repository

import (
	"context"
	"errors"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"

	authErrors "github.com/PoorMercymain/GopherEats/internal/app/auth/errors"
)

type mCollection struct {
	mongoCollection *mongo.Collection
	*sync.RWMutex
}

func (mc *mCollection) InsertOne(ctx context.Context, document interface{}) error {
	mc.Lock()
	_, err := mc.mongoCollection.InsertOne(ctx, document)
	mc.Unlock()

	if mongo.IsDuplicateKeyError(err) {
		return authErrors.ErrorUserAlreadyExists
	}

	return err
}

func (mc *mCollection) FindOne(ctx context.Context, filter interface{}) (*mongo.SingleResult, error) {
	mc.RLock()
	findHashResult := mc.mongoCollection.FindOne(ctx, filter)
	mc.RUnlock()

	if errors.Is(findHashResult.Err(), mongo.ErrNoDocuments) {
		return nil, authErrors.ErrorNoSuchUser
	}

	if findHashResult.Err() != nil {
		return nil, findHashResult.Err()
	}

	return findHashResult, nil
}
