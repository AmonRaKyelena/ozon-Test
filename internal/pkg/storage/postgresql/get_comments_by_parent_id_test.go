package postgresql

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"ozon-test-project/internal/pkg/storage/model"
)

func TestGetCommentsByParentID(t *testing.T) {
	tests := []struct {
		name             string
		parentID         int64
		limit, offset    int32
		setupMock        func(mock sqlmock.Sqlmock)
		expectedComments []*model.CommentForPagination
		expectError      bool
	}{
		{
			name:     "success",
			parentID: 10,
			limit:    2,
			offset:   0,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `(?s)^SELECT\s+.*FROM comments c\s+WHERE c\.parentIdcomment = \$1\s+ORDER BY c\.id\s+LIMIT \$2 OFFSET \$3`
				rows := sqlmock.NewRows([]string{"id", "parentIdcomment", "idPost", "text", "has_child"}).
					AddRow(1, int64(10), int64(100), "comment1", true).
					AddRow(2, int64(10), int64(100), "comment2", false)
				mock.ExpectQuery(query).
					WithArgs(int64(10), int32(2), int32(0)).
					WillReturnRows(rows)
			},
			expectedComments: []*model.CommentForPagination{
				{
					Comment: model.Comment{
						ID:       1,
						ParentID: func() *int64 { id := int64(10); return &id }(),
						PostID:   100,
						Text:     "comment1",
					},
					HasChild: true,
				},
				{
					Comment: model.Comment{
						ID:       2,
						ParentID: func() *int64 { id := int64(10); return &id }(),
						PostID:   100,
						Text:     "comment2",
					},
					HasChild: false,
				},
			},
			expectError: false,
		},
		{
			name:     "query error",
			parentID: 10,
			limit:    1,
			offset:   0,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `(?s)^SELECT\s+.*FROM comments c\s+WHERE c\.parentIdcomment = \$1\s+ORDER BY c\.id\s+LIMIT \$2 OFFSET \$3`
				mock.ExpectQuery(query).
					WithArgs(int64(10), int32(1), int32(0)).
					WillReturnError(errors.New("query error"))
			},
			expectedComments: nil,
			expectError:      true,
		},
		{
			name:     "scan error",
			parentID: 10,
			limit:    1,
			offset:   0,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `(?s)^SELECT\s+.*FROM comments c\s+WHERE c\.parentIdcomment = \$1\s+ORDER BY c\.id\s+LIMIT \$2 OFFSET \$3`
				rows := sqlmock.NewRows([]string{"id", "parentIdcomment", "idPost", "text", "has_child"}).
					AddRow("invalid", int64(10), int64(100), "comment", true)
				mock.ExpectQuery(query).
					WithArgs(int64(10), int32(1), int32(0)).
					WillReturnRows(rows)
			},
			expectedComments: nil,
			expectError:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := &postgresqlRepository{db: db}
			tc.setupMock(mock)

			comments, err := repo.GetCommentsByParentID(context.Background(), tc.parentID, tc.limit, tc.offset)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedComments, comments)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
