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

func TestGetCommentsByParentID(t *testing.T) {
	type args struct {
		ctx      context.Context
		parentID int64
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
			name: "success",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetCommentsByParentIDMock.
					Expect(args.ctx, args.parentID, args.limit, args.offset).
					Return([]*modelDB.CommentForPagination{
						{
							Comment: modelDB.Comment{
								ID:       1,
								PostID:   2,
								ParentID: &args.parentID,
								Text:     "text1",
							},
							HasChild: true,
						},
					}, nil)
			},
			args: args{
				ctx:      context.Background(),
				parentID: parentID,
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
			name: "failed get comments",
			prepare: func(m mocksData, args args) {
				m.storageMock.GetCommentsByParentIDMock.
					Expect(args.ctx, args.parentID, args.limit, args.offset).
					Return(nil, errors.New("some error"))
			},
			args: args{
				ctx:      context.Background(),
				parentID: parentID,
				limit:    10,
				offset:   0,
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
			got, err := postService.GetCommentsByParentID(test.args.ctx, test.args.parentID, test.args.limit, test.args.offset)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.response, got)
			}
		})
	}
}
