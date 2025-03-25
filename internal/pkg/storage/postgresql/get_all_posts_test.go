package postgresql

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"ozon-test-project/internal/pkg/storage/model"
)

func TestGetAllPosts(t *testing.T) {
	tests := []struct {
		name          string
		limit, offset int32
		setupMock     func(mock sqlmock.Sqlmock)
		expectedPosts []model.Post
		expectError   bool
	}{
		{
			name:   "success",
			limit:  2,
			offset: 0,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, title, text, readOnly FROM posts ORDER BY id LIMIT \$1 OFFSET \$2`
				rows := sqlmock.NewRows([]string{"id", "title", "text", "readOnly"}).
					AddRow(1, "Title1", "Text1", false).
					AddRow(2, "Title2", "Text2", true)
				mock.ExpectQuery(query).
					WithArgs(int32(2), int32(0)).
					WillReturnRows(rows)
			},
			expectedPosts: []model.Post{
				{ID: 1, Title: "Title1", Text: "Text1", ReadOnly: false},
				{ID: 2, Title: "Title2", Text: "Text2", ReadOnly: true},
			},
			expectError: false,
		},
		{
			name:   "query error",
			limit:  1,
			offset: 0,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, title, text, readOnly FROM posts ORDER BY id LIMIT \$1 OFFSET \$2`
				mock.ExpectQuery(query).
					WithArgs(int32(1), int32(0)).
					WillReturnError(errors.New("query error"))
			},
			expectedPosts: nil,
			expectError:   true,
		},
		{
			name:   "scan error",
			limit:  1,
			offset: 0,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, title, text, readOnly FROM posts ORDER BY id LIMIT \$1 OFFSET \$2`
				rows := sqlmock.NewRows([]string{"id", "title", "text", "readOnly"}).
					AddRow(1, "Title1", "Text1", "invalid_bool")
				mock.ExpectQuery(query).
					WithArgs(int32(1), int32(0)).
					WillReturnRows(rows)
			},
			expectedPosts: nil,
			expectError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := &postgresqlRepository{db: db}

			tc.setupMock(mock)

			posts, err := repo.GetAllPosts(context.Background(), tc.limit, tc.offset)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedPosts, posts)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
