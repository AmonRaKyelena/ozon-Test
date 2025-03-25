package inmemory

import (
	"context"
	"ozon-test-project/internal/pkg/storage/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	type args struct {
		ctx  context.Context
		post model.Post
	}

	tests := []struct {
		name     string
		args     args
		response int64
	}{
		{
			name:     "success",
			response: 0,
			args: args{
				ctx: context.Background(),
				post: model.Post{
					Title: "title1",
					Text:  "text1",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewInMemoryRepository()
			got, err := repo.CreatePost(test.args.ctx, test.args.post)
			require.NoError(t, err)
			require.Equal(t, test.response, got)
		})
	}
}
