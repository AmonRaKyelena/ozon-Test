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

func TestGetAllPosts(t *testing.T) {
	type args struct {
		ctx    context.Context
		limit  int32
		offset int32
	}

	tests := []struct {
		name     string
		prepare  func(m mocksData, args args)
		args     args
		response []model.PostForPagination
		wantErr  bool
	}{
		{
			name: "success",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetAllPostsMock.Expect(args.ctx, args.limit, args.offset).Return([]modelDB.Post{
					{
						ID:       1,
						Title:    "title1",
						Text:     "text1",
						ReadOnly: false,
					},
					{
						ID:       2,
						Title:    "title2",
						Text:     "text2",
						ReadOnly: true,
					},
				}, nil)
			},
			args: args{
				ctx:    context.Background(),
				limit:  10,
				offset: 0,
			},
			response: []model.PostForPagination{
				{
					ID:       1,
					Title:    "title1",
					Text:     "text1",
					ReadOnly: false,
				},
				{
					ID:       2,
					Title:    "title2",
					Text:     "text2",
					ReadOnly: true,
				},
			},
			wantErr: false,
		},
		{
			name: "failed get all post",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetAllPostsMock.Expect(args.ctx, args.limit, args.offset).Return(nil, errors.New("some error"))
			},
			args: args{
				ctx:    context.Background(),
				limit:  10,
				offset: 0,
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
			got, err := postService.GetAllPosts(test.args.ctx, test.args.limit, test.args.offset)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.response, got)
			}
		})
	}
}
