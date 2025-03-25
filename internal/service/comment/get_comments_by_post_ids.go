package comment

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/handlers/model"
)

func (s *commentService) GetCommentsByPostIDs(ctx context.Context, postIDs []int64, limit, offset int32) (map[int64][]*model.CommentForPagination, error) {
	commentsMap, err := s.storage.GetCommentsByPostIDs(ctx, postIDs, limit, offset)
	if err != nil {
		return nil, err
	}

	result := map[int64][]*model.CommentForPagination{}

	for postID, comments := range commentsMap {
		resultComments := make([]*model.CommentForPagination, 0, len(comments))
		for _, rootComment := range comments {
			resultComments = append(resultComments, &model.CommentForPagination{
				ID:       rootComment.Comment.ID,
				PostID:   rootComment.Comment.PostID,
				Text:     rootComment.Comment.Text,
				ParentID: rootComment.Comment.ParentID,
				HasChild: rootComment.HasChild,
			})
		}
		result[postID] = resultComments
	}
	return result, nil
}
