package usecase

import (
	"commentsService/internal/entity"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=CommentRepository
type CommentRepository interface {
	Save(ctx context.Context, userID, text string) error
	Change(ctx context.Context, commentID string, text string) error
	Get(ctx context.Context, offset string, sortFromOld bool) ([]entity.Comment, error)
}
