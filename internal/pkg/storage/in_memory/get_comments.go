package inmemory

import (
	"context"
	"fmt"
	"ozon-test-project/internal/pkg/storage"
	"ozon-test-project/internal/pkg/storage/model"

	"errors"
)

func (r *inMemoryRepository) GetCommentsByPostIDs(
	_ context.Context,
	postIDs []int64,
	limit, offset int32,
) (map[int64][]model.CommentForPagination, error) {
	result := map[int64][]model.CommentForPagination{}
	for _, postID := range postIDs {
		if int(postID) >= len(r.posts) {
			return nil, errors.Join(storage.ErrPostNotFound, fmt.Errorf("postID = %d", postID))
		}

		if offset != 0 && offset >= int32(len(r.posts[postID].Comments)) {
			return nil, storage.ErrOutOfRangePagination
		}

		for _, comment := range r.posts[postID].Comments[offset:min(offset+limit, int32(len(r.posts[postID].Comments)))] {
			result[postID] = append(result[postID], model.CommentForPagination{
				Comment:  *comment,
				HasChild: false,
			})

			if len(comment.Childs) > 0 {
				result[postID][len(result[postID])-1].HasChild = true
				result[postID][len(result[postID])-1].Comment.Childs = nil
			}
		}
	}
	return result, nil
}
