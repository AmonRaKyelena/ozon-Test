package postgresql

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"
)

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name        string
		post        model.Post
		setupMock   func(mock sqlmock.Sqlmock)
		expectedID  int64
		expectError bool
	}{
		{
			name: "success",
			post: model.Post{
				Title:    "Test Title",
				Text:     "Test Text",
				ReadOnly: false,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO posts \(title, text, readOnly\) VALUES \(\$1, \$2, \$3\) RETURNING id`
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(query).
					WithArgs("Test Title", "Test Text", false).
					WillReturnRows(rows)
			},
			expectedID:  1,
			expectError: false,
		},
		{
			name: "error on insert",
			post: model.Post{
				Title:    "Error Title",
				Text:     "Error Text",
				ReadOnly: true,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO posts \(title, text, readOnly\) VALUES \(\$1, \$2, \$3\) RETURNING id`
				mock.ExpectQuery(query).
					WithArgs("Error Title", "Error Text", true).
					WillReturnError(errors.New("insert error"))
			},
			expectedID:  0,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := &postgresqlRepository{db: db}
			tc.setupMock(mock)

			id, err := repo.CreatePost(context.Background(), tc.post)
			if tc.expectError {
				require.Error(t, err)
				require.Equal(t, int64(0), id)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedID, id)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
