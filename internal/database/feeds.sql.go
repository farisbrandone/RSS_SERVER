// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeeds = `-- name: CreateFeeds :one
INSERT INTO  feeds (id, created_at, updated_at, name, url , user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at, name, url, user_id, last_fetched_at
`

type CreateFeedsParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      sql.NullString
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) CreateFeeds(ctx context.Context, arg CreateFeedsParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeeds,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const getAllFeeds = `-- name: GetAllFeeds :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at
FROM feeds
`

func (q *Queries) GetAllFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getAllFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextFeedsToFetch = `-- name: GetNextFeedsToFetch :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1
`

func (q *Queries) GetNextFeedsToFetch(ctx context.Context, limit int32) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getNextFeedsToFetch, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getoneFeeds = `-- name: GetoneFeeds :one
SELECT id, created_at, updated_at, name, url , user_id
FROM feeds
WHERE id=$1
`

type GetoneFeedsRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      sql.NullString
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) GetoneFeeds(ctx context.Context, id uuid.UUID) (GetoneFeedsRow, error) {
	row := q.db.QueryRowContext(ctx, getoneFeeds, id)
	var i GetoneFeedsRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const markFeedFetched = `-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, name, url, user_id, last_fetched_at
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id uuid.UUID) (Feed, error) {
	row := q.db.QueryRowContext(ctx, markFeedFetched, id)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}