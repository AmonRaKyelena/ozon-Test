package inmemory

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCommentsByPostIDs(t *testing.T) {
	type args struct {
		ctx     context.Context
		postIDs []int64
		limit   int32
		offset  int32
	}

	tests := []struct {
		name     string
		args     args
		prepare  func(repo *inMemoryRepository, args args)
		response map[int64][]model.CommentForPagination
		wantErr  bool
	}{
		{
			name: "success",
			prepare: func(repo *inMemoryRepository, args args) {
				postID, _ := repo.CreatePost(args.ctx, model.Post{})
				parentID, _ := repo.CreateComment(args.ctx, model.Comment{
					PostID: postID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					ParentID: &parentID,
				})
			},
			response: map[int64][]model.CommentForPagination{
				0: {
					{
						Comment: model.Comment{
							ID:   1,
							Text: "",
						},
						HasChild: true,
					},
				},
			},
			args: args{
				ctx:     context.Background(),
				postIDs: []int64{0},
				limit:   5,
				offset:  0,
			},
			wantErr: false,
		},
		{
			name: "not found post",
			args: args{
				ctx:     context.Background(),
				postIDs: []int64{0},
				limit:   5,
				offset:  0,
			},
			wantErr: true,
		},
		{
			name: "out of range pagination",
			prepare: func(repo *inMemoryRepository, args args) {
				postID, _ := repo.CreatePost(args.ctx, model.Post{})
				parentID, _ := repo.CreateComment(args.ctx, model.Comment{
					PostID: postID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					ParentID: &parentID,
				})
			},
			args: args{
				ctx:     context.Background(),
				postIDs: []int64{0},
				limit:   5,
				offset:  1230,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inMemoryRepository := inMemoryRepository{
				posts: make([]postWithComments, 0),
			}

			if test.prepare != nil {
				test.prepare(&inMemoryRepository, test.args)
			}

			got, err := inMemoryRepository.GetCommentsByPostIDs(test.args.ctx, test.args.postIDs, test.args.limit, test.args.offset)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.response, got)
			}
		})
	}
}
