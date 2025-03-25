package post

import (
	"context"
	"ozon-test-project/internal/handlers/model"
	"ozon-test-project/internal/pkg/storage"
)

type PostService interface {
	CreatePost(ctx context.Context, newPost model.NewPost) (int64, error)
	GetAllPosts(ctx context.Context, limit, offset int32) ([]model.PostForPagination, error)
	GetPostByID(ctx context.Context, id int64) (*model.PostForPagination, error)
}

type postService struct {
	storage storage.Storage
}

func NewPostService(storage storage.Storage) PostService {
	return &postService{
		storage: storage,
	}
}
