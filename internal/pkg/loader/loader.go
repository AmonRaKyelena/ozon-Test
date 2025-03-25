package loader

import (
	"context"
	"errors"
	"strconv"

	"ozon-test-project/internal/service/comment"

	"github.com/graph-gophers/dataloader"
)

type commentLoaderKeyCtx struct{}

type commentBatcher struct {
	commentService comment.CommentService
}

func InsertLoaderToContext(ctx context.Context, loader *dataloader.Loader) context.Context {
	ctx = context.WithValue(ctx, commentLoaderKeyCtx{}, loader)
	return ctx
}

func LoaderFromContext(ctx context.Context) (*dataloader.Loader, error) {
	loaderVal := ctx.Value(commentLoaderKeyCtx{})
	if loaderVal == nil {
		return nil, errors.New("no loader found in context")
	}
	loaderInstance, ok := loaderVal.(*dataloader.Loader)
	if !ok {
		return nil, errors.New("loader is not an *dataloader.Loader")
	}
	return loaderInstance, nil
}

func NewCommentLoader(commentService comment.CommentService) *dataloader.Loader {
	batcher := &commentBatcher{
		commentService: commentService,
	}
	return dataloader.NewBatchedLoader(batcher.loadBatch)
}

func FillPaginatioValue(ctx context.Context, limit, offset int32) context.Context {
	ctx = context.WithValue(ctx, limitKeyCtx{}, limit)
	ctx = context.WithValue(ctx, offsetKeyCtx{}, offset)

	return ctx
}

func (b *commentBatcher) loadBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := func(err error) []*dataloader.Result {
		results := make([]*dataloader.Result, len(keys))
		for i := range keys {
			results[i] = &dataloader.Result{Error: err}
		}
		return results
	}

	postIDs := make([]int64, len(keys))
	for i, key := range keys {
		id, err := strconv.ParseInt(key.String(), 10, 64)
		if err != nil {
			handleError(err)
		}
		postIDs[i] = id
	}

	limit, offset, err := extractPaginationValue(ctx)
	if err != nil {
		return handleError(err)
	}

	commentsMap, err := b.commentService.GetCommentsByPostIDs(ctx, postIDs, limit, offset)
	if err != nil {
		return handleError(err)
	}

	results := make([]*dataloader.Result, len(keys))
	for i, key := range keys {
		id, err := strconv.ParseInt(key.String(), 10, 64)
		if err != nil {
			handleError(err)
		}

		commentsForThisPost := commentsMap[id]
		results[i] = &dataloader.Result{Data: commentsForThisPost, Error: nil}
	}
	return results
}

type limitKeyCtx struct{}
type offsetKeyCtx struct{}

func extractPaginationValue(ctx context.Context) (int32, int32, error) {
	limitVal := ctx.Value(limitKeyCtx{})
	offsetVal := ctx.Value(offsetKeyCtx{})

	limit, ok := limitVal.(int32)
	if !ok {
		return 0, 0, errors.New("limit is not an int32")
	}

	offset, ok := offsetVal.(int32)
	if !ok {
		return 0, 0, errors.New("offset is not an int32")
	}

	return limit, offset, nil
}
