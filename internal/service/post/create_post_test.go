package post

import (
	"context"
	"errors"
	"github.com/AmonRaKyelena/ozon-Test/internal/handlers/model"
	"testing"

	modelDB "github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	type args struct {
		ctx     context.Context
		newPost model.NewPost
	}

	tests := []struct {
		name     string
		prepare  func(m mocksData, args args)
		args     args
		response int64
		wantErr  bool
	}{
		{
			name: "success",
			prepare: func(m mocksData, args args) {
				m.storageMock.CreatePostMock.Expect(args.ctx, modelDB.Post{
					Text:     args.newPost.Text,
					Title:    args.newPost.Title,
					ReadOnly: args.newPost.ReadOnly,
				}).Return(1, nil)
			},
			args: args{
				ctx: context.Background(),
				newPost: model.NewPost{
					Title: "title1",
					Text:  "text1",
				},
			},
			response: 1,
			wantErr:  false,
		},
		{
			name: "failed create post",
			prepare: func(m mocksData, args args) {
				m.storageMock.CreatePostMock.Expect(args.ctx, modelDB.Post{
					Text:     args.newPost.Text,
					Title:    args.newPost.Title,
					ReadOnly: args.newPost.ReadOnly,
				}).Return(0, errors.New("some error"))
			},
			args: args{
				ctx: context.Background(),
				newPost: model.NewPost{
					Title: "title1",
					Text:  "text1",
				},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := minimock.NewController(t)
			m := newMock(ctrl)

			if test.prepare != nil {
				test.prepare(m, test.args)
			}

			postService := NewPostService(m.storageMock)
			got, err := postService.CreatePost(test.args.ctx, test.args.newPost)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.response, got)
			}
		})
	}
}
