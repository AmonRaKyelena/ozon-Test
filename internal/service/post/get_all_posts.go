package post

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/handlers/model"
)

func (s *postService) GetAllPosts(ctx context.Context, limit, offset int32) ([]model.PostForPagination, error) {
	posts, err := s.storage.GetAllPosts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]model.PostForPagination, 0, len(posts))
	for _, post := range posts {
		result = append(result, model.PostForPagination{
			ID:       post.ID,
			Title:    post.Title,
			Text:     post.Text,
			ReadOnly: post.ReadOnly,
		})
	}
	return result, nil
}
