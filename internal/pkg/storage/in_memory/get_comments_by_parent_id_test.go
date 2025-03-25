package inmemory

import (
	"context"
	"ozon-test-project/internal/pkg/storage/model"
	"testing"

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
		args     args
		prepare  func(repo *inMemoryRepository, args args)
		response []*model.CommentForPagination
		wantErr  bool
	}{
		{
			name: "success pagination=1",
			prepare: func(repo *inMemoryRepository, args args) {
				postID, _ := repo.CreatePost(args.ctx, model.Post{})
				parentID, _ := repo.CreateComment(args.ctx, model.Comment{
					PostID: postID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &parentID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &parentID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &parentID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &parentID,
				})
				forChild, _ := repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &parentID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &forChild,
				})
			},
			response: []*model.CommentForPagination{
				{
					Comment: model.Comment{
						ID:       2,
						Text:     "",
						ParentID: &parentID,
					},
					HasChild: false,
				},
				{
					Comment: model.Comment{
						ID:       3,
						Text:     "",
						ParentID: &parentID,
					},
					HasChild: false,
				},
				{
					Comment: model.Comment{
						ID:       4,
						Text:     "",
						ParentID: &parentID,
					},
					HasChild: false,
				},
				{
					Comment: model.Comment{
						ID:       5,
						Text:     "",
						ParentID: &parentID,
					},
					HasChild: false,
				},
				{
					Comment: model.Comment{
						ID:       6,
						Text:     "",
						ParentID: &parentID,
					},
					HasChild: true,
				},
			},
			args: args{
				ctx:      context.Background(),
				parentID: 1,
				limit:    5,
				offset:   0,
			},
			wantErr: false,
		},
		{
			name: "success pagination=2",
			prepare: func(repo *inMemoryRepository, args args) {
				postID, _ := repo.CreatePost(args.ctx, model.Post{})
				parentID, _ := repo.CreateComment(args.ctx, model.Comment{
					PostID: postID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &parentID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &parentID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &parentID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &parentID,
				})
				forChild, _ := repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &parentID,
				})
				repo.CreateComment(args.ctx, model.Comment{
					PostID:   postID,
					ParentID: &forChild,
				})
			},
			response: []*model.CommentForPagination{
				{
					Comment: model.Comment{
						ID:       4,
						Text:     "",
						ParentID: &parentID,
					},
					HasChild: false,
				},
				{
					Comment: model.Comment{
						ID:       5,
						Text:     "",
						ParentID: &parentID,
					},
					HasChild: false,
				},
				{
					Comment: model.Comment{
						ID:       6,
						Text:     "",
						ParentID: &parentID,
					},
					HasChild: true,
				},
			},
			args: args{
				ctx:      context.Background(),
				parentID: 1,
				limit:    5,
				offset:   2,
			},
			wantErr: false,
		},
		{
			name: "not found parent",
			args: args{
				ctx:      context.Background(),
				parentID: 1,
				limit:    5,
				offset:   0,
			},
			wantErr: true,
		},
		{
			name: "not found parent",
			prepare: func(repo *inMemoryRepository, args args) {
				postID, _ := repo.CreatePost(args.ctx, model.Post{})
				repo.CreateComment(args.ctx, model.Comment{
					PostID: postID,
				})
			},
			args: args{
				ctx:      context.Background(),
				parentID: 1,
				limit:    5,
				offset:   1230,
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

			got, err := inMemoryRepository.GetCommentsByParentID(test.args.ctx, test.args.parentID, test.args.limit, test.args.offset)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.response, got)
			}
		})
	}
}
