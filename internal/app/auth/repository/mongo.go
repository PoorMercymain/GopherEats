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

// InsertOne executes an insert command to insert a single document into the collection.
func (mc *mCollection) InsertOne(ctx context.Context, document interface{}) error {
	mc.Lock()
	_, err := mc.mongoCollection.InsertOne(ctx, document)
	mc.Unlock()

	if mongo.IsDuplicateKeyError(err) {
		return authErrors.ErrorUserAlreadyExists
	}

	return err
}

// FindOne executes a find command and returns a SingleResult for one document in the collection.
func (mc *mCollection) FindOne(ctx context.Context, filter interface{}) (*mongo.SingleResult, error) {
	mc.RLock()
	findResult := mc.mongoCollection.FindOne(ctx, filter)
	mc.RUnlock()

	if errors.Is(findResult.Err(), mongo.ErrNoDocuments) {
		return nil, authErrors.ErrorNoSuchUser
	}

	if findResult.Err() != nil {
		return nil, findResult.Err()
	}

	return findResult, nil
}

// UpdateOne executes an update command to update at most one document in the collection.
func (mc *mCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	mc.Lock()
	updateResult, err := mc.mongoCollection.UpdateOne(ctx, filter, update)
	mc.Unlock()

	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, authErrors.ErrorUserWasNotUpdated
	}

	return updateResult, nil
}
