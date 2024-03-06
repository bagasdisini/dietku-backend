package handler

import (
	"dietku-backend/cmd/auth/gear"
	"dietku-backend/cmd/user/repo"
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type AuthHandler struct {
	repo *repo.UserRepository
}

func NewAuthHandler(e *echo.Echo, db *mongo.Database) {
	h := &AuthHandler{
		repo: repo.NewUserRepository(db),
	}
	e.POST("/api/login", h.Login)
	e.POST("/api/register", h.Register)
}

// Login
// @Tags Auth
// @Summary Login
// @ID login
// @Router /api/login [post]
// @Accept json
// @Param body body LoginForm true "login body"
// @Produce json
// @Success 200
func (h *AuthHandler) Login(c echo.Context) error {
	form, err := NewLoginForm(c)
	if err != nil {
		return err
	}

	u, err := h.repo.FindOneByEmail(form.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Wrong username/email or password")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	if u.IsDeleted {
		return echo.NewHTTPError(http.StatusUnauthorized, "Wrong username/email or password")
	}

	if !gear.CheckPassword(u.Password, form.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Wrong username/email or password")
	}

	accessToken, err := gear.GenerateToken(u)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": accessToken,
	})
}

// Register
// @Tags Auth
// @Summary Register
// @ID register
// @Router /api/register [post]
// @Accept json
// @Param body body RegisterForm true "register body"
// @Produce json
// @Success 200
func (h *AuthHandler) Register(c echo.Context) error {
	form, err := NewRegisterForm(c)
	if err != nil {
		return err
	}

	_, err = h.repo.FindOneByEmail(form.Email)
	if err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The email provided is already taken")
	}

	u := &repo.User{
		ID:        primitive.NewObjectID(),
		Email:     form.Email,
		FullName:  form.FullName,
		Password:  gear.CryptPassword(form.Password),
		CreatedAt: time.Now(),
		IsDeleted: false,
	}

	_, err = h.repo.InsertOne(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	inserted, err := h.repo.FindOne(u.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, inserted)
}
