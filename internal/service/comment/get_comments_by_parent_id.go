package comment

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/handlers/model"
)

func (s *commentService) GetCommentsByParentID(ctx context.Context, parentID int64, limit, offset int32) ([]*model.CommentForPagination, error) {
	comments, err := s.storage.GetCommentsByParentID(ctx, parentID, limit, offset)
	if err != nil {
		return nil, err
	}
	result := make([]*model.CommentForPagination, 0, len(comments))
	for _, comment := range comments {
		result = append(result, &model.CommentForPagination{
			ID:       comment.Comment.ID,
			ParentID: comment.Comment.ParentID,
			Text:     comment.Comment.Text,
			PostID:   comment.Comment.PostID,
			HasChild: comment.HasChild,
		})
	}
	return result, nil
}
