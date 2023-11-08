package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/PoorMercymain/GopherEats/internal/app/dishes/domain"
	"github.com/PoorMercymain/GopherEats/internal/pkg/logger"
)

type dbStorage struct {
	pgxPool *pgxpool.Pool
}

func NewDBStorage(DSN string) (storage domain.DishesRepository, err error) {
	pg, err := sql.Open("pgx", DSN)
	if err != nil {
		logger.Logger().Infoln(err)
		return
	}
	err = goose.SetDialect("postgres")
	if err != nil {
		logger.Logger().Infoln(err)
		return
	}

	err = pg.PingContext(context.Background())
	if err != nil {
		logger.Logger().Infoln(err)

		for _, retryInterval := range domain.RepeatedAttemptsIntervals {
			time.Sleep(retryInterval)
			err = pg.PingContext(context.Background())
			if err != nil {
				logger.Logger().Infoln(err)
			} else {
				logger.Logger().Infoln("ping successful")
				break
			}

		}
		if err != nil {
			logger.Logger().Infoln(err)
			return
		}
	}

	const migrationsPath = "./internal/app/dishes/repository/migration"

	err = goose.Run("up", pg, migrationsPath)
	if err != nil {
		logger.Logger().Infoln(err)
		return
	}
	pg.Close()

	config, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		logger.Logger().Infoln("Error parsing DSN:", err)
		return
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Logger().Infoln("Error creating pgxpool:", err)
		return
	}

	logger.Logger().Infoln("Pool created", pool, err)

	storage = &dbStorage{pgxPool: pool}
	return
}

func (dbs *dbStorage) StoreIngredient(ctx context.Context, ingredient *domain.Ingredient) (err error) {
	return
}
func (dbs *dbStorage) UpdateIngredient(ctx context.Context, ingredient *domain.Ingredient) (err error) {
	return
}
func (dbs *dbStorage) GetIngredient(ctx context.Context, id uint64) (ingredient *domain.Ingredient, err error) {
	return
}
func (dbs *dbStorage) DeleteIngredient(ctx context.Context, id uint64) (err error) {
	return
}
func (dbs *dbStorage) ListIngredients(ctx context.Context) (ingredients []*domain.Ingredient, err error) {
	return
}

func (dbs *dbStorage) StoreDish(ctx context.Context, dish *domain.Dish) (err error) {
	return
}
func (dbs *dbStorage) UpdateDish(ctx context.Context, dish *domain.Dish) (err error) {
	return
}
func (dbs *dbStorage) GetDish(ctx context.Context, id uint64) (dish *domain.Dish, err error) {
	return
}
func (dbs *dbStorage) DeleteDish(ctx context.Context, id uint64) (err error) {
	return
}
func (dbs *dbStorage) ListDishes(ctx context.Context) (dishes []*domain.Dish, err error) {
	return
}

func (dbs *dbStorage) StoreBundle(ctx context.Context, bundle *domain.Bundle) (err error) {
	return
}
func (dbs *dbStorage) UpdateBundle(ctx context.Context, bundle *domain.Bundle) (err error) {
	return
}
func (dbs *dbStorage) GetBundle(ctx context.Context, id uint64) (bundle *domain.Bundle, err error) {
	return
}
func (dbs *dbStorage) DeleteBundle(ctx context.Context, id uint64) (err error) {
	return
}
func (dbs *dbStorage) ListBundles(ctx context.Context) (bundles []*domain.Bundle, err error) {
	return
}

func (dbs *dbStorage) AddBundleWeeklyDish(ctx context.Context, weekNumber, bundleId, dishId uint64) (err error) {
	return
}
func (dbs *dbStorage) DeleteBundleWeeklyDish(ctx context.Context, weekNumber, bundleId, dishId uint64) (err error) {
	return
}
func (dbs *dbStorage) GetBundleWeeklyDishes(ctx context.Context, weekNumber, bundleId uint64) (dishes []*domain.Dish, err error) {
	return
}

func (dbs *dbStorage) ClosePool() (err error) {
	dbs.pgxPool.Close()
	return
}

func (dbs *dbStorage) Ping(ctx context.Context) (err error) {
	conn, err := dbs.pgxPool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	err = conn.Ping(ctx)
	if err != nil {
		return
	}
	return
}
