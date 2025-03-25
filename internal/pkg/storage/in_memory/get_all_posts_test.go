package inmemory

import (
	"context"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllPosts(t *testing.T) {
	type args struct {
		ctx    context.Context
		limit  int32
		offset int32
	}

	tests := []struct {
		name     string
		args     args
		prepare  func(repo *inMemoryRepository, args args)
		response []model.Post
		wantErr  bool
	}{
		{
			name: "success pagination=1",
			prepare: func(repo *inMemoryRepository, args args) {
				repo.CreatePost(args.ctx, model.Post{})
				repo.CreatePost(args.ctx, model.Post{})
				repo.CreatePost(args.ctx, model.Post{})
				repo.CreatePost(args.ctx, model.Post{})
				repo.CreatePost(args.ctx, model.Post{})
			},
			response: []model.Post{
				{
					ID:       0,
					Title:    "",
					Text:     "",
					ReadOnly: false,
				},
				{
					ID:       1,
					Title:    "",
					Text:     "",
					ReadOnly: false,
				},
				{
					ID:       2,
					Title:    "",
					Text:     "",
					ReadOnly: false,
				},
				{
					ID:       3,
					Title:    "",
					Text:     "",
					ReadOnly: false,
				},
				{
					ID:       4,
					Title:    "",
					Text:     "",
					ReadOnly: false,
				},
			},
			args: args{
				ctx:    context.Background(),
				limit:  5,
				offset: 0,
			},
			wantErr: false,
		},
		{
			name: "success pagination=2",
			prepare: func(repo *inMemoryRepository, args args) {
				repo.CreatePost(args.ctx, model.Post{})
				repo.CreatePost(args.ctx, model.Post{})
				repo.CreatePost(args.ctx, model.Post{})
				repo.CreatePost(args.ctx, model.Post{})
				repo.CreatePost(args.ctx, model.Post{})
			},
			response: []model.Post{
				{
					ID:       1,
					Title:    "",
					Text:     "",
					ReadOnly: false,
				},
			},
			args: args{
				ctx:    context.Background(),
				limit:  1,
				offset: 1,
			},
			wantErr: false,
		},
		{
			name: "out of range pagination",
			args: args{
				ctx:    context.Background(),
				limit:  1,
				offset: 123,
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

			got, err := inMemoryRepository.GetAllPosts(test.args.ctx, test.args.limit, test.args.offset)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.response, got)
			}
		})
	}
}
