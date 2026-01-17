package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/wb-go/wbf/dbpg"
	wbfretry "github.com/wb-go/wbf/retry"
	"github.com/yokitheyo/wb_level3_02/internal/domain"
	internalRetry "github.com/yokitheyo/wb_level3_02/internal/retry"
)

type PostgresURLRepository struct {
	db            *dbpg.DB
	retryStrategy wbfretry.Strategy
}

func NewPostgresURLRepository(db *dbpg.DB, strategy wbfretry.Strategy) URLRepository {
	if strategy.Attempts <= 0 {
		strategy = internalRetry.DefaultStrategy
	}
	return &PostgresURLRepository{
		db:            db,
		retryStrategy: strategy,
	}
}

func (r *PostgresURLRepository) Create(ctx context.Context, u *domain.URL) error {
	q := `INSERT INTO urls (short, original, created_at, expires_at)
		  VALUES ($1, $2, $3, $4)
		  RETURNING id, created_at`

	row := r.db.QueryRowContext(ctx, q, u.Short, u.Original, u.CreatedAt, u.ExpiresAt)
	return row.Scan(&u.ID, &u.CreatedAt)
}

func (r *PostgresURLRepository) FindByShort(ctx context.Context, short string) (*domain.URL, error) {
	q := `SELECT id, short, original, created_at, expires_at, visits FROM urls WHERE short = $1`

	row := r.db.QueryRowContext(ctx, q, short)
	u := &domain.URL{}
	err := row.Scan(&u.ID, &u.Short, &u.Original, &u.CreatedAt, &u.ExpiresAt, &u.Visits)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (r *PostgresURLRepository) IncrementVisits(ctx context.Context, id int64) error {
	q := `UPDATE urls SET visits = visits + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *PostgresURLRepository) AggregateByDay(ctx context.Context, short string, from, to time.Time) (map[string]int64, error) {
	q := `SELECT DATE(occurred_at)::TEXT as date, COUNT(*) as count
		  FROM clicks WHERE short = $1 AND occurred_at BETWEEN $2 AND $3
		  GROUP BY DATE(occurred_at)
		  ORDER BY date`

	rows, err := r.db.QueryContext(ctx, q, short, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var date string
		var count int64
		if err := rows.Scan(&date, &count); err != nil {
			return nil, err
		}
		result[date] = count
	}

	return result, rows.Err()
}

func (r *PostgresURLRepository) GetDeviceStats(ctx context.Context, short string, from, to time.Time) (map[string]int64, error) {
	q := `SELECT device, COUNT(*) as count FROM clicks
		  WHERE short = $1 AND occurred_at BETWEEN $2 AND $3
		  GROUP BY device`

	rows, err := r.db.QueryContext(ctx, q, short, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var device string
		var count int64
		if err := rows.Scan(&device, &count); err != nil {
			return nil, err
		}
		result[device] = count
	}

	return result, rows.Err()
}

func (r *PostgresURLRepository) SaveClick(ctx context.Context, c *domain.Click) error {
	q := `INSERT INTO clicks (url_id, short, occurred_at, user_agent, ip, referrer, device)
		  VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, q, c.URLID, c.Short, c.OccurredAt, c.UserAgent, c.IP, c.Referrer, c.Device)
	return err
}

func (r *PostgresURLRepository) GetRecentClicks(ctx context.Context, short string, limit int) ([]*domain.Click, error) {
	q := `SELECT id, url_id, short, occurred_at, user_agent, ip, referrer, device
		  FROM clicks WHERE short = $1
		  ORDER BY occurred_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, q, short, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clicks []*domain.Click
	for rows.Next() {
		click := &domain.Click{}
		err := rows.Scan(&click.ID, &click.URLID, &click.Short, &click.OccurredAt, &click.UserAgent, &click.IP, &click.Referrer, &click.Device)
		if err != nil {
			return nil, err
		}
		clicks = append(clicks, click)
	}

	return clicks, rows.Err()
}
