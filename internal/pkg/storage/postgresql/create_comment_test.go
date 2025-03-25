package postgresql

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"ozon-test-project/internal/pkg/storage/model"
)

func TestCreateComment(t *testing.T) {
	tests := []struct {
		name        string
		comment     model.Comment
		prepare     func(mock sqlmock.Sqlmock)
		expectedID  int64
		expectError bool
	}{
		{
			name: "success",
			comment: model.Comment{
				PostID:   1,
				ParentID: func() *int64 { id := int64(10); return &id }(),
				Text:     "test comment",
			},
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO comments \(idPost, parentIdcomment, text\) VALUES \(\$1, \$2, \$3\) RETURNING id`
				rows := sqlmock.NewRows([]string{"id"}).AddRow(42)
				mock.ExpectQuery(query).
					WithArgs(1, func() *int64 { id := int64(10); return &id }(), "test comment").
					WillReturnRows(rows)
			},
			expectedID:  42,
			expectError: false,
		},
		{
			name: "error on insert",
			comment: model.Comment{
				PostID:   2,
				ParentID: nil,
				Text:     "error comment",
			},
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO comments \(idPost, parentIdcomment, text\) VALUES \(\$1, \$2, \$3\) RETURNING id`
				mock.ExpectQuery(query).
					WithArgs(2, nil, "error comment").
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
			tc.prepare(mock)

			id, err := repo.CreateComment(context.Background(), tc.comment)
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
