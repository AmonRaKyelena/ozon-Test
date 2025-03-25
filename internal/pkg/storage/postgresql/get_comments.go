package postgresql

import (
	"context"
	"database/sql"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"

	"github.com/lib/pq"
)

func (r *postgresqlRepository) GetCommentsByPostIDs(
	ctx context.Context,
	postIDs []int64,
	limit, offset int32,
) (map[int64][]model.CommentForPagination, error) {
	query := `
		WITH numbered AS (
			SELECT
				c.id,
				c.parentIdcomment,
				c.idPost,
				c.text,
				EXISTS (
					SELECT 1 FROM comments cc WHERE cc.parentIdcomment = c.id
				) AS has_child,
				row_number() OVER (PARTITION BY c.idPost ORDER BY c.id) AS rn
			FROM comments c
			WHERE c.idPost = ANY ($1)
		)
		SELECT id, parentIdcomment, idPost, text, has_child
		FROM numbered
		WHERE rn > $2 AND rn <= ($2 + $3)
		ORDER BY idPost, rn;
	`
	rows, err := r.db.QueryContext(ctx, query, pq.Array(postIDs), offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]model.CommentForPagination)
	for rows.Next() {
		var (
			id       int64
			parent   sql.NullInt64
			postID   int64
			text     string
			hasChild bool
		)
		if err := rows.Scan(&id, &parent, &postID, &text, &hasChild); err != nil {
			return nil, err
		}

		var parentID *int64
		if parent.Valid {
			parentID = &parent.Int64
		}

		comment := model.Comment{
			ID:       id,
			ParentID: parentID,
			PostID:   postID,
			Text:     text,
		}

		cfp := model.CommentForPagination{
			Comment:  comment,
			HasChild: hasChild,
		}

		result[postID] = append(result[postID], cfp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
