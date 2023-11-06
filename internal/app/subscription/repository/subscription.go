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

	"github.com/PoorMercymain/GopherEats/internal/app/subscription/domain"
	subErrors "github.com/PoorMercymain/GopherEats/internal/app/subscription/errors"
	"github.com/PoorMercymain/GopherEats/internal/pkg/logger"
	api "github.com/PoorMercymain/GopherEats/pkg/api/subscription"
)

var _ domain.SubscriptionRepository = (*subscription)(nil)

type subscription struct {
	pool *pgxpool.Pool
}

func DB(dsn string) (*pgxpool.Pool, error) {
	pg, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = goose.SetDialect("postgres")
	if err != nil {
		return nil, err
	}

	err = pg.PingContext(context.Background())
	if err != nil {
		return nil, err
	}

	err = goose.Run("up", pg, "./internal/app/subscription/repository/migrations")
	if err != nil {
		return nil, err
	}
	pg.Close()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, err
}

func New(pool *pgxpool.Pool) *subscription {
	return &subscription{pool: pool}
}

func (r *subscription) WithTransaction(ctx context.Context, txFunc func(context.Context, pgx.Tx) error) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if !errors.Is(err, pgx.ErrTxClosed) {
			logger.Logger().Errorln(err)
		}
	}()

	err = txFunc(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *subscription) WithConnection(ctx context.Context, connFunc func(context.Context, *pgxpool.Conn) error) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return connFunc(ctx, conn)
}

func (r *subscription) CreateSubscription(ctx context.Context, email string, bundleID int64) error {
	return r.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO subscriptions VALUES($1, $2, $3)", email, bundleID, false)
		if err != nil {
			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					return subErrors.ErrorUniqueViolationWhileCreating
				}
			}

			return err
		}

		return nil
	})
}

func (r *subscription) ReadSubscription(ctx context.Context, email string) (int64, error) {
	var bundleID int64

	err := r.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		err := conn.QueryRow(ctx, "SELECT bundle_id FROM subscriptions WHERE email = $1 AND is_deleted = $2", email, false).Scan(&bundleID)

		if errors.Is(err, pgx.ErrNoRows) {
			return subErrors.ErrorNoRowsWhileReading
		}

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return bundleID, nil
}

func (r *subscription) UpdateSubscription(ctx context.Context, email string, bundleID int64) error {
	return r.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		commandTag, err := tx.Exec(ctx, "UPDATE subscriptions SET bundle_id = $1 WHERE email = $2", bundleID, email)
		if err != nil {
			return err
		}

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		return nil
	})
}

func (r *subscription) DeleteSubscription(ctx context.Context, email string) error {
	return r.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		commandTag, err := tx.Exec(ctx, "UPDATE subscriptions SET is_deleted = $1 WHERE email = $2 AND is_deleted = $3", true, email, false)
		if err != nil {
			return err
		}

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		return nil
	})
}

func (r *subscription) AddBalance(ctx context.Context, email string, amount uint64) error {
	return r.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO balances (email, current_balance) VALUES($1, $2) ON CONFLICT(email) DO UPDATE SET current_balance = balances.current_balance + EXCLUDED.current_balance", email, amount)

		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, "INSERT INTO balances_history VALUES($1, $2, $3, $4)", email, amount, "replenishment", time.Now())

		return err
	})
}

func (r *subscription) ReadUserData(ctx context.Context, email string) (domain.UserData, error) {
	var userData domain.UserData

	err := r.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		err := conn.QueryRow(ctx, "SELECT bundle_id FROM subscriptions WHERE email = $1", email).Scan(&userData.BundleID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		err = conn.QueryRow(ctx, "SELECT current_balance FROM balances WHERE email = $1", email).Scan(&userData.Balance)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		return nil
	})

	if err != nil {
		return domain.UserData{}, err
	}

	return userData, nil
}

func (r *subscription) ReadBalanceHistory(ctx context.Context, email string, page uint64) ([]*api.HistoryElementV1, error) {
	var history = make([]*api.HistoryElementV1, 0)
	err := r.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		rows, err := conn.Query(ctx, "SELECT email, amount, operation, made_at FROM balances_history WHERE email = $1 ORDER BY made_at DESC LIMIT 15 OFFSET $2", email, (page-1)*15)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		for rows.Next() {
			var historyElem api.HistoryElementV1
			var t time.Time
			err := rows.Scan(&historyElem.Email, &historyElem.Amount, &historyElem.Operation, &t)
			if err != nil {
				return err
			}

			historyElem.MadeAt = t.Format(time.RFC3339)

			history = append(history, &historyElem)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(history) == 0 {
		return nil, subErrors.ErrorNoRowsWhileReadingHistory
	}
	return history, nil
}
