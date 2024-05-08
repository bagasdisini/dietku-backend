package handler

import (
	"context"
	"dietku-backend/cmd/auth/gear"
	"dietku-backend/cmd/user/repo"
	"dietku-backend/config"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"time"
)

type AuthHandler struct {
	repo *repo.UserRepository
	conf *config.Config
}

func NewAuthHandler(e *echo.Echo, db *mongo.Database, conf *config.Config) {
	h := &AuthHandler{
		repo: repo.NewUserRepository(db),
		conf: conf,
	}
	e.POST("/api/login", h.Login)
	e.POST("/api/register", h.Register)

	e.GET("/api/login-google", h.loginGoogle)
	e.GET("/api/callback-google", h.callbackGoogle)
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
		FirstName: form.FirstName,
		LastName:  form.LastName,
		BirthDay:  form.BirthDay,
		Phone:     form.Phone,
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

func (h *AuthHandler) loginGoogle(c echo.Context) error {
	var oauthConfGl = &oauth2.Config{
		ClientID:     h.conf.GoogleClientID,
		ClientSecret: h.conf.GoogleClientSecret,
		RedirectURL:  "https://" + h.conf.SwaggerHost + "/api/callback-google",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	URL := oauthConfGl.AuthCodeURL(h.conf.StateString)

	return c.Redirect(http.StatusTemporaryRedirect, URL)
}

func (h *AuthHandler) callbackGoogle(c echo.Context) error {
	var oauthConfGl = &oauth2.Config{
		ClientID:     h.conf.GoogleClientID,
		ClientSecret: h.conf.GoogleClientSecret,
		RedirectURL:  "https://" + h.conf.SwaggerHost + "/api/callback-google",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	state := c.FormValue("state")
	if state != h.conf.StateString {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid state")
	}

	code := c.FormValue("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Code not found.")
	}

	token, err := oauthConfGl.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Code-Token Exchange Failed")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "User Data Fetch Failed")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "JSON Parsing Failed")
	}

	var userDoc map[string]interface{}
	if err := json.Unmarshal(userData, &userDoc); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "JSON Unmarshal Failed")
	}

	email, ok := userDoc["email"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Email not found")
	}

	u, err := h.repo.FindOneByEmail(email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server exception: "+err.Error()).SetInternal(err)
	}

	if u == nil {
		newID := primitive.NewObjectID()
		user := &repo.User{
			ID:        newID,
			FirstName: userDoc["name"].(string),
			Email:     email,
			CreatedAt: time.Now(),
			IsDeleted: false,
		}

		_, err = h.repo.InsertOne(user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server exception: "+err.Error()).SetInternal(err)
		}

		u = user
	}

	accessToken, err := gear.GenerateToken(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server exception: "+err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      accessToken,
		"userDetail": string(userData),
	})
}
