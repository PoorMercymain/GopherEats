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
			ingredient.Name, ingredient.Unit, ingredient.ID)
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

	ingredient = &domain.Ingredient{ID: id}

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

		if err != nil {
			return subErrors.ErrorNoRowsWhileListing
		}

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
			err = rows.Scan(&i.ID, &i.Name, &i.Unit)
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

		if dish.Ingredients == nil || len(dish.Ingredients) == 0 {
			return
		}

		//save new ingredients info

		///batch
		batch := &pgx.Batch{}

		for _, value := range dish.Ingredients {
			batch.Queue("INSERT INTO dishes_ingredients (id, dish_id, ingredient_id, qty) VALUES(DEFAULT, $1, $2, $3)",
				dish.ID, value.ID, value.Qty)
		}

		br := tx.SendBatch(ctx, batch)
		defer br.Close()

		for range dish.Ingredients {
			_, err = br.Exec()
			if err != nil {
				logger.Logger().Errorln(err)
				return
			}
		}

		err = br.Close()
		if err != nil {
			logger.Logger().Errorln(err)
			return
		}
		//end batch

		return
	})
}
func (dbs *dbStorage) UpdateDish(ctx context.Context, dish *domain.Dish) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {

		commandTag, err := tx.Exec(ctx, "UPDATE dishes SET name = $1, description = $2, "+
			"updated_at = CURRENT_TIMESTAMP WHERE id = $3",
			dish.Name, dish.Description, dish.ID)
		if err != nil {
			return
		}

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		//in case we don't want to update ingredients info just pass nil
		if dish.Ingredients == nil {
			return
		}

		//delete old ingredients info
		_, err = tx.Exec(ctx, "DELETE FROM dishes_ingredients WHERE dish_id = $1", dish.ID)
		if err != nil {
			logger.Logger().Errorln(err)
		}

		if len(dish.Ingredients) == 0 {
			return
		}

		//save new ingredients info

		///batch
		batch := &pgx.Batch{}

		for _, value := range dish.Ingredients {
			batch.Queue("INSERT INTO dishes_ingredients (id, dish_id, ingredient_id, qty) VALUES(DEFAULT, $1, $2, $3)",
				dish.ID, value.ID, value.Qty)
		}

		br := tx.SendBatch(ctx, batch)
		defer br.Close()

		for range dish.Ingredients {
			_, err = br.Exec()
			if err != nil {
				logger.Logger().Errorln(err)
				return
			}
		}

		err = br.Close()
		if err != nil {
			logger.Logger().Errorln(err)
			return
		}
		//end batch

		return
	})
}

func (dbs *dbStorage) GetDish(ctx context.Context, id uint64) (dish *domain.Dish, err error) {
	dish = &domain.Dish{ID: id}

	err = dbs.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) (err error) {
		var totalRows int
		err = conn.QueryRow(ctx, "SELECT count(*) FROM dishes_ingredients WHERE dish_id = $1", id).
			Scan(&totalRows)

		if err != nil {
			return subErrors.ErrorWhileScanning
		}

		rows, err := conn.Query(ctx, "SELECT dishes.name, dishes.description, dishes_ingredients.qty, "+
			"ingredients.id, ingredients.name, ingredients.unit FROM dishes "+
			"LEFT OUTER JOIN dishes_ingredients ON dishes.ingredients.dish_id = dishes.id "+
			"JOIN ingredients ON ingredients.id = dishes_ingredients.ingredient_id "+
			"WHERE dishes.id = $1 ", id)

		dish.Ingredients = make([]*domain.DishIngredient, totalRows)

		counter := 0
		for rows.Next() {
			i := &domain.DishIngredient{}
			err = rows.Scan(&dish.Name, &dish.Description, &i.Qty, &i.ID, &i.Name, &i.Unit)
			if err != nil {
				return subErrors.ErrorWhileScanning
			}
			dish.Ingredients[counter] = i
			counter++
		}

		return
	})

	return
}
func (dbs *dbStorage) DeleteDish(ctx context.Context, id uint64) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {

		commandTagIng, err := tx.Exec(ctx, "DELETE FROM dishes_ingredients WHERE dish_id = $1", id)
		if err != nil {
			logger.Logger().Errorln(err)
		}

		commandTag, err := tx.Exec(ctx, "DELETE FROM dishes WHERE id = $1", id)

		if commandTag.RowsAffected() == 0 && commandTagIng.RowsAffected() == 0 {
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

		if err != nil {
			return subErrors.ErrorNoRowsWhileListing
		}

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
			d := &domain.Dish{}
			err = rows.Scan(&d.ID, &d.Name, &d.Description)
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
			bundle.Name, bundle.Price, bundle.ID)
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
	bundle = &domain.Bundle{ID: id}

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

		if err != nil {
			return subErrors.ErrorNoRowsWhileListing
		}

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
			err = rows.Scan(&b.ID, &b.Name, &b.Price)
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

func (dbs *dbStorage) AddBundleWeeklyDish(ctx context.Context, weekNumber, bundleID, dishID uint64) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {
		_, err = tx.Exec(ctx, "INSERT INTO bundles_dishes VALUES(DEFAULT, $1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			weekNumber, bundleID, dishID)
		return
	})
}
func (dbs *dbStorage) DeleteBundleWeeklyDish(ctx context.Context, weekNumber, bundleID, dishID uint64) (err error) {
	return dbs.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) (err error) {
		commandTag, err := tx.Exec(ctx, "DELETE FROM bundles_dishes WHERE week_number = $1 AND bundle_id = $2 AND dish_id = $3",
			weekNumber, bundleID, dishID)

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		return
	})
}
func (dbs *dbStorage) GetBundleWeeklyDishes(ctx context.Context, weekNumber, bundleID uint64) (dishes []*domain.Dish, err error) {
	err = dbs.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) (err error) {

		var totalRows int

		err = conn.QueryRow(ctx, "SELECT count(*) FROM bundles_dishes WHERE week_number = $1 AND bundle_id = $2",
			weekNumber, bundleID).
			Scan(&totalRows)

		if err != nil {
			return subErrors.ErrorNoRowsWhileListing
		}

		if totalRows == 0 {
			return subErrors.ErrorNoRowsWhileListing
		}

		rows, err := conn.Query(ctx, "SELECT dish_id FROM bundles_dishes WHERE week_number = $1 AND bundle_id = $2",
			weekNumber, bundleID)

		if err != nil {
			return subErrors.ErrorNoRowsWhileListing
		}

		dishes = make([]*domain.Dish, totalRows)
		counter := 0
		for rows.Next() {
			var dishID uint64
			err = rows.Scan(&dishID)
			if err != nil {
				return subErrors.ErrorWhileScanning
			}

			//not sure if it's the most optimised way to get the info, but join over 5 tables looks questionable either
			currDish, getErr := dbs.GetDish(ctx, dishID)
			if getErr != nil {
				return getErr
			}
			dishes[counter] = currDish
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
