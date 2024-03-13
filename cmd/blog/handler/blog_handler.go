package handler

import (
	"dietku-backend/cmd/auth/gear"
	"dietku-backend/cmd/blog/repo"
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type BlogHandler struct {
	repo *repo.BlogRepository
}

func NewBlogApi(e *echo.Echo, db *mongo.Database) *BlogHandler {
	b := &BlogHandler{
		repo: repo.NewBlogRepository(db),
	}
	bGroup := e.Group("")
	{
		bGroup.GET("/api/blog", b.Blogs)
		bGroup.GET("/api/blog/:id", b.Blog)
		bGroup.GET("/api/blog/user/:userId", b.BlogsByUser)
		bGroup.GET("/api/blog/category/:category", b.BlogsByCategory)

		bGroup.POST("/api/blog", b.Create, gear.IsLoggedIn(db))

		bGroup.PUT("/api/blog/:id", b.Update, gear.IsLoggedIn(db))

		bGroup.DELETE("/api/blog/:id", b.Delete, gear.IsLoggedIn(db))
	}
	return b
}

// Blogs
// @Tags Blog
// @Summary Get All Blogs
// @ID blog
// @Router /api/blog [get]
// @Produce json
// @Success 200
func (h *BlogHandler) Blogs(c echo.Context) error {
	blogs, err := h.repo.FindAll()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusBadRequest, "Blogs not found!", c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while getting blog.", c)
	}
	return c.JSON(http.StatusOK, blogs)
}

// Blog
// @Tags Blog
// @Summary Get Blog
// @ID blog-get
// @Router /api/blog/{id} [get]
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200
func (h *BlogHandler) Blog(c echo.Context) error {
	id := c.Param("id")
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid blog id", c)
	}

	blog, err := h.repo.FindOne(oId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusBadRequest, "Blog not found!", c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while getting blog.", c)
	}
	return c.JSON(http.StatusOK, blog)
}

// BlogsByUser
// @Tags Blog
// @Summary Get Blogs By User
// @ID blog-user
// @Router /api/blog/user/{userId} [get]
// @Produce json
// @Param userId path string true "User ID"
// @Success 200
func (h *BlogHandler) BlogsByUser(c echo.Context) error {
	userId := c.Param("userId")
	oId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id", c)
	}

	blogs, err := h.repo.FindByUser(oId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusBadRequest, "Blogs not found!", c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while getting blog.", c)
	}
	return c.JSON(http.StatusOK, blogs)
}

// BlogsByCategory
// @Tags Blog
// @Summary Get Blogs By Category
// @ID blog-category
// @Router /api/blog/category/{category} [get]
// @Produce json
// @Param category path string true "Category"
// @Success 200
func (h *BlogHandler) BlogsByCategory(c echo.Context) error {
	category := c.Param("category")

	blogs, err := h.repo.FindByCategory(category)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusBadRequest, "Blogs not found!", c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while getting blog.", c)
	}
	return c.JSON(http.StatusOK, blogs)
}

// Create
// @Tags Blog
// @Summary Create Blog
// @ID blog-create
// @Router /api/blog [post]
// @Accept json
// @Param body body BlogForm true "blog body"
// @Produce json
// @Success 200
// @Security ApiKeyAuth
func (h *BlogHandler) Create(c echo.Context) error {
	form, err := NewBlogForm(c)
	if err != nil {
		return err
	}

	tokenData := c.Get("me").(*gear.UserClaims)

	b := &repo.Blog{
		ID:       primitive.NewObjectID(),
		Header:   form.Header,
		Content:  form.Content,
		Category: form.Category,
		CreatedBy: repo.By{
			ID:       tokenData.ID,
			Email:    tokenData.Email,
			FullName: tokenData.FullName,
			At:       time.Now(),
		},
	}

	_, err = h.repo.InsertOne(b)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while creating blog.", c)
	}

	docs, err := h.repo.FindOne(b.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while getting blog.", c)
	}
	return c.JSON(http.StatusOK, docs)
}

// Update
// @Tags Blog
// @Summary Update Blog
// @ID blog-update
// @Router /api/blog/{id} [put]
// @Accept json
// @Param id path string true "Blog ID"
// @Param body body BlogForm true "blog body"
// @Produce json
// @Success 200
// @Security ApiKeyAuth
func (h *BlogHandler) Update(c echo.Context) error {
	id := c.Param("id")
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid blog id", c)
	}

	blog, err := h.repo.FindOne(oId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusBadRequest, "Blog not found!", c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while getting blog.", c)
	}

	tokenData := c.Get("me").(*gear.UserClaims)

	if tokenData.ID != blog.CreatedBy.ID {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized to update this blog", c)
	}

	form, err := NewUpdateBlogForm(c)
	if err != nil {
		return err
	}

	if form.Header == "" && form.Content == "" && len(form.Category) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Nothing to update", c)
	}

	if form.Header != "" {
		blog.Header = form.Header
	}
	if form.Content != "" {
		blog.Content = form.Content
	}
	if len(form.Category) > 0 {
		blog.Category = form.Category
	}

	blog.UpdatedBy = &repo.By{
		ID:       tokenData.ID,
		Email:    tokenData.Email,
		FullName: tokenData.FullName,
	}

	docs, err := h.repo.UpdateOne(blog)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while updating blog.", c)
	}
	return c.JSON(http.StatusOK, docs)
}

// Delete
// @Tags Blog
// @Summary Delete Blog
// @ID blog-delete
// @Router /api/blog/{id} [delete]
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200
// @Security ApiKeyAuth
func (h *BlogHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid blog id", c)
	}

	blog, err := h.repo.FindOne(oId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusBadRequest, "Blog not found!", c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while getting blog.", c)
	}

	tokenData := c.Get("me").(*gear.UserClaims)

	if tokenData.ID != blog.CreatedBy.ID {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized to delete this blog", c)
	}

	docs, err := h.repo.DeleteOne(blog.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while deleting blog.", c)
	}
	return c.JSON(http.StatusOK, docs)
}
