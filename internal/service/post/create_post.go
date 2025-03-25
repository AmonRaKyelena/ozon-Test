package post

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/handlers/model"
	modelDB "github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"
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
