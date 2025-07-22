package main

import (
	"fmt"
	"log"
	"net/http"

	_ "main/docs" // Import generated docs by swag
	"main/handler"
	"main/internal/store"
	"main/service"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title 2FA Code Generator API
// @version 2.0
// @description This is a server for generating and verifying 2FA codes using a layered architecture.
// @host localhost:8080
// @BasePath /
func main() {
	codeStore := store.NewInMemoryStore()
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
