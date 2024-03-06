package handler

import (
	"dietku-backend/cmd/auth/gear"
	"dietku-backend/cmd/user/repo"
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type UserHandler struct {
	repo *repo.UserRepository
}

func NewUserApi(e *echo.Echo, db *mongo.Database) *UserHandler {
	me := &UserHandler{
		repo: repo.NewUserRepository(db),
	}
	meGroup := e.Group("")
	meGroup.Use(gear.IsLoggedIn(db))
	{
		meGroup.GET("/api/user", me.Me)

		meGroup.PUT("/api/user", me.UpdateMe)
	}
	return me
}

// Me
// @Tags User
// @Summary User
// @ID user
// @Router /api/user [get]
// @Produce json
// @Success 200
// @Security ApiKeyAuth
func (h *UserHandler) Me(c echo.Context) error {
	tokenData := c.Get("me").(*gear.UserClaims)
	meData, err := h.repo.FindOne(tokenData.ID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusBadRequest, "User not found!", c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while getting user.", c)
	}
	return c.JSON(http.StatusOK, meData)
}

// UpdateMe godoc
// @Tags User
// @Summary Update me
// @ID user-update
// @Router /api/user [put]
// @Param body body UserUpdateForm true "update body"
// @Accept  json
// @Produce  json
// @Success 200
// @Security ApiKeyAuth
func (h *UserHandler) UpdateMe(c echo.Context) error {
	updateParam, err := NewUserUpdateForm(c)
	if err != nil {
		return err
	}
	tokenData := c.Get("me").(*gear.UserClaims)
	meData, err := h.repo.FindOne(tokenData.ID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusBadRequest, "User not found.", c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while getting user.", c)
	}
	if updateParam.Email != "" {
		meData.Email = updateParam.Email
	}
	if updateParam.FullName != "" {
		meData.FullName = updateParam.FullName
	}
	if updateParam.Password != "" {
		meData.Password = updateParam.Password
	}

	checkEmail, err := h.repo.FindOneByEmail(meData.Email)
	if err == nil && checkEmail.ID != meData.ID {
		return echo.NewHTTPError(http.StatusBadRequest, "The email provided is already taken.", c)
	}

	result, err := h.repo.UpdateOne(meData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while updating user.", c)
	}
	return c.JSON(http.StatusOK, result)
}
