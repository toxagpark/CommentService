package storage

import (
	"commentsService/internal/entity"
	"commentsService/internal/usecase"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const stringsLimit = 10

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) usecase.CommentRepository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) Save(ctx context.Context, userID string, text string) error {
	query := "INSERT INTO comments (user_id, text) VALUES ($1, $2)"
	_, err := repo.db.Exec(ctx, query, userID, text)
	return err
}

func (repo *Repository) Change(ctx context.Context, commentID string, text string) error {
	query := "UPDATE comments SET text = $1, is_edited = TRUE, created_at = CURRENT_TIMESTAMP WHERE id = $2;"
	_, err := repo.db.Exec(ctx, query, text, commentID)
	return err
}

func (repo *Repository) Get(ctx context.Context, offset string, sortFromOld bool) ([]entity.Comment, error) {
	direction := "DESC"
	if sortFromOld {
		direction = "ASC"
	}
	query := fmt.Sprintf(`
		SELECT 
			c.id, 
			c.user_id, 
			u.name,
			c.text, 
			c.created_at, 
			c.is_edited
		FROM comments c
		JOIN users u ON c.user_id = u.id
		ORDER BY c.created_at %s
		LIMIT $1 OFFSET $2`, direction)

	rows, err := repo.db.Query(ctx, query, stringsLimit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entity.Comment
	for rows.Next() {
		var c entity.Comment
		err := rows.Scan(
			&c.ID,
			&c.User.ID,
			&c.User.Name,
			&c.Text,
			&c.CreatedAt,
			&c.IsEdited,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
