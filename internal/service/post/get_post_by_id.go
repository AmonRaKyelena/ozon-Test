package post

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/handlers/model"
)

func (s *postService) GetPostByID(ctx context.Context, id int64) (*model.PostForPagination, error) {
	post, err := s.storage.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.PostForPagination{
		ID:       post.ID,
		Title:    post.Title,
		Text:     post.Text,
		ReadOnly: post.ReadOnly,
	}, nil
}
