package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"
)

func (r *postgresqlRepository) GetPostByID(ctx context.Context, id int64) (*model.Post, error) {
	query := `SELECT id, title, text, readOnly FROM posts WHERE id = $1`
	var post model.Post

	err := r.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Title, &post.Text, &post.ReadOnly)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Join(storage.ErrPostNotFound, fmt.Errorf("postID = %d", id))
		}
		return nil, err
	}
	return &post, nil
}
