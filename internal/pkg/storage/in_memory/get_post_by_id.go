package inmemory

import (
	"context"
	"errors"
	"fmt"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"
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
