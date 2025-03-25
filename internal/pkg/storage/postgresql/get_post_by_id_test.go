package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"ozon-test-project/internal/pkg/storage"
	"ozon-test-project/internal/pkg/storage/model"
)

func TestGetPostByID(t *testing.T) {
	tests := []struct {
		name         string
		id           int64
		setupMock    func(mock sqlmock.Sqlmock)
		expectedPost *model.Post
		expectError  bool
	}{
		{
			name: "success",
			id:   1,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, title, text, readOnly FROM posts WHERE id = \$1`
				rows := sqlmock.NewRows([]string{"id", "title", "text", "readOnly"}).
					AddRow(1, "Title 1", "Text 1", false)
				mock.ExpectQuery(query).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectedPost: &model.Post{
				ID:       1,
				Title:    "Title 1",
				Text:     "Text 1",
				ReadOnly: false,
			},
			expectError: false,
		},
		{
			name: "no rows",
			id:   2,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, title, text, readOnly FROM posts WHERE id = \$1`
				mock.ExpectQuery(query).
					WithArgs(int64(2)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedPost: nil,
			expectError:  true,
		},
		{
			name: "query error",
			id:   3,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, title, text, readOnly FROM posts WHERE id = \$1`
				mock.ExpectQuery(query).
					WithArgs(int64(3)).
					WillReturnError(errors.New("query error"))
			},
			expectedPost: nil,
			expectError:  true,
		},
		{
			name: "scan error",
			id:   4,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, title, text, readOnly FROM posts WHERE id = \$1`
				rows := sqlmock.NewRows([]string{"id", "title", "text", "readOnly"}).
					AddRow(4, "Title 4", "Text 4", "invalid_bool")
				mock.ExpectQuery(query).
					WithArgs(int64(4)).
					WillReturnRows(rows)
			},
			expectedPost: nil,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := &postgresqlRepository{db: db}
			tc.setupMock(mock)

			post, err := repo.GetPostByID(context.Background(), tc.id)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, post)
				if errors.Is(err, sql.ErrNoRows) {
					require.Contains(t, err.Error(), fmt.Sprintf("postID = %d", tc.id))
					require.ErrorIs(t, err, storage.ErrPostNotFound)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedPost, post)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
