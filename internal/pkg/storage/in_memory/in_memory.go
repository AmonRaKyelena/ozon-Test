package inmemory

import (
	"ozon-test-project/internal/pkg/storage"
	"ozon-test-project/internal/pkg/storage/model"
	"sync"
	"sync/atomic"
)

type postWithComments struct {
	Post     model.Post
	Comments []*model.Comment
}

type inMemoryRepository struct {
	posts []postWithComments

	commentMap     sync.Map
	commentCounter atomic.Int64
}

func NewInMemoryRepository() storage.Storage {
	return &inMemoryRepository{
		posts:          make([]postWithComments, 0),
		commentCounter: atomic.Int64{},
	}
}
