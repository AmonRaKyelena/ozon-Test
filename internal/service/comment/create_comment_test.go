package comment

import (
	"context"
	"errors"
	"testing"

	modelDB "ozon-test-project/internal/pkg/storage/model"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreateComment(t *testing.T) {
	type args struct {
		ctx      context.Context
		postID   int64
		text     string
		parentID *int64
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
				m.storageMock.GetPostByIDMock.Expect(args.ctx, args.postID).Return(&modelDB.Post{
					ID:       args.postID,
					Title:    "title1",
					Text:     "text1",
					ReadOnly: false,
				}, nil)
				m.storageMock.CreateCommentMock.Expect(args.ctx, modelDB.Comment{
					PostID:   args.postID,
					Text:     args.text,
					ParentID: args.parentID,
				}).Return(1, nil)
			},
			args: args{
				ctx:      context.Background(),
				postID:   1,
				text:     "some text",
				parentID: nil,
			},
			response: 1,
			wantErr:  false,
		},
		{
			name: "failed get post from storage",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetPostByIDMock.Expect(args.ctx, args.postID).Return(nil, errors.New("some error"))
			},
			args: args{
				ctx:      context.Background(),
				postID:   1,
				text:     "some text",
				parentID: nil,
			},
			wantErr: true,
		},
		{
			name: "can't comment read only post",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetPostByIDMock.Expect(args.ctx, args.postID).Return(&modelDB.Post{
					ID:       args.postID,
					Title:    "title1",
					Text:     "text1",
					ReadOnly: true,
				}, nil)
			},
			args: args{
				ctx:      context.Background(),
				postID:   1,
				text:     "some text",
				parentID: nil,
			},
			wantErr: true,
		},
		{
			name: "failed create comment",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetPostByIDMock.Expect(args.ctx, args.postID).Return(&modelDB.Post{
					ID:       args.postID,
					Title:    "title1",
					Text:     "text1",
					ReadOnly: false,
				}, nil)
				m.storageMock.CreateCommentMock.Expect(args.ctx, modelDB.Comment{
					PostID:   args.postID,
					Text:     args.text,
					ParentID: args.parentID,
				}).Return(0, errors.New("some error"))
			},
			args: args{
				ctx:      context.Background(),
				postID:   1,
				text:     "some text",
				parentID: nil,
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

			postService := NewCommentService(m.storageMock)
			got, err := postService.CreateComment(test.args.ctx, test.args.postID, test.args.text, test.args.parentID)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.response, got)
			}
		})
	}
}
