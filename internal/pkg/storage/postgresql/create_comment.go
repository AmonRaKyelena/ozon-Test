package postgresql

import (
	"context"
	"ozon-test-project/internal/pkg/storage/model"
)

func (r *postgresqlRepository) CreateComment(ctx context.Context, comment model.Comment) (int64, error) {
	query := `INSERT INTO comments (idPost, parentIdcomment, text) VALUES ($1, $2, $3) RETURNING id`
	var id int64

	err := r.db.QueryRowContext(ctx, query, comment.PostID, comment.ParentID, comment.Text).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
