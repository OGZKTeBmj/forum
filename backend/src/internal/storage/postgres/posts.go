package postgres

import (
	"context"
	"errors"

	models "github.com/OGZKTeBmj/forum/domain"
	"github.com/OGZKTeBmj/forum/utils"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) PostsBatch(ctx context.Context, offset int64, limit int64) (models.Posts, error) {
	const op = "postgres.GetPostsBatch"

	batch := models.Posts{Posts: make([]*models.Post, 0)}

	rows, err := s.db.Query(ctx, QueryPosts, offset, limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Posts{}, nil
		}
		return models.Posts{}, utils.ErrWrap(op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.Id, &post.Title, &post.Content,
			&post.AuthorId, &post.Rating, &post.TimeStamp, &post.CommentsCount); err != nil {
			return models.Posts{}, utils.ErrWrap(op, err)
		}
		batch.Posts = append(batch.Posts, &post)
	}
	return batch, nil
}

func (s *Storage) SavePost(ctx context.Context, post *models.Post) (int64, error) {
	const op = "postgres.SavePost"

	var id int64

	if err := s.db.QueryRow(ctx, QuerySavePost, post.Title,
		post.Content, post.AuthorId).Scan(&id); err != nil {
		return -1, utils.ErrWrap(op, err)
	}
	return id, nil
}

func (s *Storage) CommentsBatch(
	ctx context.Context, postID int64, offset int64, limit int64) (
	commets models.Comments, err error) {

	const op = "postgres.CommentsBatchByPostID"

	batch := models.Comments{Comments: make([]*models.Comment, 0)}

	rows, err := s.db.Query(ctx, QueryComments, postID, offset, limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Comments{}, nil
		}
		return models.Comments{}, utils.ErrWrap(op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.Id, &comment.PostId,
			&comment.AuthorId, &comment.Content,
			&comment.TimeStamp); err != nil {
			return models.Comments{}, utils.ErrWrap(op, err)
		}
		batch.Comments = append(batch.Comments, &comment)
	}
	return batch, nil
}

func (s *Storage) SaveComment(ctx context.Context, comment *models.Comment) (int64, error) {
	const op = "postgres.SaveComment"

	var id int64

	if err := s.db.QueryRow(ctx, QuerySaveComment, comment.PostId,
		comment.AuthorId, comment.Content).Scan(&id); err != nil {
		return -1, utils.ErrWrap(op, err)
	}
	return id, nil
}

func (s *Storage) Vote(ctx context.Context, vote models.Vote) error {
	const op = "postgres.Vote"

	if _, err := s.db.Exec(ctx, QueryVote, vote.PostId, vote.AuthorId, vote.Value); err != nil {
		return utils.ErrWrap(op, err)
	}
	return nil
}

func (s *Storage) DeleteVote(ctx context.Context, vote models.Vote) error {
	const op = "postgres.DeleteVote"

	if _, err := s.db.Exec(ctx, QueryDeleteVote, vote.PostId, vote.AuthorId); err != nil {
		return utils.ErrWrap(op, err)
	}
	return nil
}

func (s *Storage) GetVote(ctx context.Context, authorId []byte, postId int64) (vote models.Vote, err error) {
	const op = "postgres.GetVote"

	if err = s.db.QueryRow(ctx, QueryGetVote, postId, authorId).Scan(&vote.PostId, &vote.AuthorId, &vote.Value); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return models.Vote{}, utils.ErrWrap(op, err)
		}
		return models.Vote{}, nil
	}
	return vote, nil
}

const (
	QuerySavePost = `
	INSERT INTO posts (title, content, author_id, time_stamp)
	VALUES ($1, $2, $3, NOW())
	RETURNING id;
	`

	QueryPosts = `
	SELECT p.id, p.title, p.content, p.author_id,
     COALESCE(v.rating, 0) AS rating,
     p.time_stamp,
     COALESCE(c.comments_count, 0) AS comments_count
	FROM posts p
	LEFT JOIN (
     SELECT post_id, SUM(value) AS rating
     FROM votes
     GROUP BY post_id
	) v ON v.post_id = p.id
	LEFT JOIN (
     SELECT post_id, COUNT(*) AS comments_count
     FROM comments
     GROUP BY post_id
	) c ON c.post_id = p.id
	ORDER BY p.id DESC
	LIMIT $2 OFFSET $1;
	`

	QueryComments = `
	SELECT id, post_id, author_id, content, time_stamp
	FROM comments
	WHERE post_id = $1
	ORDER BY id DESC
	LIMIT $3 OFFSET $2
	`

	QuerySaveComment = `
	INSERT INTO comments (post_id, author_id, content, time_stamp)
	VALUES ($1, $2, $3, NOW())
	RETURNING id
	`

	QueryVote = `
	INSERT INTO votes (post_id, author_id, value)
	VALUES ($1, $2, $3)
	ON CONFLICT (post_id, author_id)
	DO UPDATE SET value = EXCLUDED.value;
	`
	QueryDeleteVote = `
	DELETE FROM votes
    WHERE post_id = $1 AND author_id = $2
	`

	QueryGetVote = `
	SELECT post_id, author_id, value
	FROM votes
	WHERE post_id = $1 AND author_id = $2
	`
)
