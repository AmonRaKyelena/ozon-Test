package inmemory

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateComment(t *testing.T) {
	type args struct {
		ctx     context.Context
		comment model.Comment
	}

	parentID := int64(1)
	parentID2 := int64(122333)

	tests := []struct {
		name     string
		args     args
		prepare  func(repo *inMemoryRepository, args args)
		response int64
		wantErr  bool
	}{
		{
			name: "success without parent",
			prepare: func(repo *inMemoryRepository, args args) {
				repo.CreatePost(args.ctx, model.Post{})
			},
			response: 1,
			args: args{
				ctx: context.Background(),
				comment: model.Comment{
					PostID: 0,
					Text:   "text1",
				},
			},
			wantErr: false,
		},
		{
			name: "success with parent",
			prepare: func(repo *inMemoryRepository, args args) {
				postID, _ := repo.CreatePost(args.ctx, model.Post{})
				repo.CreateComment(args.ctx, model.Comment{
					PostID: postID,
				})
			},
			response: 2,
			args: args{
				ctx: context.Background(),
				comment: model.Comment{
					PostID:   0,
					Text:     "text1",
					ParentID: &parentID,
				},
			},
			wantErr: false,
		},
		{
			name: "failed with parent",
			prepare: func(repo *inMemoryRepository, args args) {
				repo.CreatePost(args.ctx, model.Post{})
			},
			args: args{
				ctx: context.Background(),
				comment: model.Comment{
					PostID:   0,
					Text:     "text1",
					ParentID: &parentID2,
				},
			},
			wantErr: true,
		},
		{
			name: "post not found",
			args: args{
				ctx: context.Background(),
				comment: model.Comment{
					PostID:   0,
					Text:     "text1",
					ParentID: &parentID2,
				},
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

			got, err := inMemoryRepository.CreateComment(test.args.ctx, test.args.comment)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.response, got)
			}
		})
	}
}
