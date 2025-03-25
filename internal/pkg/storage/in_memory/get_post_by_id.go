package inmemory

import (
	"context"
	"errors"
	"fmt"
	"ozon-test-project/internal/pkg/storage"
	"ozon-test-project/internal/pkg/storage/model"
)

func (r *inMemoryRepository) GetPostByID(ctx context.Context, id int64) (*model.Post, error) {
	if int(id) >= len(r.posts) {
		return nil, errors.Join(storage.ErrPostNotFound, fmt.Errorf("postID = %d", id))
	}
	postWithComment := r.posts[id]
	return &model.Post{
		ID:       postWithComment.Post.ID,
		Title:    postWithComment.Post.Title,
		Text:     postWithComment.Post.Text,
		ReadOnly: postWithComment.Post.ReadOnly,
	}, nil
}
