package post

import (
	"context"
	"errors"
	"testing"

	"github.com/AmonRaKyelena/ozon-Test/internal/handlers/model"
	modelDB "github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGetPostByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}

	tests := []struct {
		name     string
		prepare  func(m mocksData, args args)
		args     args
		response *model.PostForPagination
		wantErr  bool
	}{
		{
			name: "success",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetPostByIDMock.Expect(args.ctx, args.id).Return(&modelDB.Post{
					ID:       args.id,
					Title:    "title1",
					Text:     "text1",
					ReadOnly: false,
				}, nil)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			response: &model.PostForPagination{
				ID:       1,
				Title:    "title1",
				Text:     "text1",
				ReadOnly: false,
			},
			wantErr: false,
		},
		{
			name: "failed get post by id",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetPostByIDMock.Expect(args.ctx, args.id).Return(nil, errors.New("some error"))
			},
			args: args{
				ctx: context.Background(),
				id:  1,
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
			got, err := postService.GetPostByID(test.args.ctx, test.args.id)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.response, got)
			}
		})
	}
}
