package service

import (
	"context"
	"log/slog"

	models "github.com/OGZKTeBmj/forum/domain"
	"github.com/OGZKTeBmj/forum/utils"
)

type PostsService struct {
	Provider PostProvider
	Log      *slog.Logger
}

type PostProvider interface {
	PostsBatch(ctx context.Context, offset int64, limit int64) (posts models.Posts, err error)
	SavePost(ctx context.Context, post *models.Post) (id int64, err error)
	CommentsBatch(ctx context.Context, postID int64, offset int64, limit int64) (commets models.Comments, err error)
	SaveComment(ctx context.Context, comment *models.Comment) (id int64, err error)
	Vote(ctx context.Context, vote models.Vote) error
	DeleteVote(ctx context.Context, vote models.Vote) error
	GetVote(ctx context.Context, authorId []byte, postId int64) (vote models.Vote, err error)
	DeletePost(ctx context.Context, postId int64) error
}

func (p *PostsService) GetPostsBatch(ctx context.Context, offset int64, limit int64) (models.Posts, error) {
	const op = "postsService.GetPosts"
	log := p.Log.With("op", op)

	posts, err := p.Provider.PostsBatch(ctx, offset, limit)
	if err != nil {
		log.Error("Failed to get posts")
		return models.Posts{}, utils.ErrWrap(op, err)
	}
	if posts.Posts == nil {
		return models.Posts{}, utils.ErrWrap(op, ErrNotFound)
	}
	return posts, nil
}

func (p *PostsService) SavePost(ctx context.Context, post *models.Post) (int64, error) {
	const op = "postsService.SavePost"

	log := p.Log.With("op", op)

	id, err := p.Provider.SavePost(ctx, post)
	if err != nil {
		log.Error("Failed save post", utils.SlogErr(err))
		return -1, utils.ErrWrap(op, err)
	}
	return id, nil
}

func (p *PostsService) CommentsBatch(
	ctx context.Context, postID int64, offset int64, limit int64) (
	commets models.Comments, err error) {

	const op = "postService.CommentsBatch"

	log := p.Log.With("op", op)
	comments, err := p.Provider.CommentsBatch(ctx, postID, offset, limit)
	if err != nil {
		log.Error("Failed to get comments")
		return models.Comments{}, utils.ErrWrap(op, err)
	}
	if comments.Comments == nil {
		return models.Comments{}, utils.ErrWrap(op, ErrNotFound)
	}
	return comments, nil
}

func (p *PostsService) SaveComment(ctx context.Context, comment *models.Comment) (int64, error) {
	const op = "postsService.SaveComment"

	log := p.Log.With("op", op)

	id, err := p.Provider.SaveComment(ctx, comment)
	if err != nil {
		log.Error("Failed save comment", utils.SlogErr(err))
		return -1, utils.ErrWrap(op, err)
	}
	return id, nil
}

func (p *PostsService) Vote(ctx context.Context, vote models.Vote) error {
	const op = "postsService.Vote"

	log := p.Log.With("op", op)

	err := p.Provider.Vote(ctx, vote)
	if err != nil {
		log.Error("Failed save vote", utils.SlogErr(err))
		return utils.ErrWrap(op, err)
	}
	return nil
}

func (p *PostsService) DeleteVote(ctx context.Context, vote models.Vote) error {
	const op = "postsService.DeleteVote"

	log := p.Log.With("op", op)

	err := p.Provider.DeleteVote(ctx, vote)
	if err != nil {
		log.Error("Failed delete vote", utils.SlogErr(err))
		return utils.ErrWrap(op, err)
	}
	return nil
}

func (p *PostsService) GetVote(ctx context.Context, authorId []byte, postId int64) (models.Vote, error) {
	const op = "postsService.UserIsVote"

	log := p.Log.With("op", op)

	vote, err := p.Provider.GetVote(ctx, authorId, postId)
	if err != nil {
		log.Error("Failed check user vote", utils.SlogErr(err))
		return vote, utils.ErrWrap(op, err)
	}
	return vote, nil
}

func (p *PostsService) DeletePost(ctx context.Context, postId int64) error {
	const op = "postsService.DeletePost"

	log := p.Log.With("op", op)

	err := p.Provider.DeletePost(ctx, postId)
	if err != nil {
		log.Error("Failed delete post", utils.SlogErr(err))
		return utils.ErrWrap(op, err)
	}
	return err
}
