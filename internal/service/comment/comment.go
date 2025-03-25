package comment

import (
	"context"
	"ozon-test-project/internal/handlers/model"
	"ozon-test-project/internal/pkg/storage"
)

type CommentService interface {
	GetCommentsByPostIDs(ctx context.Context, postIDs []int64, limit, offset int32) (map[int64][]*model.CommentForPagination, error)
	GetCommentsByParentID(ctx context.Context, parentID int64, limit, offset int32) ([]*model.CommentForPagination, error)

	GetComments(ctx context.Context, postID int64, parentID *int64, limit, offset int32) ([]*model.CommentForPagination, error)
	CreateComment(ctx context.Context, postID int64, text string, parentID *int64) (int64, error)
}

type commentService struct {
	storage storage.Storage
}

func NewCommentService(storage storage.Storage) CommentService {
	return &commentService{
		storage: storage,
	}
}
