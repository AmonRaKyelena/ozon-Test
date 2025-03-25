package comment

import (
	"context"
	"errors"
	"testing"

	"ozon-test-project/internal/handlers/model"
	modelDB "ozon-test-project/internal/pkg/storage/model"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGetComments(t *testing.T) {
	type args struct {
		ctx      context.Context
		postID   int64
		parentID *int64
		limit    int32
		offset   int32
	}

	parentID := int64(1)

	tests := []struct {
		name     string
		prepare  func(m mocksData, args args)
		args     args
		response []*model.CommentForPagination
		wantErr  bool
	}{
		{
			name: "success without parent",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetCommentsByPostIDsMock.
					Expect(args.ctx, []int64{args.postID}, args.limit, args.offset).
					Return(map[int64][]modelDB.CommentForPagination{
						args.postID: {
							{
								Comment: modelDB.Comment{
									ID:       1,
									PostID:   args.postID,
									ParentID: nil,
									Text:     "text1",
								},
								HasChild: false,
							},
							{
								Comment: modelDB.Comment{
									ID:       2,
									PostID:   args.postID,
									ParentID: nil,
									Text:     "text2",
								},
								HasChild: false,
							},
						},
					}, nil)
			},
			args: args{
				ctx:    context.Background(),
				postID: 1,
				limit:  10,
				offset: 0,
			},
			response: []*model.CommentForPagination{
				{
					ID:       1,
					PostID:   1,
					ParentID: nil,
					Text:     "text1",
					HasChild: false,
				},
				{
					ID:       2,
					PostID:   1,
					ParentID: nil,
					Text:     "text2",
					HasChild: false,
				},
			},
			wantErr: false,
		},
		{
			name: "success with parent",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetCommentsByParentIDMock.
					Expect(args.ctx, *args.parentID, args.limit, args.offset).
					Return([]*modelDB.CommentForPagination{
						{
							Comment: modelDB.Comment{
								ID:       1,
								PostID:   2,
								ParentID: args.parentID,
								Text:     "text1",
							},
							HasChild: true,
						},
					}, nil)
			},
			args: args{
				ctx:      context.Background(),
				postID:   1,
				parentID: &parentID,
				limit:    10,
				offset:   0,
			},
			response: []*model.CommentForPagination{
				{
					ID:       1,
					PostID:   2,
					ParentID: &parentID,
					Text:     "text1",
					HasChild: true,
				},
			},
			wantErr: false,
		},
		{
			name: "success without parent",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetCommentsByPostIDsMock.
					Expect(args.ctx, []int64{args.postID}, args.limit, args.offset).
					Return(nil, errors.New("some error"))
			},
			args: args{
				ctx:    context.Background(),
				postID: 1,
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

			postService := NewCommentService(m.storageMock)
			got, err := postService.GetComments(test.args.ctx, test.args.postID, test.args.parentID, test.args.limit, test.args.offset)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.response, got)
			}
		})
	}
}
