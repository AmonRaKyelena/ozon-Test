package inmemory

import (
	"context"
	"errors"
	"fmt"
	"ozon-test-project/internal/pkg/storage"
	"ozon-test-project/internal/pkg/storage/model"
)

func (r *inMemoryRepository) CreateComment(ctx context.Context, comment model.Comment) (int64, error) {
	if int(comment.PostID) >= len(r.posts) {
		return 0, errors.Join(storage.ErrPostNotFound, fmt.Errorf("postID = %d", comment.PostID))
	}

	r.commentCounter.Add(1)
	comment.ID = r.commentCounter.Load()

	r.commentMap.Store(comment.ID, &comment)
	if comment.ParentID == nil {
		r.posts[comment.PostID].Comments = append(r.posts[comment.PostID].Comments, &comment)
	} else if parentComment, ok := r.commentMap.Load(*comment.ParentID); ok {
		parentComment.(*model.Comment).Childs = append(parentComment.(*model.Comment).Childs, &comment)
	} else {
		return 0, errors.Join(storage.ErrCommentNotFound, fmt.Errorf("commentID = %d", comment.ParentID))
	}

	return comment.ID, nil
}
