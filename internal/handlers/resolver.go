package handlers

import (
	"ozon-test-project/internal/service/comment"
	"ozon-test-project/internal/service/post"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	postService    post.PostService
	commentService comment.CommentService
}

func NewResolver(
	postService post.PostService,
	commentService comment.CommentService,
) *Resolver {
	return &Resolver{
		postService:    postService,
		commentService: commentService,
	}
}
