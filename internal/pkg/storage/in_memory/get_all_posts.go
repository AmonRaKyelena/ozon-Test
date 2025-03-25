package inmemory

import (
	"context"
	"ozon-test-project/internal/pkg/storage"
	"ozon-test-project/internal/pkg/storage/model"
)

func (r *inMemoryRepository) GetAllPosts(ctx context.Context, limit, offset int32) ([]model.Post, error) {
	if offset > int32(len(r.posts)) {
		return nil, storage.ErrOutOfRangePagination
	}

	postsWithComments := r.posts[offset:min(limit+offset, int32(len(r.posts)))]
	result := make([]model.Post, 0, len(postsWithComments))
	for _, postWithComment := range postsWithComments {
		result = append(result, model.Post{
			ID:       postWithComment.Post.ID,
			Title:    postWithComment.Post.Title,
			Text:     postWithComment.Post.Text,
			ReadOnly: postWithComment.Post.ReadOnly,
		})
	}

	return result, nil
}
