package comment

import (
	"commentsService/internal/entity"
	"commentsService/internal/usecase"
	"context"
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrSaveComment   = errors.New("save comment error")
	ErrChangeComment = errors.New("change comment error")
	ErrGetComment    = errors.New("get comment error")
)

type UseCase struct {
	commentRepo usecase.CommentRepository
}

func NewUseCase(commentRepo usecase.CommentRepository) *UseCase {
	return &UseCase{
		commentRepo: commentRepo,
	}
}

func (uc *UseCase) SaveComment(ctx context.Context, userID string, text string) error {
	err := uc.commentRepo.Save(ctx, userID, text)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSaveComment, err)
	}

	return nil
}

func (uc *UseCase) ChangeComment(ctx context.Context, commentID string, text string) error {
	if err := uc.commentRepo.Change(ctx, commentID, text); err != nil {
		return fmt.Errorf("%w: %w", ErrChangeComment, err)
	}
	return nil
}

func (uc *UseCase) Get(ctx context.Context, offsetStr string, sortFromOld bool) (*entity.PaginationInfo, error) {
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offsetStr = "0"
		offset = 0
	}

	hasPrev := true
	hasNext := true
	if offset < 10 {
		hasPrev = false
	}

	comments, err := uc.commentRepo.Get(ctx, offsetStr, sortFromOld)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetComment, err)
	}
	if len(comments) < 10 {
		hasNext = false
	}
	return &entity.PaginationInfo{
		Comments:       comments,
		HasPrev:        hasPrev,
		HasNext:        hasNext,
		PrevOffset:     strconv.Itoa(offset - 10),
		NextOffset:     strconv.Itoa(offset + 10),
		SortFromOldest: sortFromOld,
	}, nil
}
