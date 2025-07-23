package main

import (
	"fmt"
	"log"
	"net/http"

	_ "auth-code-generator/docs" // Import generated docs by swag
	"auth-code-generator/handler"
	"auth-code-generator/internal/store"
	"auth-code-generator/service"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title 2FA Code Generator API
// @version 2.0
// @description This is a server for generating and verifying 2FA codes using a layered architecture.
// @host localhost:8080
// @BasePath /
func main() {
	codeStore, err := store.NewSqliteStore("./2fa.db")
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	codeService := service.NewCodeService(codeStore)
	codeHandler := handler.NewCodeHandler(codeService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /code", codeHandler.GenerateCode)
	mux.HandleFunc("POST /verify", codeHandler.VerifyCode)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)
	fmt.Printf("Swagger UI available at http://localhost:%s/swagger/index.html\n", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
