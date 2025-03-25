package postgresql

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/model"
)

func TestGetCommentsByPostIDs(t *testing.T) {
	tests := []struct {
		name           string
		postIDs        []int64
		limit, offset  int32
		setupMock      func(mock sqlmock.Sqlmock)
		expectedResult map[int64][]model.CommentForPagination
		expectError    bool
	}{
		{
			name:    "success: comments for multiple posts",
			postIDs: []int64{1, 2},
			limit:   2,
			offset:  1,
			setupMock: func(mock sqlmock.Sqlmock) {
				queryRegex := `(?s)^.*FROM numbered\s+WHERE rn > \$2 AND rn <= \(\$2 \+ \$3\)\s+ORDER BY idPost, rn;`
				rows := sqlmock.NewRows([]string{"id", "parentIdcomment", "idPost", "text", "has_child"}).
					AddRow(10, int64(5), int64(1), "comment1", true).
					AddRow(11, nil, int64(2), "comment2", false)
				mock.ExpectQuery(queryRegex).
					WithArgs(pq.Array([]int64{1, 2}), int32(1), int32(2)).
					WillReturnRows(rows)
			},
			expectedResult: map[int64][]model.CommentForPagination{
				1: {
					{
						Comment: model.Comment{
							ID:       10,
							ParentID: func() *int64 { i := int64(5); return &i }(),
							PostID:   1,
							Text:     "comment1",
						},
						HasChild: true,
					},
				},
				2: {
					{
						Comment: model.Comment{
							ID:       11,
							ParentID: nil,
							PostID:   2,
							Text:     "comment2",
						},
						HasChild: false,
					},
				},
			},
			expectError: false,
		},
		{
			name:    "query error",
			postIDs: []int64{3},
			limit:   1,
			offset:  0,
			setupMock: func(mock sqlmock.Sqlmock) {
				queryRegex := `(?s)^.*FROM numbered\s+WHERE rn > \$2 AND rn <= \(\$2 \+ \$3\)\s+ORDER BY idPost, rn;`
				mock.ExpectQuery(queryRegex).
					WithArgs(pq.Array([]int64{3}), int32(0), int32(1)).
					WillReturnError(errors.New("query error"))
			},
			expectedResult: nil,
			expectError:    true,
		},
		{
			name:    "scan error",
			postIDs: []int64{4},
			limit:   1,
			offset:  0,
			setupMock: func(mock sqlmock.Sqlmock) {
				queryRegex := `(?s)^.*FROM numbered\s+WHERE rn > \$2 AND rn <= \(\$2 \+ \$3\)\s+ORDER BY idPost, rn;`
				rows := sqlmock.NewRows([]string{"id", "parentIdcomment", "idPost", "text", "has_child"}).
					AddRow("invalid", nil, int64(4), "comment", true)
				mock.ExpectQuery(queryRegex).
					WithArgs(pq.Array([]int64{4}), int32(0), int32(1)).
					WillReturnRows(rows)
			},
			expectedResult: nil,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := &postgresqlRepository{db: db}
			tc.setupMock(mock)

			result, err := repo.GetCommentsByPostIDs(context.Background(), tc.postIDs, tc.limit, tc.offset)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
