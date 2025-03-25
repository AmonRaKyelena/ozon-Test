package handlers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.68

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/handlers/model"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/logger"

	"go.uber.org/zap"
)

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input *model.NewComment) (int64, error) {
	log := logger.LoggerFromContext(ctx)

	commentID, err := r.commentService.CreateComment(ctx, input.PostID, input.Text, input.ParentID)
	if err != nil {
		log.Error("failed create comment", zap.Error(err))
		return 0, err
	}
	return commentID, nil
}

// CommentsOnPost is the resolver for the commentsOnPost field.
func (r *queryResolver) CommentsOnPost(ctx context.Context, postID int64, parentID *int64, limit *int32, offset *int32) ([]*model.CommentForPagination, error) {
	log := logger.LoggerFromContext(ctx)

	comments, err := r.commentService.GetComments(ctx, postID, parentID, *limit, *offset)
	if err != nil {
		log.Error("failed get post", zap.Int64("post_id", postID), zap.Error(err))
		return nil, err
	}
	return comments, nil
}
