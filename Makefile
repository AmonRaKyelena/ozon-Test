DB_CONN = "user=alexei dbname=ozon-test sslmode=disable password=1234"
MIGRATIONS_DIR = migrations

.PHONY: migrations-up migrations-down

migrations-up:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_CONN) up

migrations-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_CONN) down

mock:
	minimock -i ozon-test-project/internal/pkg/storage.Storage -o internal/pkg/storage/mocks/storage_mock.go
	minimock -i ozon-test-project/internal/service/comment.CommentService -o internal/service/comment/mocks/comment_service_mock.go
	minimock -i ozon-test-project/internal/service/post.PostService -o internal/service/post/mocks/post_service_mock.go
