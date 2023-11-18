// Package service contains business logic for dishes service.
package service

import (
	"context"
	"fmt"

	"github.com/PoorMercymain/GopherEats/internal/pkg/logger"

	"github.com/PoorMercymain/GopherEats/internal/app/dishes/domain"
	pb "github.com/PoorMercymain/GopherEats/pkg/api/dishes"
)

type dishesService struct {
	repo domain.DishesRepository
}

// NewDishesServiceWithRepository returns pointer to new DishesService with passed repo.
func NewDishesServiceWithRepository(repo domain.DishesRepository) (service domain.DishesService, err error) {
	service = &dishesService{repo: repo}
	return
}

// StoreIngredient saves new Ingredient.
func (d *dishesService) StoreIngredient(ctx context.Context, req *pb.CreateIngredientRequestV1) (err error) {
	i := &domain.Ingredient{
		Name: req.GetName(),
		Unit: req.GetUnit(),
	}

	err = d.repo.StoreIngredient(ctx, i)
	return
}

// UpdateIngredient updates existing Ingredient.
func (d *dishesService) UpdateIngredient(ctx context.Context, req *pb.UpdateIngredientRequestV1) (err error) {
	i := &domain.Ingredient{
		ID:   req.GetId(),
		Name: req.GetName(),
		Unit: req.GetUnit(),
	}

	err = d.repo.UpdateIngredient(ctx, i)
	return
}

// GetIngredient retrieves existing ingredient data.
func (d *dishesService) GetIngredient(ctx context.Context, id uint64) (resp *pb.IngredientV1, err error) {
	i, err := d.repo.GetIngredient(ctx, id)
	if err != nil {
		logger.Logger().Infoln("Failed to get ingredient ", id, " : ", err)
		return nil, fmt.Errorf("failed to get ingredient %d: %w", id, err)
	}
	resp = &pb.IngredientV1{
		Id:   i.ID,
		Name: i.Name,
		Unit: i.Unit,
	}
	return
}

// DeleteIngredient deletes Ingredient.
func (d *dishesService) DeleteIngredient(ctx context.Context, id uint64) (err error) {
	err = d.repo.DeleteIngredient(ctx, id)
	return
}

// ListIngredients returns list of all Ingredients.
func (d *dishesService) ListIngredients(ctx context.Context) (resp []*pb.IngredientV1, err error) {
	ingredients, err := d.repo.ListIngredients(ctx)
	if err != nil {
		logger.Logger().Infoln("Failed to list ingredients: ", err)
		return nil, fmt.Errorf("failed to list ingredients: %w", err)
	}
	resp = make([]*pb.IngredientV1, len(ingredients))
	for index, i := range ingredients {
		resp[index] = &pb.IngredientV1{
			Id:   i.ID,
			Name: i.Name,
			Unit: i.Unit,
		}
	}
	return
}

// StoreDish saves new Dish.
func (d *dishesService) StoreDish(ctx context.Context, req *pb.CreateDishRequestV1) (err error) {
	currDish := &domain.Dish{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}

	currDish.Ingredients = make([]*domain.DishIngredient, len(req.Ingredients))

	for index, i := range req.Ingredients {
		currDish.Ingredients[index] = &domain.DishIngredient{
			ID:   i.GetId(),
			Name: i.GetName(),
			Unit: i.GetName(),
			Qty:  i.GetQty(),
		}
	}

	err = d.repo.StoreDish(ctx, currDish)
	return
}

// UpdateDish updates existing Dish.
func (d *dishesService) UpdateDish(ctx context.Context, req *pb.UpdateDishRequestV1) (err error) {
	currDish := &domain.Dish{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}

	currDish.Ingredients = make([]*domain.DishIngredient, len(req.Ingredients))

	for index, i := range req.Ingredients {
		currDish.Ingredients[index] = &domain.DishIngredient{
			ID:   i.GetId(),
			Name: i.GetName(),
			Unit: i.GetName(),
			Qty:  i.GetQty(),
		}
	}

	err = d.repo.UpdateDish(ctx, currDish)
	return
}

// GetDish retrieves Dish info.
func (d *dishesService) GetDish(ctx context.Context, id uint64) (resp *pb.DishV1, err error) {
	currDish, err := d.repo.GetDish(ctx, id)
	if err != nil {
		logger.Logger().Infoln("Failed to get dish ", id, " : ", err)
		return nil, fmt.Errorf("failed to get dish %d: %w", id, err)
	}
	resp = &pb.DishV1{
		Id:          currDish.ID,
		Name:        currDish.Name,
		Description: currDish.Description,
	}

	resp.Ingredients = make([]*pb.DishIngredientV1, len(currDish.Ingredients))

	for index, i := range currDish.Ingredients {
		resp.Ingredients[index] = &pb.DishIngredientV1{
			Id:   i.ID,
			Name: i.Name,
			Unit: i.Unit,
			Qty:  i.Qty,
		}
	}

	return
}

// DeleteDish deletes Dish.
func (d *dishesService) DeleteDish(ctx context.Context, id uint64) (err error) {
	err = d.repo.DeleteDish(ctx, id)
	return
}

// ListDishes returns list of all Dishes, including nested Ingredients.
func (d *dishesService) ListDishes(ctx context.Context) (resp []*pb.DishV1, err error) {
	dishes, err := d.repo.ListDishes(ctx)
	if err != nil {
		logger.Logger().Infoln("Failed to list dishes: ", err)
		return nil, fmt.Errorf("failed to list dishes: %w", err)
	}
	resp = make([]*pb.DishV1, len(dishes))

	for dishIndex, dish := range dishes {
		resp[dishIndex] = &pb.DishV1{
			Id:          dish.ID,
			Name:        dish.Name,
			Description: dish.Description,
		}
		resp[dishIndex].Ingredients = make([]*pb.DishIngredientV1, len(dish.Ingredients))
		for index, i := range dish.Ingredients {
			resp[dishIndex].Ingredients[index] = &pb.DishIngredientV1{
				Id:   i.ID,
				Name: i.Name,
				Unit: i.Unit,
				Qty:  i.Qty,
			}
		}
	}
	return
}

// StoreBundle saves new Bundle.
func (d *dishesService) StoreBundle(ctx context.Context, req *pb.CreateBundleRequestV1) (err error) {
	bundle := &domain.Bundle{
		Name:  req.GetName(),
		Price: req.GetPrice(),
	}

	err = d.repo.StoreBundle(ctx, bundle)
	return
}

// UpdateBundle updates existing Bundle.
func (d *dishesService) UpdateBundle(ctx context.Context, req *pb.UpdateBundleRequestV1) (err error) {
	bundle := &domain.Bundle{
		ID:    req.GetId(),
		Name:  req.GetName(),
		Price: req.GetPrice(),
	}

	err = d.repo.StoreBundle(ctx, bundle)
	return
}

// GetBundle returns Bundle with passed id.
func (d *dishesService) GetBundle(ctx context.Context, id uint64) (resp *pb.BundleV1, err error) {
	bundle, err := d.repo.GetBundle(ctx, id)
	if err != nil {
		logger.Logger().Infoln("Failed to get bundle ", id, " : ", err)
		return nil, fmt.Errorf("failed to get bundle %d: %w", id, err)
	}
	resp = &pb.BundleV1{
		Id:    bundle.ID,
		Name:  bundle.Name,
		Price: bundle.Price,
	}
	return
}

// DeleteBundle removes Bundle.
func (d *dishesService) DeleteBundle(ctx context.Context, id uint64) (err error) {
	err = d.repo.DeleteBundle(ctx, id)
	return
}

// ListBundles returns array of all Bundles.
func (d *dishesService) ListBundles(ctx context.Context) (resp []*pb.BundleV1, err error) {
	bundles, err := d.repo.ListBundles(ctx)
	if err != nil {
		logger.Logger().Infoln("Failed to list bundles: ", err)
		return nil, fmt.Errorf("failed to list bundles: %w", err)
	}
	resp = make([]*pb.BundleV1, len(bundles))
	for index, bundle := range bundles {
		resp[index] = &pb.BundleV1{
			Id:    bundle.ID,
			Name:  bundle.Name,
			Price: bundle.Price,
		}
	}
	return
}

// AddBundleWeeklyDish adds Dish to Bundle for particular week.
func (d *dishesService) AddBundleWeeklyDish(ctx context.Context, weekNumber, bundleID, dishID uint64) (err error) {
	err = d.repo.AddBundleWeeklyDish(ctx, weekNumber, bundleID, dishID)
	return
}

// DeleteBundleWeeklyDish removes Dish from Bundle for particular week.
func (d *dishesService) DeleteBundleWeeklyDish(ctx context.Context, weekNumber, bundleID, dishID uint64) (err error) {
	err = d.repo.DeleteBundleWeeklyDish(ctx, weekNumber, bundleID, dishID)
	return
}

// GetBundleWeeklyDishes returns all Dishes from Bundle for particular week, including nested Ingredients info.
func (d *dishesService) GetBundleWeeklyDishes(ctx context.Context, weekNumber, bundleID uint64) (resp []*pb.DishV1, err error) {
	dishes, err := d.repo.GetBundleWeeklyDishes(ctx, weekNumber, bundleID)
	if err != nil {
		logger.Logger().Infoln("Failed to list bundle %d dishes for week: ", bundleID, weekNumber, err)
		return nil, fmt.Errorf("failed to list bundle %d dishes for week %d: %w", bundleID, weekNumber, err)
	}
	resp = make([]*pb.DishV1, len(dishes))
	for dishIndex, dish := range dishes {
		resp[dishIndex] = &pb.DishV1{
			Id:          dish.ID,
			Name:        dish.Name,
			Description: dish.Description,
		}
		resp[dishIndex].Ingredients = make([]*pb.DishIngredientV1, len(dish.Ingredients))
		for index, i := range dish.Ingredients {
			resp[dishIndex].Ingredients[index] = &pb.DishIngredientV1{
				Id:   i.ID,
				Name: i.Name,
				Unit: i.Unit,
				Qty:  i.Qty,
			}
		}
	}
	return
}

/*
	StoreResource(ctx context.Context, resource *pb.ResourceV1) error
	UpdateResource(ctx context.Context, resource *pb.ResourceV1) error
	GetResource(ctx context.Context, id uint64) (*pb.ResourceV1, error)
	DeleteResource(ctx context.Context, id uint64) error
	ListResources(ctx context.Context) ([]*pb.ResourceV1, error)

	AttachResourceToDish(ctx context.Context, resourceId, dishID uint64) error
	DetachResourceFromDish(ctx context.Context, resourceId, dishID uint64) error
	ListDishResources(ctx context.Context, dishID uint64) ([]*pb.ResourceV1, error)
*/
