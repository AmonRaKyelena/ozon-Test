package post

import (
	"context"
	"ozon-test-project/internal/handlers/model"
	modelDB "ozon-test-project/internal/pkg/storage/model"
)

func (s *postService) CreatePost(ctx context.Context, newPost model.NewPost) (int64, error) {
	postID, err := s.storage.CreatePost(ctx, modelDB.Post{
		Title:    newPost.Title,
		Text:     newPost.Text,
		ReadOnly: newPost.ReadOnly,
	})
	if err != nil {
		return 0, err
	}
	return postID, nil
}
