// Package service contains business logic for subscription service.
package service

import (
	"context"

	"github.com/PoorMercymain/GopherEats/internal/app/subscription/domain"
	api "github.com/PoorMercymain/GopherEats/pkg/api/subscription"
)

var _ domain.SubscriptionService = (*subscription)(nil)

type subscription struct {
	repo domain.SubscriptionRepository
}

// New returns pointer to new instance of subscription struct with repo.
func New(repo domain.SubscriptionRepository) *subscription {
	return &subscription{repo: repo}
}

// CreateSubscription creates new subscription.
func (s *subscription) CreateSubscription(ctx context.Context, email string, bundleID int64) error {
	return s.repo.CreateSubscription(ctx, email, bundleID)
}

// ReadSubscription returns subscription info.
func (s *subscription) ReadSubscription(ctx context.Context, email string) (int64, bool, error) {
	return s.repo.ReadSubscription(ctx, email)
}

// UpdateSubscription updates existing subscription.
func (s *subscription) UpdateSubscription(ctx context.Context, email string, bundleID int64, isDeleted bool) error {
	return s.repo.UpdateSubscription(ctx, email, bundleID, isDeleted)
}

// CancelSubscription deletes subscription.
func (s *subscription) CancelSubscription(ctx context.Context, email string) error {
	return s.repo.DeleteSubscription(ctx, email)
}

// AddBalance adds funds to user balance.
func (s *subscription) AddBalance(ctx context.Context, email string, amount uint64) error {
	return s.repo.AddBalance(ctx, email, amount)
}

// ReadUserData returns user data.
func (s *subscription) ReadUserData(ctx context.Context, email string) (domain.UserData, error) {
	return s.repo.ReadUserData(ctx, email)
}

// ReadBalanceHistory returns history of balance funding and charging.
func (s *subscription) ReadBalanceHistory(ctx context.Context, email string, page uint64) ([]*api.HistoryElementV1, error) {
	return s.repo.ReadBalanceHistory(ctx, email, page)
}

// ChargeForSubscription charges user for weekly subscription cost.
func (s *subscription) ChargeForSubscription(ctx context.Context, notEnoughFundsEmailsChan chan<- string) error { // TODO: use it every thursday
	return s.repo.ChargeForSubscription(ctx, notEnoughFundsEmailsChan)
}
