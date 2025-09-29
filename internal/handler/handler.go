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
