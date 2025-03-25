package comment

import (
	"context"
	"errors"
	storage_model "ozon-test-project/internal/pkg/storage/model"
)

func (s *commentService) CreateComment(ctx context.Context, postID int64, text string, parentID *int64) (int64, error) {
	post, err := s.storage.GetPostByID(ctx, postID)
	if err != nil {
		return 0, err
	}

	if post.ReadOnly {
		return 0, errors.New("can't create comment on readOnly post")
	}

	commentID, err := s.storage.CreateComment(ctx, storage_model.Comment{
		PostID:   postID,
		Text:     text,
		ParentID: parentID,
	})
	if err != nil {
		return 0, err
	}

	return commentID, nil
}
