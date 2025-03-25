package inmemory

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"
)

func (r *inMemoryRepository) GetCommentsByParentID(ctx context.Context, parentID int64, limit, offset int32) ([]*model.CommentForPagination, error) {
	val, ok := r.commentMap.Load(parentID)
	if !ok {
		return nil, storage.ErrCommentNotFound
	}
	parent := val.(*model.Comment)

	if offset != 0 && offset >= int32(len(parent.Childs)) {
		return nil, storage.ErrOutOfRangePagination
	}

	childs := parent.Childs[offset:min(offset+limit, int32(len(parent.Childs)))]
	result := make([]*model.CommentForPagination, 0, len(childs))
	for _, child := range childs {
		result = append(result, &model.CommentForPagination{
			Comment:  *child,
			HasChild: false,
		})

		if len(child.Childs) > 0 {
			result[len(result)-1].HasChild = true
			result[len(result)-1].Comment.Childs = nil
		}
	}

	return result, nil
}
