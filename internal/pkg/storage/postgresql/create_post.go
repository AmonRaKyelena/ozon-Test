package postgresql

import (
	"context"
	"ozon-test-project/internal/pkg/storage/model"
)

func (r *postgresqlRepository) CreatePost(ctx context.Context, post model.Post) (int64, error) {
	query := `INSERT INTO posts (title, text, readOnly) VALUES ($1, $2, $3) RETURNING id`
	var id int64

	err := r.db.QueryRowContext(ctx, query, post.Title, post.Text, post.ReadOnly).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
