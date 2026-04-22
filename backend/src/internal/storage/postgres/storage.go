package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/OGZKTeBmj/forum/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCantExecInit = errors.New("can't exec init")
)

type Storage struct {
	db *pgxpool.Pool
}

func (s *Storage) MustConnect(ctx context.Context, url, user, password, dbName string) {
	if err := s.connect(ctx, url, user, password, dbName); err != nil {
		panic(err)
	}
}

func (s *Storage) Stop(ctx context.Context) {
	s.db.Close()
}

func (s *Storage) Init(ctx context.Context) error {
	const op = "postgres.Init"

	if _, err := s.db.Exec(ctx, QueryInit); err != nil {
		return utils.ErrWrap(op, ErrCantExecInit)
	}
	return nil
}

func (s *Storage) connect(ctx context.Context, url, user, password, dbName string) error {
	const op = "postgres.Connect"

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, url, dbName)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return utils.ErrWrap(op, err)
	}

	s.db = pool
	return nil
}

const (
	QueryInit = `
	CREATE TABLE IF NOT EXISTS users (
    	id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    	name      VARCHAR(35) NOT NULL,
    	pass_hash BYTEA NOT NULL,
    	avatar_original_path TEXT DEFAULT '',
    	avatar_thumbnail_path TEXT DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS posts (
    	id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    	author_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    	title      TEXT NOT NULL,
    	content    TEXT NOT NULL,
    	time_stamp TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS comments (
    	id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    	post_id    BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    	author_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    	content    TEXT NOT NULL,
    	time_stamp TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS votes (
    	post_id    BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    	author_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    	value      SMALLINT NOT NULL CHECK (value IN (-1, 1)),
    	PRIMARY KEY (post_id, author_id)
	);
	`
)
