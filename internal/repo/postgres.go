package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wb-go/wbf/dbpg"
	wbfretry "github.com/wb-go/wbf/retry"
	"github.com/yokitheyo/wb_level3_02/internal/model"
	internalRetry "github.com/yokitheyo/wb_level3_02/internal/retry"
)

type PostgresRepo struct {
	db            *dbpg.DB
	retryStrategy wbfretry.Strategy
}

func NewPostgresRepo(db *dbpg.DB, strategy wbfretry.Strategy) *PostgresRepo {
	if strategy.Attempts <= 0 {
		strategy = internalRetry.DefaultStrategy
	}
	return &PostgresRepo{db: db, retryStrategy: strategy}
}

func (r *PostgresRepo) Create(ctx context.Context, u *model.URL) error {
	q := `INSERT INTO urls (short, original, created_at, expires_at)
		  VALUES ($1, $2, now(), $3)
		  RETURNING id, created_at
		  `

	var lastErr error

	attempts := r.retryStrategy.Attempts
	if attempts <= 0 {
		attempts = internalRetry.DefaultStrategy.Attempts
		if attempts <= 0 {
			attempts = 1
		}
	}

	delay := r.retryStrategy.Delay
	if delay <= 0 {
		delay = internalRetry.DefaultStrategy.Delay
		if delay <= 0 {
			delay = 100 * time.Millisecond
		}
	}

	backoff := r.retryStrategy.Backoff
	if backoff <= 0 {
		backoff = internalRetry.DefaultStrategy.Backoff
		if backoff <= 0 {
			backoff = 2.0
		}
	}

	for i := 0; i < attempts; i++ {
		row := r.db.Master.QueryRowContext(ctx, q, u.Short, u.Original, u.ExpiresAt)
		if err := row.Scan(&u.ID, &u.CreatedAt); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		time.Sleep(delay)
		delay = time.Duration(float64(delay) * backoff)
	}

	return fmt.Errorf("create url failed after %d attempts: %w", attempts, lastErr)
}

func (r *PostgresRepo) FindByShort(ctx context.Context, short string) (*model.URL, error) {
	q := `SELECT id, short, original, created_at, expires_at, visits
          FROM urls 
          WHERE short = $1 AND is_disabled = false
          LIMIT 1
	     `

	rows, err := r.db.QueryWithRetry(ctx, r.retryStrategy, q, short)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	u := &model.URL{}
	if err := rows.Scan(&u.ID, &u.Short, &u.Original, &u.CreatedAt, &u.ExpiresAt, &u.Visits); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepo) IncrementVisits(ctx context.Context, urlID int64) error {
	_, err := r.db.ExecWithRetry(ctx, r.retryStrategy, `UPDATE urls SET visits = visits + 1 WHERE id = $1`, urlID)
	return err
}

func (r *PostgresRepo) InsertClick(ctx context.Context, c *model.Click) error {
	_, err := r.db.ExecWithRetry(ctx, r.retryStrategy, `
	INSERT INTO clicks (url_id, short, occurred_at, user_agent, ip, referrer, device)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, c.URLID, c.Short, c.Occurred, c.UserAgent, c.IP, c.Referrer, c.Device)
	return err
}

func (r *PostgresRepo) AggregateByDay(ctx context.Context, short string, from, to time.Time) (map[string]int64, error) {
	rows, err := r.db.QueryWithRetry(ctx, r.retryStrategy, `
    SELECT date_trunc('day', occurred_at) AS day, count(*)
	FROM clicks 
    WHERE short = $1 AND occurred_at BETWEEN $2 AND $3
    GROUP BY day
    ORDER BY day
	`, short, from, to)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make(map[string]int64)
	for rows.Next() {
		var day time.Time
		var cnt int64
		if err := rows.Scan(&day, &cnt); err != nil {
			return nil, err
		}
		res[day.Format("2006-01-02")] = cnt
	}
	return res, nil
}
