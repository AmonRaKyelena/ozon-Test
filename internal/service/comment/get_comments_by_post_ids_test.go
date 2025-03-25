package comment

import (
	"context"
	"errors"
	"testing"

	"github.com/AmonRaKyelena/ozon-Test/internal/handlers/model"
	modelDB "github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGetCommentsByPostIDs(t *testing.T) {
	type args struct {
		ctx     context.Context
		postIDs []int64
		limit   int32
		offset  int32
	}

	parentID := int64(1)

	tests := []struct {
		name     string
		prepare  func(m mocksData, args args)
		args     args
		response map[int64][]*model.CommentForPagination
		wantErr  bool
	}{
		{
			name: "success",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetCommentsByPostIDsMock.
					Expect(args.ctx, args.postIDs, args.limit, args.offset).
					Return(map[int64][]modelDB.CommentForPagination{
						args.postIDs[0]: {
							{
								Comment: modelDB.Comment{
									ID:       1,
									PostID:   args.postIDs[0],
									ParentID: nil,
									Text:     "text1",
								},
								HasChild: false,
							},
							{
								Comment: modelDB.Comment{
									ID:       2,
									PostID:   args.postIDs[0],
									ParentID: nil,
									Text:     "text2",
								},
								HasChild: false,
							},
						},
						args.postIDs[1]: {
							{
								Comment: modelDB.Comment{
									ID:       3,
									PostID:   args.postIDs[1],
									ParentID: &parentID,
									Text:     "text3",
								},
								HasChild: false,
							},
							{
								Comment: modelDB.Comment{
									ID:       4,
									PostID:   args.postIDs[1],
									ParentID: &parentID,
									Text:     "text4",
								},
								HasChild: false,
							},
						},
					}, nil)
			},
			args: args{
				ctx:     context.Background(),
				postIDs: []int64{1, 2},
				limit:   10,
				offset:  0,
			},
			response: map[int64][]*model.CommentForPagination{
				1: {
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
				2: {
					{
						ID:       3,
						PostID:   2,
						ParentID: &parentID,
						Text:     "text3",
						HasChild: false,
					},
					{
						ID:       4,
						PostID:   2,
						ParentID: &parentID,
						Text:     "text4",
						HasChild: false,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "failed get posts by ids",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetCommentsByPostIDsMock.
					Expect(args.ctx, args.postIDs, args.limit, args.offset).
					Return(nil, errors.New("some error"))
			},
			args: args{
				ctx:     context.Background(),
				postIDs: []int64{1, 2},
				limit:   10,
				offset:  0,
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
			got, err := postService.GetCommentsByPostIDs(test.args.ctx, test.args.postIDs, test.args.limit, test.args.offset)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.response, got)
			}
		})
	}
}
