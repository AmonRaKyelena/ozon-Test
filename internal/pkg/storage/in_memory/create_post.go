package inmemory

import (
	"context"
	"ozon-test-project/internal/pkg/storage/model"
)

func (r *inMemoryRepository) CreatePost(_ context.Context, post model.Post) (int64, error) {
	post.ID = int64(len(r.posts))
	r.posts = append(r.posts, postWithComments{
		Post: post,
	})
	return int64(len(r.posts) - 1), nil
}
