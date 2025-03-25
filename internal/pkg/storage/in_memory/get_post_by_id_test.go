package inmemory

import (
	"context"
	"testing"

	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"

	"github.com/stretchr/testify/require"
)

func TestGetPostByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		prepare func(repo *inMemoryRepository)
		args    args
		want    *model.Post
		wantErr bool
	}{
		{
			name: "post exists",
			prepare: func(repo *inMemoryRepository) {
				// Добавляем один пост
				repo.posts = append(repo.posts, postWithComments{
					Post: model.Post{
						ID:       0,
						Title:    "Title1",
						Text:     "Text1",
						ReadOnly: false,
					},
				})
			},
			args: args{
				ctx: context.Background(),
				id:  0,
			},
			want: &model.Post{
				ID:       0,
				Title:    "Title1",
				Text:     "Text1",
				ReadOnly: false,
			},
			wantErr: false,
		},
		{
			name: "post not found (empty repository)",
			args: args{
				ctx: context.Background(),
				id:  0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "id out of range",
			prepare: func(repo *inMemoryRepository) {
				repo.posts = []postWithComments{
					{
						Post: model.Post{
							ID:       0,
							Title:    "Title1",
							Text:     "Text1",
							ReadOnly: false,
						},
					},
				}
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &inMemoryRepository{
				posts: []postWithComments{},
			}
			if tt.prepare != nil {
				tt.prepare(repo)
			}
			got, err := repo.GetPostByID(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
