package handler

import (
	"errors"
	"net/http"

	"github.com/OGZKTeBmj/forum/internal/service"
	"github.com/OGZKTeBmj/forum/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) uploadAvatar(ctx *gin.Context) {
	const op = "handler.uploadAvatar"
	log := h.log.With("op", op)

	userIDStr := ctx.GetString(userIdCtx)
	userID := []byte(userIDStr)

	fileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}

	avatarPath, err := h.imageService.UploadAvatar(ctx.Request.Context(), userID, fileHeader)
	if err != nil {
		log.Error("can't upload avatar", utils.SlogErr(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	user, err := h.authService.User(ctx.Request.Context(), userID)
	if err != nil {
		log.Error("can't get user", utils.SlogErr(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	user.AvatarPath = avatarPath
	err = h.authService.UpdateUser(ctx.Request.Context(), user)
	if err != nil {
		log.Error("cannot update user avatar", utils.SlogErr(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"avatar": Avatar{
		Original:  avatarPath.Original,
		Thumbnail: avatarPath.Thumbnail,
	}})
}

func (h *Handler) profile(ctx *gin.Context) {
	const op = "handler.profile"
	log := h.log.With("op", op)

	userIDStr := ctx.GetString(userIdCtx)
	userID := []byte(userIDStr)

	user, err := h.authService.User(ctx, userID)
	if err != nil {
		if !errors.Is(err, service.ErrUserIsNotExist) {
			log.Error("can't get user")
			ctx.Status(http.StatusInternalServerError)
		}
		ctx.Status(http.StatusUnauthorized)
	}

	ctx.JSON(http.StatusOK, User{
		Name: user.Name,
		AvatarPath: Avatar{
			Original:  user.AvatarPath.Original,
			Thumbnail: user.AvatarPath.Thumbnail,
		},
	})
}
