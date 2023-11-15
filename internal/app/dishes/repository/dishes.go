package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/PoorMercymain/GopherEats/internal/app/dishes/domain"
	subErrors "github.com/PoorMercymain/GopherEats/internal/app/dishes/errors"
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

func (dbs *dbStorage) WithTransaction(ctx context.Context, txFunc func(context.Context, pgx.Tx) error) (err error) {
	conn, err := dbs.pgxPool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		err = tx.Rollback(ctx)
		if !errors.Is(err, pgx.ErrTxClosed) && err != nil {
			logger.Logger().Errorln(err)
		}
	}()

	err = txFunc(ctx, tx)
	if err != nil {
		return
	}

	return tx.Commit(ctx)
}

func (dbs *dbStorage) WithConnection(ctx context.Context, connFunc func(context.Context, *pgxpool.Conn) error) (err error) {
	conn, err := dbs.pgxPool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	return connFunc(ctx, conn)
}

func (dbs *dbStorage) StoreIngredient(ctx context.Context, ingredient *domain.Ingredient) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {
		_, err = tx.Exec(ctx, "INSERT INTO ingredients VALUES(DEFAULT, $1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			ingredient.Name, ingredient.Unit)
		if err != nil {
			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					return subErrors.ErrorUniqueViolationWhileStoring
				}
			}

			return
		}

		return
	})
}

func (dbs *dbStorage) UpdateIngredient(ctx context.Context, ingredient *domain.Ingredient) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {

		commandTag, err := tx.Exec(ctx, "UPDATE ingredients SET name = $1, unit = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3",
			ingredient.Name, ingredient.Unit, ingredient.Id)
		if err != nil {
			return
		}

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		return
	})
}
func (dbs *dbStorage) GetIngredient(ctx context.Context, id uint64) (ingredient *domain.Ingredient, err error) {

	ingredient = &domain.Ingredient{Id: id}

	err = dbs.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) (err error) {
		err = conn.QueryRow(ctx, "SELECT name, unit FROM ingredients WHERE id = $1 LIMIT 1", id).
			Scan(&ingredient.Name, &ingredient.Unit)

		if errors.Is(err, pgx.ErrNoRows) {
			return subErrors.ErrorNoRowsWhileGetting
		}

		return
	})

	return
}

func (dbs *dbStorage) DeleteIngredient(ctx context.Context, id uint64) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {

		commandTag, err := tx.Exec(ctx, "DELETE FROM ingredients WHERE id = $1", id)
		if err != nil {
			return
		}

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		return
	})
}
func (dbs *dbStorage) ListIngredients(ctx context.Context) (ingredients []*domain.Ingredient, err error) {

	err = dbs.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) (err error) {

		var totalRows int

		err = conn.QueryRow(ctx, "SELECT count(*) FROM ingredients").
			Scan(&totalRows)

		if totalRows == 0 {
			return subErrors.ErrorNoRowsWhileListing
		}

		rows, err := conn.Query(ctx, "SELECT id, name, unit FROM ingredients")

		if err != nil {
			return subErrors.ErrorNoRowsWhileListing
		}

		ingredients = make([]*domain.Ingredient, totalRows)
		counter := 0
		for rows.Next() {
			i := &domain.Ingredient{}
			err = rows.Scan(&i.Id, &i.Name, &i.Unit)
			if err != nil {
				return subErrors.ErrorWhileScanning
			}
			ingredients[counter] = i
			counter++
		}

		return
	})

	return
}

func (dbs *dbStorage) StoreDish(ctx context.Context, dish *domain.Dish) (err error) {
	//TODO: process DishIngredients
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {
		var id int
		err = tx.QueryRow(ctx, "INSERT INTO dishes VALUES(DEFAULT, $1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) "+
			"RETURNING id",
			dish.Name, dish.Description).Scan(&id)
		if err != nil {
			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					return subErrors.ErrorUniqueViolationWhileStoring
				}
			}

			return
		}

		return
	})
}
func (dbs *dbStorage) UpdateDish(ctx context.Context, dish *domain.Dish) (err error) {

	//TODO: process DishIngredients
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {

		commandTag, err := tx.Exec(ctx, "UPDATE dishes SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3", dish.Name, dish.Description, dish.Id)
		if err != nil {
			return
		}

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		return
	})
}
func (dbs *dbStorage) GetDish(ctx context.Context, id uint64) (dish *domain.Dish, err error) {
	//TODO: process DishIngredients
	dish = &domain.Dish{Id: id}

	err = dbs.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) (err error) {
		err = conn.QueryRow(ctx, "SELECT name, description FROM dishes WHERE id = $1 LIMIT 1", id).
			Scan(&dish.Name, &dish.Description)

		if errors.Is(err, pgx.ErrNoRows) {
			return subErrors.ErrorNoRowsWhileGetting
		}

		return
	})

	return
}
func (dbs *dbStorage) DeleteDish(ctx context.Context, id uint64) (err error) {
	//TODO: process DishIngredients
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {

		commandTag, err := tx.Exec(ctx, "DELETE FROM dishes WHERE id = $1", id)

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		return
	})
}
func (dbs *dbStorage) ListDishes(ctx context.Context) (dishes []*domain.Dish, err error) {
	err = dbs.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) (err error) {

		var totalRows int

		err = conn.QueryRow(ctx, "SELECT count(*) FROM dishes").
			Scan(&totalRows)

		if totalRows == 0 {
			return subErrors.ErrorNoRowsWhileListing
		}

		rows, err := conn.Query(ctx, "SELECT id, name, description FROM dishes")

		if err != nil {
			return subErrors.ErrorNoRowsWhileListing
		}

		dishes = make([]*domain.Dish, totalRows)
		counter := 0
		for rows.Next() {
			i := &domain.Dish{}
			err = rows.Scan(&i.Id, &i.Name, &i.Description)
			if err != nil {
				return subErrors.ErrorWhileScanning
			}
			dishes[counter] = i
			counter++
		}

		return
	})

	return
}

func (dbs *dbStorage) StoreBundle(ctx context.Context, bundle *domain.Bundle) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {
		_, err = tx.Exec(ctx, "INSERT INTO bundles VALUES(DEFAULT, $1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			bundle.Name, bundle.Price)
		if err != nil {
			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					return subErrors.ErrorUniqueViolationWhileStoring
				}
			}

			return
		}

		return
	})
}
func (dbs *dbStorage) UpdateBundle(ctx context.Context, bundle *domain.Bundle) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {

		commandTag, err := tx.Exec(ctx, "UPDATE bundles SET name = $1, price = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3",
			bundle.Name, bundle.Price, bundle.Id)
		if err != nil {
			return
		}

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		return
	})
}
func (dbs *dbStorage) GetBundle(ctx context.Context, id uint64) (bundle *domain.Bundle, err error) {
	bundle = &domain.Bundle{Id: id}

	err = dbs.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) (err error) {
		err = conn.QueryRow(ctx, "SELECT name, price FROM bundles WHERE id = $1 LIMIT 1", id).
			Scan(&bundle.Name, &bundle.Price)

		if errors.Is(err, pgx.ErrNoRows) {
			return subErrors.ErrorNoRowsWhileGetting
		}

		return
	})

	return
}
func (dbs *dbStorage) DeleteBundle(ctx context.Context, id uint64) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {

		commandTag, err := tx.Exec(ctx, "DELETE FROM bundles WHERE id = $1", id)

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		return
	})
}
func (dbs *dbStorage) ListBundles(ctx context.Context) (bundles []*domain.Bundle, err error) {
	err = dbs.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) (err error) {

		var totalRows int

		err = conn.QueryRow(ctx, "SELECT count(*) FROM bundles").
			Scan(&totalRows)

		if totalRows == 0 {
			return subErrors.ErrorNoRowsWhileListing
		}

		rows, err := conn.Query(ctx, "SELECT id, name, price FROM bundles")

		if err != nil {
			return subErrors.ErrorNoRowsWhileListing
		}

		bundles = make([]*domain.Bundle, totalRows)
		counter := 0
		for rows.Next() {
			b := &domain.Bundle{}
			err = rows.Scan(&b.Id, &b.Name, &b.Price)
			if err != nil {
				return subErrors.ErrorWhileScanning
			}
			bundles[counter] = b
			counter++
		}

		return
	})

	return
}

func (dbs *dbStorage) AddBundleWeeklyDish(ctx context.Context, weekNumber, bundleId, dishId uint64) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {
		_, err = tx.Exec(ctx, "INSERT INTO bundles_dishes VALUES(DEFAULT, $1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			weekNumber, bundleId, dishId)
		return
	})
}
func (dbs *dbStorage) DeleteBundleWeeklyDish(ctx context.Context, weekNumber, bundleId, dishId uint64) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {
		commandTag, err := tx.Exec(ctx, "DELETE FROM bundles_dishes WHERE week_number = $1 AND bundle_id = $2 AND dish_id = $3",
			weekNumber, bundleId, dishId)

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		return
	})
}
func (dbs *dbStorage) GetBundleWeeklyDishes(ctx context.Context, weekNumber, bundleId uint64) (dishes []*domain.Dish, err error) {
	err = dbs.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) (err error) {

		var totalRows int

		err = conn.QueryRow(ctx, "SELECT count(*) FROM bundles_dishes WHERE week_number = $1 AND bundle_id = $2",
			weekNumber, bundleId).
			Scan(&totalRows)

		if totalRows == 0 {
			return subErrors.ErrorNoRowsWhileListing
		}

		//TODO: JOIN to get full info from dishes and dishes_ingredients table
		rows, err := conn.Query(ctx, "SELECT dish_id FROM bundles_dishes WHERE week_number = $1 AND bundle_id = $2",
			weekNumber, bundleId)

		if err != nil {
			return subErrors.ErrorNoRowsWhileListing
		}

		dishes = make([]*domain.Dish, totalRows)
		counter := 0
		for rows.Next() {
			d := &domain.Dish{}
			err = rows.Scan(&d.Id)
			if err != nil {
				return subErrors.ErrorWhileScanning
			}
			dishes[counter] = d
			counter++
		}

		return
	})

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
