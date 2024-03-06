package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type LoginForm struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func NewLoginForm(c echo.Context) (*LoginForm, error) {
	form := new(LoginForm)
	if err := c.Bind(form); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid form: "+err.Error())
	}

	form.Email = strings.TrimSpace(form.Email)
	if !govalidator.IsEmail(form.Email) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid email format")
	}

	if form.Password == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Password is required")
	}
	return form, nil
}

type RegisterForm struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	FullName string `json:"fullname" bson:"fullname"`
}

func NewRegisterForm(c echo.Context) (*RegisterForm, error) {
	form := new(RegisterForm)
	if err := c.Bind(form); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid form: "+err.Error())
	}

	form.Email = strings.TrimSpace(form.Email)
	if !govalidator.IsEmail(form.Email) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid email format")
	}

	if len(form.Password) < 6 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 6 characters")
	}

	form.FullName = strings.TrimSpace(form.FullName)
	if form.FullName == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Fullname is required")
	}
	return form, nil
}
