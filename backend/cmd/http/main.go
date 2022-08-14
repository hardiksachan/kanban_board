package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	goredis "github.com/go-redis/redis/v9"
	"github.com/gorilla/mux"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/ports"
	"github.com/hardiksachan/kanban_board/backend/internal/users/handlers/auth"
	"github.com/hardiksachan/kanban_board/backend/internal/users/handlers/user"
	"github.com/hardiksachan/kanban_board/backend/internal/users/repository/jwt"
	"github.com/hardiksachan/kanban_board/backend/internal/users/repository/postgres"
	"github.com/hardiksachan/kanban_board/backend/internal/users/repository/postgres/user/dao"
	"github.com/hardiksachan/kanban_board/backend/internal/users/repository/redis"
	"github.com/hardiksachan/kanban_board/backend/shared/logging"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := logging.NewDefaultLogger(log.New(os.Stdout, "backend", log.LstdFlags))

	port := os.Getenv("PORT")
	if port == "" {
		logger.Error("environment variable PORT not set")
		os.Exit(1)
	}
	port = ":" + port

	pgUrl := os.Getenv("PG_URL")
	logger.Debug(fmt.Sprintf("Postgres URL in env: %s", pgUrl))

	rAddr := os.Getenv("REDIS_ADDR")
	logger.Debug(fmt.Sprintf("Redis addr in env: %s", rAddr))

	rPass := os.Getenv("REDIS_PASS")
	logger.Debug(fmt.Sprintf("Redis password in env: %s", rPass))

	jwtKey := os.Getenv("JWT_SIGNING_KEY")
	logger.Debug(fmt.Sprintf("jwt key in env: %s", jwtKey))

	pg, err := pgxpool.Connect(context.Background(), pgUrl)
	if err != nil {
		logger.Error(fmt.Sprintf("Can't connect to database. %s", err.Error()))
		os.Exit(1)
	}

	pgq := dao.New(pg)

	authHandler := auth.NewAuthHandler(
		ports.NewAuthService(
			postgres.NewUserStore(pgq),
			jwt.NewAccessTokenStore(jwtKey, time.Minute*10),
			redis.NewRefreshTokenStore(goredis.NewClient(&goredis.Options{
				Addr:     rAddr,
				Password: rPass,
			})),
		),
		logger,
		validator.New(),
	)

	userHandler := user.NewUserHandler(
		ports.NewUserService(postgres.NewUserMetadataStore(pgq)),
		logger,
		validator.New(),
	)

	router := mux.NewRouter()

	router.HandleFunc("/users/signup", authHandler.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/users/login", authHandler.LogIn).Methods(http.MethodPost)
	router.HandleFunc("/users/refresh", authHandler.RefreshAccessToken).Methods(http.MethodPost)
	router.Handle("/users/logout", authHandler.AuthMiddleware(http.HandlerFunc(authHandler.LogOut))).Methods(http.MethodPost)

	router.HandleFunc("/users/{user_id}", userHandler.Get).Methods(http.MethodGet)
	router.Handle("/users/{user_id}", authHandler.AuthMiddleware(http.HandlerFunc(userHandler.Update))).Methods(http.MethodPut)

	logger.Debug(fmt.Sprintf("Starting server on port: %s", port))

	// todo: graceful shutdown
	log.Fatal(http.ListenAndServe(port, router))
}
