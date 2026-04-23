package handler

import (
	"errors"
	"net/http"
	"strconv"

	models "github.com/OGZKTeBmj/forum/domain"
	"github.com/OGZKTeBmj/forum/internal/service"
	"github.com/OGZKTeBmj/forum/internal/storage"
	"github.com/OGZKTeBmj/forum/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) posts(ctx *gin.Context) {
	const op = "handler.posts"
	log := h.log.With("op", op)

	offset, _ := strconv.Atoi(ctx.Query("offset"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	if offset < 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "offset must be positive or zero"})
		return
	}
	if limit <= 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "limit must be positive"})
		return
	}

	posts, err := h.postsService.GetPostsBatch(ctx.Request.Context(), int64(offset), int64(limit))
	if err != nil {
		if !errors.Is(err, service.ErrNotFound) {
			log.Error("can't get posts", utils.SlogErr(err))
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "can't get posts"})
			return
		}
	}
	userIdString := ctx.GetString(userIdCtx)

	responce := PostsBatchResponce{Batch: make([]*Post, 0)}
	for _, post := range posts.Posts {
		author, err := h.authService.User(ctx, post.AuthorId)
		if err != nil {
			if !errors.Is(err, storage.ErrIsNotExist) {
				log.Error("can't get user")
			}
		}
		postResponce := Post{
			Id:      post.Id,
			Title:   post.Title,
			Content: post.Content,
			Author: User{
				Name: author.Name,
				AvatarPath: Avatar{
					Original:  author.AvatarPath.Original,
					Thumbnail: author.AvatarPath.Thumbnail,
				},
			},
			Rating:        post.Rating,
			TimeStamp:     post.TimeStamp,
			CommentsCount: post.CommentsCount,
		}
		if len(userIdString) > 0 {
			userVote, err := h.postsService.GetVote(ctx, []byte(userIdString), post.Id)
			if err != nil {
				log.Error("can't get vote", utils.SlogErr(err))
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
			postResponce.UserVote = userVote.Value
		}
		responce.Batch = append(responce.Batch, &postResponce)
	}
	ctx.IndentedJSON(http.StatusOK, responce)
}

func (h *Handler) savePost(ctx *gin.Context) {
	const op = "handler.savePost"
	log := h.log.With("op", op)

	request := &CreatePostRequest{}
	if err := ctx.BindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	authorId := ctx.GetString(userIdCtx)

	var post *models.Post = &models.Post{
		Title:    request.Title,
		Content:  request.Content,
		AuthorId: []byte(authorId),
	}

	id, err := h.postsService.SavePost(ctx.Request.Context(), post)
	if err != nil {
		log.Error("can't save post", utils.SlogErr(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.IndentedJSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) comments(ctx *gin.Context) {
	const op = "handler.comments"
	log := h.log.With("op", op)

	postIdStr := ctx.Param("post_id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	offset, _ := strconv.Atoi(ctx.Query("offset"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	if offset < 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "offset must be positive or zero"})
		return
	}
	if limit <= 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "limit must be positive"})
		return
	}

	comments, err := h.postsService.CommentsBatch(ctx.Request.Context(), int64(postId), int64(offset), int64(limit))
	if err != nil {
		if !errors.Is(err, service.ErrNotFound) {
			log.Error("can't get comments", utils.SlogErr(err))
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "can't get comments"})
			return
		}
	}

	responce := CommentsBatchResponce{Batch: make([]*Comment, 0)}
	for _, comment := range comments.Comments {
		author, err := h.authService.User(ctx, comment.AuthorId)
		if err != nil {
			if !errors.Is(err, storage.ErrIsNotExist) {
				log.Error("can't get user")
			}
		}
		commentResponce := Comment{
			Id:      comment.Id,
			Content: comment.Content,
			Author: User{
				Name: author.Name,
				AvatarPath: Avatar{
					Original:  author.AvatarPath.Original,
					Thumbnail: author.AvatarPath.Thumbnail,
				},
			},
			TimeStamp: comment.TimeStamp,
		}
		responce.Batch = append(responce.Batch, &commentResponce)
	}
	ctx.IndentedJSON(http.StatusOK, responce)
}

func (h *Handler) saveComment(ctx *gin.Context) {
	const op = "handler.saveComment"
	log := h.log.With("op", op)

	postIdStr := ctx.Param("post_id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	request := &CreateCommentRequest{}
	if err := ctx.BindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	authorId := ctx.GetString(userIdCtx)

	var comment *models.Comment = &models.Comment{
		Content:  request.Content,
		PostId:   postId,
		AuthorId: []byte(authorId),
	}

	id, err := h.postsService.SaveComment(ctx.Request.Context(), comment)
	if err != nil {
		log.Error("can't save comment", utils.SlogErr(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.IndentedJSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) vote(ctx *gin.Context) {
	const op = "handler.vote"
	log := h.log.With("op", op)

	postIdStr := ctx.Param("post_id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var request VoteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	userId := ctx.GetString(userIdCtx)

	vote := models.Vote{
		PostId:   postId,
		AuthorId: []byte(userId),
		Value:    request.Value,
	}

	err = h.postsService.Vote(ctx.Request.Context(), vote)
	if err != nil {
		log.Error("can't vote", utils.SlogErr(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) deleteVote(ctx *gin.Context) {
	const op = "handler.deleteVote"
	log := h.log.With("op", op)

	postIdStr := ctx.Param("post_id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	userId := ctx.GetString(userIdCtx)

	vote := models.Vote{PostId: postId, AuthorId: []byte(userId)}
	err = h.postsService.DeleteVote(ctx, vote)
	if err != nil {
		log.Error("can't delete vote", utils.SlogErr(err))
		ctx.Status(http.StatusInternalServerError)
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) deletePost(ctx *gin.Context) {
	const op = "handler.rating"
	log := h.log.With("op", op)

	postIdStr := ctx.Param("post_id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	err = h.postsService.DeletePost(ctx.Request.Context(), postId)
	if err != nil {
		log.Error("can't delete post", utils.SlogErr(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
