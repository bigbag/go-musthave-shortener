package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

type pgRepository struct {
	ctx         context.Context
	conn        *sql.DB
	connTimeout time.Duration
}

func NewPGRepository(
	ctx context.Context,
	databaseDSN string,
	connTimeout time.Duration,
) (StorageRepository, error) {
	conn, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		return nil, err
	}

	repo := &pgRepository{ctx: ctx, conn: conn, connTimeout: connTimeout * time.Second}
	err = repo.init()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *pgRepository) init() error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	_, err := r.conn.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS
			urls(
				key VARCHAR NOT NULL,
				value VARCHAR NOT NULL,
				user_id VARCHAR NOT NULL,
				correlation_id VARCHAR NULL,
				removed BOOL DEFAULT 'f',
				PRIMARY KEY (key),
				UNIQUE (value)
			);`,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *pgRepository) GetByKey(key string) (*Record, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	record := &Record{}

	sqlStatement := `SELECT key, value, user_id, removed FROM urls WHERE key=$1;`
	row := r.conn.QueryRowContext(ctx, sqlStatement, key)
	switch err := row.Scan(
		&record.Key,
		&record.Value,
		&record.UserID,
		&record.Removed,
	); err {
	case sql.ErrNoRows:
		return nil, errors.New("not found url")
	case nil:
		return record, nil
	default:
		return record, err
	}
}

func (r *pgRepository) GetByValue(value string) (*Record, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	record := &Record{}

	sqlStatement := `SELECT key, value, user_id FROM urls WHERE value=$1;`
	row := r.conn.QueryRowContext(ctx, sqlStatement, value)
	switch err := row.Scan(&record.Key, &record.Value, &record.UserID); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return record, nil
	default:
		return record, err
	}
}

func (r *pgRepository) GetAllByUserID(userID string) ([]*Record, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	sqlStatement := `SELECT key, value, user_id FROM urls WHERE user_id=$1;`
	rows, err := r.conn.QueryContext(ctx, sqlStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var record *Record
	result := make([]*Record, 0, 100)
	for rows.Next() {
		record = &Record{}
		err = rows.Scan(&record.Key, &record.Value, &record.UserID)
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *pgRepository) Save(record *Record) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	query := `INSERT INTO urls(key, value, user_id) 
          			VALUES($1, $2, $3);`
	_, err := r.conn.ExecContext(ctx, query, record.Key, record.Value, record.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *pgRepository) SaveBatchOfURL(records []*Record) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO urls(key, value, user_id, correlation_id)
				VALUES($1, $2, $3, $4);`,
	)
	if err != nil {
		return err
	}

	for _, record := range records {
		if _, err = stmt.ExecContext(
			ctx, record.Key,
			record.Value,
			record.UserID,
			record.CorrelationID,
		); err != nil {
			return err
		}
	}

	defer stmt.Close()

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *pgRepository) DeleteByUserID(userID string, keys []string) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	stmt, err := tx.PrepareContext(
		ctx,
		`UPDATE urls 
				SET removed = true 
				WHERE user_id = $1 and key = $2;`,
	)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if _, err = stmt.ExecContext(
			ctx, userID, key); err != nil {
			return err
		}
	}

	defer stmt.Close()

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *pgRepository) Status() error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	return r.conn.PingContext(ctx)
}

func (r *pgRepository) Close() error {
	return r.conn.Close()
}
