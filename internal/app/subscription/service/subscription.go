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

func New(repo domain.SubscriptionRepository) *subscription {
	return &subscription{repo: repo}
}

func (s *subscription) CreateSubscription(ctx context.Context, email string, bundleID int64) error {
	return s.repo.CreateSubscription(ctx, email, bundleID)
}

func (s *subscription) ReadSubscription(ctx context.Context, email string) (int64, error) {
	return s.repo.ReadSubscription(ctx, email)
}

func (s *subscription) UpdateSubscription(ctx context.Context, email string, bundleID int64) error {
	return s.repo.UpdateSubscription(ctx, email, bundleID)
}

func (s *subscription) CancelSubscription(ctx context.Context, email string) error {
	return s.repo.DeleteSubscription(ctx, email)
}

func (s *subscription) AddBalance(ctx context.Context, email string, amount uint64) error {
	return s.repo.AddBalance(ctx, email, amount)
}

func (s *subscription) ReadUserData(ctx context.Context, email string) (domain.UserData, error) {
	return s.repo.ReadUserData(ctx, email)
}

func (s *subscription) ReadBalanceHistory(ctx context.Context, email string, page uint64) ([]*api.HistoryElementV1, error) {
	return s.repo.ReadBalanceHistory(ctx, email, page)
}
