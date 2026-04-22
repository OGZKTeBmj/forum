package postgres

import (
	"context"
	"errors"

	models "github.com/OGZKTeBmj/forum/domain"
	"github.com/OGZKTeBmj/forum/internal/storage"
	"github.com/OGZKTeBmj/forum/utils"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) User(ctx context.Context, id []byte) (models.User, error) {
	const op = "postgres.User"

	var user models.User

	if err := s.db.QueryRow(ctx, QueryUser, id).Scan(&user.Id, &user.Name,
		&user.PassHash, &user.AvatarPath.Original,
		&user.AvatarPath.Thumbnail); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, utils.ErrWrap(op, storage.ErrIsNotExist)
		}
		return models.User{}, utils.ErrWrap(op, err)
	}
	return user, nil
}

func (s *Storage) UserByName(ctx context.Context, name string) (models.User, error) {
	const op = "postgres.UserByName"

	var user models.User

	if err := s.db.QueryRow(ctx, QueryUserByName, name).Scan(&user.Id, &user.Name,
		&user.PassHash, &user.AvatarPath.Original, &user.AvatarPath.Thumbnail); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, utils.ErrWrap(op, storage.ErrIsNotExist)
		}
		return models.User{}, utils.ErrWrap(op, err)
	}
	return user, nil
}

func (s *Storage) SaveUser(ctx context.Context, name string, passhash []byte) ([]byte, error) {
	const op = "postgres.SaveUser"

	var id []byte

	if err := s.db.QueryRow(ctx, QuerySaveUser, name, passhash).Scan(&id); err != nil {
		return nil, utils.ErrWrap(op, err)
	}
	return id, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user models.User) error {
	const op = "postgres.UpdateUser"

	res, err := s.db.Exec(ctx, QueryUpdateUser, user.Id, user.Name, user.PassHash,
		user.AvatarPath.Original, user.AvatarPath.Thumbnail)
	if err != nil {
		return utils.ErrWrap(op, err)
	}
	if rows := res.RowsAffected(); rows == 0 {
		return storage.ErrIsNotExist
	}
	return nil
}

const (
	QueryUser = `
	SELECT id, name, pass_hash, avatar_original_path, avatar_thumbnail_path FROM users
	WHERE id = $1;
	`

	QueryUserByName = `
	SELECT id, name, pass_hash, avatar_original_path, avatar_thumbnail_path FROM users
	WHERE name = $1;
	`

	QuerySaveUser = `
	INSERT INTO users (name, pass_hash)
	VALUES ($1, $2)
	RETURNING id;
	`

	QueryUpdateUser = `
	UPDATE users
	SET name   = $2,
    	pass_hash   = $3,
    	avatar_original_path = $4,
		avatar_thumbnail_path = $5
	WHERE id = $1;
	`
)
