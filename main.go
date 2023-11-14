package main

import (
	"log"
	"net/http"
	"swiftShare/backend/database"
	"swiftShare/backend/handlers"
	"swiftShare/backend/handlers/middleware"
)

func main() {
	if err := database.Start(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
	if err := database.DB.Ping(); err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/signup", middleware.Logger(http.HandlerFunc(handlers.SignUp)))
	mux.Handle("/login", middleware.Logger(http.HandlerFunc(handlers.SignIn)))
	mux.Handle("/logout", middleware.Logger(http.HandlerFunc(handlers.LogOut)))
	mux.Handle("/delete", middleware.RequireAuth(http.HandlerFunc(handlers.DeleteAccount)))
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
