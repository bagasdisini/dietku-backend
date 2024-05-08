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
	FirstName string `form:"firstName" json:"firstName"`
	LastName  string `form:"lastName" json:"lastName"`
	BirthDay  string `form:"birthDay" json:"birthDay"`
	Phone     string `form:"phone" json:"phone"`
	Email     string `form:"email" json:"email"`
	Password  string `form:"password" json:"password"`
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

	form.FirstName = strings.TrimSpace(form.FirstName)
	if form.FirstName == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "FirstName is required")
	}

	form.LastName = strings.TrimSpace(form.LastName)
	if form.LastName == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "LastName is required")
	}

	form.BirthDay = strings.TrimSpace(form.BirthDay)
	if form.BirthDay == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "BirthDay is required")
	}

	form.Phone = strings.TrimSpace(form.Phone)
	if form.Phone == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Phone is required")
	}
	return form, nil
}
