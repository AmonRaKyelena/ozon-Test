package model

type Comment struct {
	ID       int64
	ParentID *int64
	PostID   int64
	Text     string
	Childs   []*Comment
}

type CommentForPagination struct {
	Comment  Comment
	HasChild bool
}
