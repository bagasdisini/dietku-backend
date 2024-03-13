package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type BlogForm struct {
	Header   string   `json:"header" bson:"header"`
	Content  string   `json:"content" bson:"content"`
	Category []string `json:"category" bson:"category"`
}

func NewBlogForm(c echo.Context) (*BlogForm, error) {
	form := new(BlogForm)
	if err := c.Bind(form); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid form: "+err.Error())
	}

	if len(form.Header) < 1 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Header is required.")
	}

	if len(form.Content) < 1 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Content is required.")
	}

	return form, nil
}

func NewUpdateBlogForm(c echo.Context) (*BlogForm, error) {
	form := new(BlogForm)
	if err := c.Bind(form); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid form: "+err.Error())
	}
	return form, nil
}
