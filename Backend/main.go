package main

import (
	"Backend/db"
	"Backend/handler"
	"Backend/middleware"
	"Backend/repo"
	"Backend/usecase"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Setup database
	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("error connect DB: %s", err)
	}
	defer db.Close()

	// Setup server and middleware
	r := gin.Default()
	middleware := middleware.NewMiddleware()
	r.Use(middleware.Error())

	// Setup app (in layers)
	rp := repo.NewRepo(db)
	uc := usecase.NewUsecase(rp)
	hd := handler.NewHandler(uc)

	// Get synbols
	r.GET("/symbols", hd.GetSymbols)

	// Run server
	srv := &http.Server{
		Addr:    os.Getenv("SERVER_PORT"),
		Handler: r.Handler(),
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
