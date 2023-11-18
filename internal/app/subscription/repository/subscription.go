// Package repository contains repository handling for subscription service.
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

// DB returns new connection pool.
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

// New returns pointer to subscription repo struct with connection pool.
func New(pool *pgxpool.Pool) *subscription {
	return &subscription{pool: pool}
}

// WithTransaction wraps postgres SQL requests into transaction.
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
		if !errors.Is(err, pgx.ErrTxClosed) && err != nil {
			logger.Logger().Errorln(err)
		}
	}()

	err = txFunc(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// WithConnection wraps postgres SQL queries with connection.
func (r *subscription) WithConnection(ctx context.Context, connFunc func(context.Context, *pgxpool.Conn) error) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return connFunc(ctx, conn)
}

// CreateSubscription creates new subscription.
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

// ReadSubscription returns subscription info.
func (r *subscription) ReadSubscription(ctx context.Context, email string) (int64, bool, error) {
	var bundleID int64
	var isDeleted bool

	err := r.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		err := conn.QueryRow(ctx, "SELECT bundle_id, is_deleted FROM subscriptions WHERE email = $1", email).Scan(&bundleID, &isDeleted)

		if errors.Is(err, pgx.ErrNoRows) {
			return subErrors.ErrorNoRowsWhileReading
		}

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, false, err
	}

	return bundleID, isDeleted, nil
}

// UpdateSubscription updates subscription.
func (r *subscription) UpdateSubscription(ctx context.Context, email string, bundleID int64, isDeleted bool) error {
	return r.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		commandTag, err := tx.Exec(ctx, "UPDATE subscriptions SET bundle_id = $1, is_deleted = $2 WHERE email = $3", bundleID, isDeleted, email)
		if err != nil {
			return err
		}

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNoRowsUpdated
		}

		return nil
	})
}

// DeleteSubscription removes subscription for user with passed email.
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

// AddBalance allows to add funds for user balance.
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

// ReadUserData returns user data.
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

// ReadBalanceHistory returns history of balance funding and purchases for user with passed email.
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

// TODO: use struct of four channels: one to send bundle_ids to handler (to request their prices from dishes),
// another - to receive the prices, the next - to send emails to handler to cancel unpaid subscriptions, and the last - to send email, bundle_id to handler,
// combine them with week_number and address and send to dishes service

// ChargeForSubscription allows to charge user balance.
func (r *subscription) ChargeForSubscription(ctx context.Context, notEnoughFundsEmailsChan chan<- string) error {
	return r.WithConnection(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		rows, err := conn.Query(ctx, "SELECT email, bundle_id FROM subscriptions WHERE is_deleted = $1", false)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		var email string
		var bundleID uint64

		for rows.Next() {
			err := rows.Scan(&email, &bundleID)
			if err != nil {
				return err
			}

			err = r.RemoveAmountFromBalance(ctx, email, 10) // TODO: use actual price
			logger.Logger().Infoln(email)

			if errors.Is(err, subErrors.ErrorNotEnoughFunds) { //TODO: send to email chan
				logger.Logger().Infoln("opa")
				notEnoughFundsEmailsChan <- email
				err = r.DeleteSubscription(ctx, email)
				if err != nil {
					return err
				}
				logger.Logger().Infoln(email)
			}

			if !errors.Is(err, subErrors.ErrorNotEnoughFunds) && err != nil {
				logger.Logger().Infoln(email)
				return err
			}
			// TODO: send info to dishes service
		}

		return nil
	})
}

// RemoveAmountFromBalance decreases balance.
func (r *subscription) RemoveAmountFromBalance(ctx context.Context, email string, amount uint64) error {
	return r.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		commandTag, err := tx.Exec(ctx, "UPDATE balances SET current_balance = current_balance - $1 WHERE email = $2 AND current_balance > $1", amount, email)
		if err != nil {
			return err
		}

		if commandTag.RowsAffected() == 0 {
			return subErrors.ErrorNotEnoughFunds
		}

		_, err = tx.Exec(ctx, "INSERT INTO balances_history VALUES($1, $2, $3, $4)", email, amount, "debit", time.Now())
		if err != nil {
			return err
		}

		return nil
	})
}
