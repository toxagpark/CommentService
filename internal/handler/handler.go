package handler

import (
	"commentsService/internal/usecase/comment"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	useCase *comment.UseCase
}

func NewCommentHandler(useCase *comment.UseCase) *CommentHandler {
	return &CommentHandler{
		useCase: useCase,
	}
}

// GetComments retrieves a paginated list of comments and renders them as HTML.
// @Summary Get comments (HTML page)
// @Tags comments
// @Produce html
// @Param offset query string false "Offset for pagination" default(0)
// @Param sort query string false "Sort order: 'sortFromOldest' or 'sortFromNew'" default(sortFromNew) Enums(sortFromNew,sortFromOldest)
// @Success 200 {string} string "HTML page with comments"
// @Failure 500 {string} string "Failed to load comments"
// @Router /comments [get]
func (h *CommentHandler) GetComments(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sortFromOldest := c.DefaultQuery("sort", "sortFromNew") == "sortFromOldest"

	paginated, err := h.useCase.Get(c.Request.Context(), offsetStr, sortFromOldest)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to load comments")
		return
	}

	c.HTML(http.StatusOK, "comments.html", paginated)
}

// SaveComment creates a new comment from form data and redirects to /comments.
// @Summary Create a new comment (via form)
// @Tags comments
// @Accept x-www-form-urlencoded
// @Produce html
// @Param text formData string true "Comment text"
// @Success 303 {string} string "Redirect to /comments"
// @Failure 400 {string} string "Текст комментария обязателен"
// @Failure 500 {string} string "Не удалось сохранить комментарий"
// @Router /comments [post]
func (h *CommentHandler) SaveComment(c *gin.Context) {
	text := c.PostForm("text")
	if text == "" {
		c.String(http.StatusBadRequest, "Текст комментария обязателен")
		return
	}
	userID := "4" //заглушка

	err := h.useCase.SaveComment(c.Request.Context(), userID, text)
	if errors.Is(err, comment.ErrSaveComment) {
		c.String(http.StatusInternalServerError, "Не удалось сохранить комментарий")
		return
	}

	redirectURL := "/comments"
	c.Redirect(http.StatusSeeOther, redirectURL)
}

// ChangeComment updates an existing comment via form submission and redirects to /comments.
// @Summary Update a comment
// @Tags comments
// @Accept x-www-form-urlencoded
// @Produce html
// @Param commentID formData string true "Comment ID to update"
// @Param text formData string true "New comment text"
// @Success 303 {string} string "Redirect to /comments"
// @Failure 400 {string} string "Текст или ID комментария обязателен"
// @Failure 404 {string} string "Комментарий не найден"
// @Failure 500 {string} string "Не удалось обновить комментарий"
// @Router /comments [put]
func (h *CommentHandler) ChangeComment(c *gin.Context) {
	text := c.PostForm("text")
	if text == "" {
		c.String(http.StatusBadRequest, "Текст комментария обязателен")
		return
	}

	commentID := c.PostForm("commentID")
	if commentID == "" {
		c.String(http.StatusBadRequest, "ID комментария обязателен")
		return
	}

	err := h.useCase.ChangeComment(c.Request.Context(), commentID, text)
	if errors.Is(err, comment.ErrChangeComment) {
		c.String(http.StatusInternalServerError, "Не удалось обновить комментарий")
		return
	}

	redirectURL := "/comments"
	c.Redirect(http.StatusSeeOther, redirectURL)
}
