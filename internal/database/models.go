// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          sql.NullString
	Url           string
	UserID        uuid.UUID
	LastFetchedAt sql.NullTime
}

type FeedFollow struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FeedID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt sql.NullTime
	FeedID      uuid.UUID
}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      sql.NullString
	ApiKey    string
}
