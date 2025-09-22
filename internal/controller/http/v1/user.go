package v1

import (
	"github.com/gin-gonic/gin"
	"test_go/internal/utils"

	"net/http"
	"sync"
	"test_go/internal/controller/http/errors"
	"test_go/internal/controller/http/middleware"
	"test_go/internal/controller/http/v1/request"
	"test_go/internal/usecase"
	httpError "test_go/pkg/httpserver"
	"test_go/pkg/jwt"
	"test_go/pkg/logger"
)

type userRoutes struct {
	l  logger.Interface
	uc usecase.User
}

func NewUserRoutes(privateGroup *gin.RouterGroup, l logger.Interface, uc usecase.User, jwtManager *jwt.JWTManager) {
	r := &userRoutes{l, uc}
	{
		h := privateGroup.Group("/users")
		h.GET("/profile", r.getProfile)
		h.PATCH("/:id/change-password", r.changePassword)
		h.PUT("/:id/photo", r.setProfilePhoto)
		h.PATCH("/:id/rating", r.updateRating)
		h.PATCH("/:id/concurrent-test", r.simulateConcurrentUpdates)
	}
}

func (r *userRoutes) getProfile(c *gin.Context) {

	currentUser, err := middleware.GetCurrentUser(c)
	if err != nil {
		errors.ErrorResponse(c, err)
		return
	}

	res, err := r.uc.GetUser(c.Request.Context(), currentUser.ID)
	if err != nil {
		r.l.Error(err, "http - v1 - getProfile")
		errors.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *userRoutes) changePassword(c *gin.Context) {
	var req request.ChangePasswordRequest

	id, err := utils.ParsePathParam(utils.ParseParams{Context: c, Key: "id"}, utils.ParseInt64)
	if err != nil {
		r.l.Error(err, "http - v1 - changePassword")
		errors.ErrorResponse(c, httpError.NewBadPathParamsError(err))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - changePassword")
		errors.ErrorResponse(c, httpError.NewBadRequestBodyError(err))
		return
	}
	inp := req.ToEntity()
	inp.ID = id
	if err := r.uc.ChangePassword(c.Request.Context(), inp); err != nil {
		r.l.Error(err, "http - v1 - changePassword")
		errors.ErrorResponse(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *userRoutes) setProfilePhoto(c *gin.Context) {
	var req request.UploadFileRequest
	if err := c.ShouldBind(&req); err != nil {
		r.l.Error(err, "http - v1 - setProfilePhoto")
		errors.ErrorResponse(c, httpError.NewBadRequestBodyError(err))
		return
	}

	id, err := utils.ParsePathParam(utils.ParseParams{Context: c, Key: "id"}, utils.ParseInt64)
	if err != nil {
		r.l.Error(err, "http - v1 - setProfilePhoto")
		errors.ErrorResponse(c, httpError.NewBadPathParamsError(err))
		return
	}

	if err := r.uc.SetProfilePhoto(c.Request.Context(), id, req.File); err != nil {
		r.l.Error(err, "http - v1 - setProfilePhoto")
		errors.ErrorResponse(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *userRoutes) updateRating(c *gin.Context) {
	var req request.UpdateRatingRequest
	if err := c.ShouldBind(&req); err != nil {
		r.l.Error(err, "http - v1 - updateRating")
		errors.ErrorResponse(c, httpError.NewBadRequestBodyError(err))
		return
	}

	id, err := utils.ParsePathParam(utils.ParseParams{Context: c, Key: "id"}, utils.ParseInt64)
	if err != nil {
		r.l.Error(err, "http - v1 - updateRating")
		errors.ErrorResponse(c, httpError.NewBadPathParamsError(err))
		return
	}

	if err := r.uc.UpdateRating(c.Request.Context(), id, req.Rating); err != nil {
		r.l.Error(err, "http - v1 - updateRating")
		errors.ErrorResponse(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *userRoutes) simulateConcurrentUpdates(c *gin.Context) {
	id, err := utils.ParsePathParam(utils.ParseParams{Context: c, Key: "id"}, utils.ParseInt64)
	if err != nil {
		r.l.Error(err, "http - v1 - simulateConcurrentUpdates")
		errors.ErrorResponse(c, httpError.NewBadPathParamsError(err))
		return
	}

	numGoroutines := 2
	deltas := []float32{15.0, 30.0}

	var wg sync.WaitGroup
	ctx := c.Request.Context()

	for i := 0; i < numGoroutines; i++ {
		delta := deltas[i]
		wg.Add(1)
		go func(d float32) {
			defer wg.Done()
			if err := r.uc.UpdateRating(ctx, id, d); err != nil {
				r.l.Error(err, "http - v1 - simulateConcurrentUpdates - goroutine error")
			}
		}(delta)
	}
	wg.Wait()

	user, err := r.uc.GetUser(ctx, id)
	if err != nil {
		r.l.Error(err, "http - v1 - simulateConcurrentUpdates - get user after update")
		errors.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
