package comment

import (
	"commentsService/internal/usecase/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommentService_SaveComment_Validation(t *testing.T) {
	testCases := []struct {
		name        string
		userID      string
		text        string
		expectedErr error
		shouldCall  bool
	}{
		{
			name:        "valid comment",
			userID:      "user-123",
			text:        "Нормальный комментарий",
			expectedErr: nil,
			shouldCall:  true,
		},
		{
			name:        "empty text",
			userID:      "user-123",
			text:        "",
			expectedErr: ErrSaveComment,
			shouldCall:  false,
		},
		{
			name:        "text too long",
			userID:      "user-123",
			text:        string(make([]byte, 1001)),
			expectedErr: ErrSaveComment,
			shouldCall:  false,
		},
		{
			name:        "empty user ID",
			userID:      "",
			text:        "Текст комментария",
			expectedErr: nil,
			shouldCall:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mocks.CommentRepository{}
			service := NewUseCase(mockRepo)

			if tc.shouldCall {
				mockRepo.On("Save", mock.Anything, tc.userID, tc.text).Return(tc.expectedErr)
			}

			err := service.SaveComment(context.Background(), tc.userID, tc.text)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			if tc.shouldCall {
				mockRepo.AssertCalled(t, "Save", mock.Anything, tc.userID, tc.text)
			} else {
				mockRepo.AssertNotCalled(t, "Save")
			}
		})
	}
}
