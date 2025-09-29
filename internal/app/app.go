package app

import (
	_ "commentsService/docs"
	"commentsService/internal/config"
	"commentsService/internal/handler"
	"commentsService/internal/repo/storage"
	"commentsService/internal/usecase/comment"
	"commentsService/pkg"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func dsn(cfg *config.Config) string {
	user := cfg.Postgres.POSTGRES_USER
	password := cfg.Postgres.POSTGRES_PASSWORD
	host := cfg.Postgres.POSTGRES_HOST
	port := cfg.Postgres.POSTGRES_PORT
	dbname := cfg.Postgres.POSTGRES_DB
	sslmode := cfg.Postgres.POSTGRES_SSLMODE

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)
}

func App(cfg *config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pkg.NewPostgres(dsn(cfg), ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repository := storage.NewRepository(pool)
	uc := comment.NewUseCase(repository)

	h := handler.NewCommentHandler(uc)

	r := gin.Default()
	r.LoadHTMLGlob("cmd/app/templates/*")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.PUT("/comments", h.ChangeComment)
	r.GET("/comments", h.GetComments)
	r.POST("/comments", h.SaveComment)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
