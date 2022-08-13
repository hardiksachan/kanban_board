package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"kanban_board/internal/shared/logging"
	"kanban_board/internal/users/core/ports"
	"kanban_board/internal/users/handlers"
	"kanban_board/internal/users/repository/native"
	"kanban_board/internal/users/repository/postgres"
	"kanban_board/internal/users/repository/postgres/user/dao"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := logging.NewDefaultLogger(log.New(os.Stdout, "api", log.LstdFlags))
	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error(fmt.Sprintf("Can't connect to database. %s", err.Error()))
		os.Exit(1)
	}

	usersHandler := handlers.NewUsersHandler(
		ports.NewAuthService(
			postgres.NewUserStore(dao.New(db)),
			native.NewSessionStore(),
		),
		logger,
		validator.New(),
	)

	router := mux.NewRouter()

	router.HandleFunc("/users/signup", usersHandler.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/users/login", usersHandler.LogIn).Methods(http.MethodPost)
	router.Handle("/users/logout", usersHandler.AuthMiddleware(http.HandlerFunc(usersHandler.LogOut))).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":9090", router))
}
