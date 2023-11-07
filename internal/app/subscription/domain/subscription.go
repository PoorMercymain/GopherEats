package domain

import (
	"context"

	"github.com/PoorMercymain/GopherEats/pkg/api/subscription"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, email string, bundleID int64) error
	ReadSubscription(ctx context.Context, email string) (int64, bool, error)
	UpdateSubscription(ctx context.Context, email string, bundleID int64, isDeleted bool) error
	CancelSubscription(ctx context.Context, email string) error
	AddBalance(ctx context.Context, email string, balance uint64) error
	ReadUserData(ctx context.Context, email string) (UserData, error)
	ReadBalanceHistory(ctx context.Context, email string, page uint64) ([]*subscription.HistoryElementV1, error)
	ChargeForSubscription(ctx context.Context) error
}

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, email string, bundleID int64) error
	ReadSubscription(ctx context.Context, email string) (int64, bool, error)
	UpdateSubscription(ctx context.Context, email string, bundleID int64, isDeleted bool) error
	DeleteSubscription(ctx context.Context, email string) error
	AddBalance(ctx context.Context, email string, balance uint64) error
	ReadUserData(ctx context.Context, email string) (UserData, error)
	ReadBalanceHistory(ctx context.Context, email string, page uint64) ([]*subscription.HistoryElementV1, error)
	ChargeForSubscription(ctx context.Context) error
	RemoveAmountFromBalance(ctx context.Context, email string, amount uint64) error
}
