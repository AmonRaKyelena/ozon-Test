package storage

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"
)

type Storage interface {
	CreatePost(ctx context.Context, post model.Post) (int64, error)
	GetAllPosts(ctx context.Context, limit, offset int32) ([]model.Post, error)
	GetPostByID(ctx context.Context, id int64) (*model.Post, error)

	GetCommentsByPostIDs(ctx context.Context, postIDs []int64, limit, offset int32) (map[int64][]model.CommentForPagination, error)
	GetCommentsByParentID(ctx context.Context, parentID int64, limit, offset int32) ([]*model.CommentForPagination, error)
	CreateComment(ctx context.Context, comment model.Comment) (int64, error)
}
