package storage

import "errors"

var (
	ErrPostNotFound         = errors.New("post not found")
	ErrCommentNotFound      = errors.New("comment not found")
	ErrOutOfRangePagination = errors.New("out of range pagination")
)
