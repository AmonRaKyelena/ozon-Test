package postgresql

import (
	"context"
	"ozon-test-project/internal/pkg/storage/model"
)

func (r *postgresqlRepository) GetAllPosts(ctx context.Context, limit, offset int32) ([]model.Post, error) {
	query := `SELECT id, title, text, readOnly FROM posts ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Text, &post.ReadOnly); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
