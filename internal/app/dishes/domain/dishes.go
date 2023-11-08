package domain

import (
	"context"
	pb "github.com/PoorMercymain/GopherEats/pkg/api/dishes"
)

type DishesService interface {
	StoreIngredient(ctx context.Context, ingredient *pb.IngredientV1) error
	UpdateIngredient(ctx context.Context, ingredient *pb.IngredientV1) error
	GetIngredient(ctx context.Context, id uint64) (*pb.IngredientV1, error)
	DeleteIngredient(ctx context.Context, id uint64) error
	ListIngredients(ctx context.Context) ([]*pb.IngredientV1, error)

	StoreDish(ctx context.Context, dish *pb.DishV1) error
	UpdateDish(ctx context.Context, dish *pb.DishV1) error
	GetDish(ctx context.Context, id uint64) (*pb.DishV1, error)
	DeleteDish(ctx context.Context, id uint64) error
	ListDishes(ctx context.Context) ([]*pb.DishV1, error)

	StoreBundle(ctx context.Context, bundle *pb.BundleV1) error
	UpdateBundle(ctx context.Context, bundle *pb.BundleV1) error
	GetBundle(ctx context.Context, id uint64) (*pb.BundleV1, error)
	DeleteBundle(ctx context.Context, id uint64) error
	ListBundles(ctx context.Context) ([]*pb.BundleV1, error)

	AddBundleWeeklyDish(ctx context.Context, weekNumber, bundleId, dishId uint64) error
	DeleteBundleWeeklyDish(ctx context.Context, weekNumber, bundleId, dishId uint64) error
	GetBundleWeeklyDishes(ctx context.Context, weekNumber, bundleId uint64) (*pb.DishV1, error)
	/*
		StoreResource(ctx context.Context, resource *pb.ResourceV1) error
		UpdateResource(ctx context.Context, resource *pb.ResourceV1) error
		GetResource(ctx context.Context, id uint64) (*pb.ResourceV1, error)
		DeleteResource(ctx context.Context, id uint64) error
		ListResources(ctx context.Context) ([]*pb.ResourceV1, error)

		AttachResourceToDish(ctx context.Context, resourceId, dishId uint64) error
		DetachResourceFromDish(ctx context.Context, resourceId, dishId uint64) error
		ListDishResources(ctx context.Context, dishId uint64) ([]*pb.ResourceV1, error)
	*/
	SendOrder(ctx context.Context, order *pb.SendOrderRequestV1) error
	CancelOrder(ctx context.Context, order *pb.CancelOrderRequestV1) error
}

type DishesRepository interface {
	StoreIngredient(ctx context.Context, ingredient *Ingredient) error
	UpdateIngredient(ctx context.Context, ingredient *Ingredient) error
	GetIngredient(ctx context.Context, id uint64) (*Ingredient, error)
	DeleteIngredient(ctx context.Context, id uint64) error
	ListIngredients(ctx context.Context) ([]*Ingredient, error)

	StoreDish(ctx context.Context, dish *Dish) error
	UpdateDish(ctx context.Context, dish *Dish) error
	GetDish(ctx context.Context, id uint64) (*Dish, error)
	DeleteDish(ctx context.Context, id uint64) error
	ListDishes(ctx context.Context) ([]*Dish, error)

	StoreBundle(ctx context.Context, bundle *Bundle) error
	UpdateBundle(ctx context.Context, bundle *Bundle) error
	GetBundle(ctx context.Context, id uint64) (*Bundle, error)
	DeleteBundle(ctx context.Context, id uint64) error
	ListBundles(ctx context.Context) ([]*Bundle, error)

	AddBundleWeeklyDish(ctx context.Context, weekNumber, bundleId, dishId uint64) error
	DeleteBundleWeeklyDish(ctx context.Context, weekNumber, bundleId, dishId uint64) error
	GetBundleWeeklyDishes(ctx context.Context, weekNumber, bundleId uint64) ([]*Dish, error)
	/*
		StoreResource(ctx context.Context, resource *Resource) error
		UpdateResource(ctx context.Context, resource *Resource) error
		GetResource(ctx context.Context, id uint64) (*Resource, error)
		DeleteResource(ctx context.Context, id uint64) error
		ListResources(ctx context.Context) ([]*Resource, error)

		AttachResourceToDish(ctx context.Context, resourceId, dishId uint64) error
		DetachResourceFromDish(ctx context.Context, resourceId, dishId uint64) error
		ListDishResources(ctx context.Context, dishId uint64) ([]*Resource, error)
	*/
	Ping(ctx context.Context) error
	ClosePool() error
}
