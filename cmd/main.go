package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"github.com/AmonRaKyelena/ozon-Test/internal/graph"
	"github.com/AmonRaKyelena/ozon-Test/internal/handlers"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/config"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/loader"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/logger"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage"
	inmemory "github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/in_memory"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/postgresql"
	"github.com/AmonRaKyelena/ozon-Test/internal/service/comment"
	"github.com/AmonRaKyelena/ozon-Test/internal/service/post"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func middlewareHandler(
	gqlHandler http.Handler,
	commentService comment.CommentService,
	baseLogger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loaderInstance := loader.NewCommentLoader(commentService)

		ctx := context.Background()
		ctx = loader.InsertLoaderToContext(ctx, loaderInstance)
		ctx = logger.InsertLoggerToContext(ctx, baseLogger)

		baseLogger.Debug("Incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Any("headers", r.Header),
		)

		gqlHandler.ServeHTTP(w, r.WithContext(ctx))
	}
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	level := zapcore.InfoLevel
	if err := level.Set(cfg.LogLevel); err != nil {
		log.Printf("unknown log level %q, set to INFO by default", cfg.LogLevel)
		level = zapcore.InfoLevel
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	baseLogger, err := zapCfg.Build()
	if err != nil {
		log.Fatalf("failed to build logger: %v", err)
	}
	defer func() {
		if err := baseLogger.Sync(); err != nil {
			fmt.Printf("failed to sync logger: %v\n", err)
		}
	}()

	baseLogger.Info("Logger initialized",
		zap.String("log_level", cfg.LogLevel),
	)

	var store storage.Storage

	baseLogger.Info("start init repository", zap.String("mode", cfg.RepositoryMode))

	if cfg.RepositoryMode == string(config.PostgresqlMode) {
		db, err := sql.Open("postgres", cfg.PsqlInfo)
		if err != nil {
			log.Fatalf("failed to create connection db: %v", err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatalf("failed to ping db: %v", err)
		}
		store = postgresql.NewPostgresqlRepository(db)
		defer db.Close()
	} else {
		store = inmemory.NewInMemoryRepository()
	}

	postService := post.NewPostService(store)
	commentSerice := comment.NewCommentService(store)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: handlers.NewResolver(postService, commentSerice),
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	graphqlHandler := middlewareHandler(srv, commentSerice, baseLogger)
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", graphqlHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
