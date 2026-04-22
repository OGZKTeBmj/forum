package handler

import (
	"context"
	"errors"
	"log/slog"
	"mime/multipart"

	models "github.com/OGZKTeBmj/forum/domain"
	"github.com/OGZKTeBmj/forum/utils"
	"github.com/gin-gonic/gin"
)

var (
	ErrPleaseInit = errors.New("please init handler")
)

type Handler struct {
	log          *slog.Logger
	router       *gin.Engine
	authService  AuthService
	postsService PostsService
	imageService ImageService
	initStatus   bool
}

type AuthService interface {
	SignUp(ctx context.Context, name string, password string) (id []byte, err error)
	SignIn(ctx context.Context, name string, password string) (token string, err error)
	User(ctx context.Context, id []byte) (user models.User, err error)
	UpdateUser(ctx context.Context, user models.User) error
	ParseToken(token string) (userId []byte, err error)
}

type PostsService interface {
	GetPostsBatch(ctx context.Context, offset int64, limit int64) (posts models.Posts, err error)
	SavePost(ctx context.Context, post *models.Post) (id int64, err error)
	CommentsBatch(ctx context.Context, postID int64, offset int64, limit int64) (commets models.Comments, err error)
	SaveComment(ctx context.Context, comment *models.Comment) (id int64, err error)
	Vote(ctx context.Context, vote models.Vote) error
	DeleteVote(ctx context.Context, vote models.Vote) error
	GetVote(ctx context.Context, authorId []byte, postId int64) (vote models.Vote, err error)
}

type ImageService interface {
	UploadAvatar(ctx context.Context, userId []byte, fileHeader *multipart.FileHeader) (avatarPath models.AvatarPath, err error)
}

func New(log *slog.Logger, authService AuthService, postsService PostsService, imageService ImageService) *Handler {
	return &Handler{
		log:          log,
		router:       gin.New(),
		authService:  authService,
		postsService: postsService,
		imageService: imageService,
	}
}

func (h *Handler) Init() {
	api := h.router.Group("/api")
	{
		posts := api.Group("/posts")
		{
			posts.GET("", h.userIdentityWithoutAbort, h.posts)
			posts.POST("", h.userIdentity, h.savePost)

			ac := posts.Group("/:post_id")
			{
				ac.PUT("/vote", h.userIdentity, h.vote)
				ac.DELETE("/vote", h.userIdentity, h.deleteVote)

				ac.GET("/comments", h.comments)
				ac.POST("/comments", h.userIdentity, h.saveComment)
			}
		}
		user := api.Group("/profile")
		{
			user.GET("", h.userIdentity, h.profile)
			user.POST("/avatar", h.userIdentity, h.uploadAvatar)
		}
	}

	a := h.router.Group("/auth")
	{
		a.GET("/token-valid", h.tokenValid)
		a.POST("/sign-up", h.signUp)
		a.POST("/sign-in", h.signIn)
	}
	h.initStatus = true
}

func (h *Handler) Run(addr string) error {
	const op = "handler.Run"

	if !h.initStatus {
		return utils.ErrWrap(op, ErrPleaseInit)
	}

	if err := h.router.Run(addr); err != nil {
		return utils.ErrWrap(op, err)
	}

	return nil
}
