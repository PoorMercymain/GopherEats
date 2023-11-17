package mockgen

import (
	"context"

	"github.com/PoorMercymain/GopherEats/internal/app/dishes/domain"
)

//go:generate mockgen -destination=../repository/mocks/storage_mock.gen.go -package=mocks . StorageGen
type StorageGen interface {
	StoreIngredient(ctx context.Context, ingredient *domain.Ingredient) error
	UpdateIngredient(ctx context.Context, ingredient *domain.Ingredient) error
	GetIngredient(ctx context.Context, id uint64) error
	DeleteIngredient(ctx context.Context, id uint64) error
	ListIngredients(ctx context.Context) ([]*domain.Ingredient, error)

	StoreDish(ctx context.Context, dish *domain.Dish) error
	UpdateDish(ctx context.Context, dish *domain.Dish) error
	GetDish(ctx context.Context, id uint64) error
	DeleteDish(ctx context.Context, id uint64) error
	ListDishes(ctx context.Context) ([]*domain.Dish, error)

	StoreBundle(ctx context.Context, bundle *domain.Bundle) error
	UpdateBundle(ctx context.Context, bundle *domain.Bundle) error
	GetBundle(ctx context.Context, id uint64) error
	DeleteBundle(ctx context.Context, id uint64) error
	ListBundles(ctx context.Context) ([]*domain.Bundle, error)

	AddBundleWeeklyDish(ctx context.Context, weekNumber, bundleID, dishID uint64) error
	DeleteBundleWeeklyDish(ctx context.Context, weekNumber, bundleID, dishID uint64) error
	GetBundleWeeklyDishes(ctx context.Context, weekNumber, bundleID uint64) ([]*domain.Dish, error)
	/*
		StoreResource(ctx context.Context, resource *domain.Resource) error
		UpdateResource(ctx context.Context, resource *domain.Resource) error
		GetResource(ctx context.Context, id uint64) error
		DeleteResource(ctx context.Context, id uint64) error
		ListResources(ctx context.Context) ([]*domain.Resource, error)

		AttachResourceToDish(ctx context.Context, resourceId, dishID uint64) error
		DetachResourceFromDish(ctx context.Context, resourceId, dishID uint64) error
		ListDishResources(ctx context.Context, dishID uint64) ([]*domain.Resource, error)
	*/
	Ping(ctx context.Context) error
	ClosePool() error
}
