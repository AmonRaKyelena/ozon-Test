package comment

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/handlers/model"
)

func (s *commentService) GetComments(ctx context.Context, postID int64, parentID *int64, limit, offset int32) ([]*model.CommentForPagination, error) {
	if parentID == nil {
		comments, err := s.GetCommentsByPostIDs(ctx, []int64{postID}, limit, offset)
		if err != nil {
			return nil, err
		}
		return comments[postID], nil
	}

	return s.GetCommentsByParentID(ctx, *parentID, limit, offset)
}
