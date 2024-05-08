package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserUpdateForm struct {
	Email     string `form:"email" json:"email"`
	Password  string `form:"password" json:"password"`
	FirstName string `form:"firstName" json:"firstName"`
	LastName  string `form:"lastName" json:"lastName"`
}

func NewUserUpdateForm(c echo.Context) (*UserUpdateForm, error) {
	form := new(UserUpdateForm)
	if err := c.Bind(form); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid form: "+err.Error())
	}

	if len(form.Email) > 0 && !govalidator.IsEmail(form.Email) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid email address format.")
	}

	if len(form.Password) > 0 && len(form.Password) < 6 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 6 characters")
	}
	return form, nil
}
