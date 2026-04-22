package handler

import (
	"errors"
	"net/http"

	"github.com/OGZKTeBmj/forum/internal/service"
	"github.com/OGZKTeBmj/forum/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(ctx *gin.Context) {
	const op = "handler.signUp"
	log := h.log.With("op", op)

	var request UserInput
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	if (len(request.Name) < 3) || (len(request.Password) < 8) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "username lenght < 3 or password lenght < 8"})
		return
	}

	id, err := h.authService.SignUp(ctx, request.Name, request.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserIsExist) {
			ctx.JSON(http.StatusConflict, gin.H{"message": "user is exist"})
			return
		}
		log.Error("failed register user", utils.SlogErr(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed register user"})
		return
	}
	ctx.JSON(http.StatusCreated, id)
}

func (h *Handler) signIn(ctx *gin.Context) {
	const op = "handler.signIn"
	log := h.log.With("op", op)

	var request UserInput
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	if (len(request.Name) < 3) || (len(request.Password) < 8) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "username lenght < 3 or password lenght < 8"})
		return
	}

	token, err := h.authService.SignIn(ctx, request.Name, request.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentails) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentails"})
			return
		}
		log.Error("failed login user: ", utils.SlogErr(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed login user"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) tokenValid(ctx *gin.Context) {
	userId := h._userIdentityWithoutAbort(ctx)
	if userId != nil {
		ctx.IndentedJSON(http.StatusOK, gin.H{"valid": true})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"valid": false})
}
