package main

import (
	"context"
	"log"
	"net/http"

	"github.com/syarifid/bankx/internal/cfg"
	"github.com/syarifid/bankx/internal/handler"
	"github.com/syarifid/bankx/internal/repo"
	"github.com/syarifid/bankx/internal/service"
	"github.com/syarifid/bankx/pkg/env"
	"github.com/syarifid/bankx/pkg/postgre"
	"github.com/syarifid/bankx/pkg/router"
	"github.com/syarifid/bankx/pkg/validator"
)

func main() {
	env.LoadEnv()

	ctx := context.Background()
	router := router.NewRouter()
	conn := postgre.GetConn(ctx)
	defer conn.Close()
	validator := validator.New()

	cfg := cfg.Load()
	repo := repo.NewRepo(conn)
	service := service.NewService(repo, validator, cfg)
	handler.NewHandler(router, service, cfg)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln("fail start server:", err)
	}
}
