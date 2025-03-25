package postgresql

import (
	"context"
	"database/sql"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"
)

func (r *postgresqlRepository) GetCommentsByParentID(ctx context.Context, parentID int64, limit, offset int32) ([]*model.CommentForPagination, error) {
	query := `
		SELECT 
			c.id,
			c.parentIdcomment,
			c.idPost,
			c.text,
			EXISTS (
				SELECT 1 FROM comments WHERE parentIdcomment = c.id
			) as has_child
		FROM comments c
		WHERE c.parentIdcomment = $1
		ORDER BY c.id
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, parentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.CommentForPagination
	for rows.Next() {
		var cfp model.CommentForPagination
		var parent sql.NullInt64
		if err := rows.Scan(&cfp.Comment.ID, &parent, &cfp.Comment.PostID, &cfp.Comment.Text, &cfp.HasChild); err != nil {
			return nil, err
		}
		if parent.Valid {
			cfp.Comment.ParentID = &parent.Int64
		} else {
			cfp.Comment.ParentID = nil
		}
		comments = append(comments, &cfp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
