package main

import (
	"html/template"
	"log"
	"net/http"
	"swiftShare/database"
	"swiftShare/handlers"
	"swiftShare/handlers/middleware"
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

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	templates := template.Must(template.ParseFiles("static/main.html"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Render the HTML template
		if err := templates.ExecuteTemplate(w, "main.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	signupHandler := middleware.RequireAuth(middleware.Logger(http.HandlerFunc(handlers.SignUp)))
	loginHandler := middleware.RequireAuth(middleware.Logger(http.HandlerFunc(handlers.SignIn)))
	logoutHandler := middleware.Logger(http.HandlerFunc(handlers.LogOut))
	deleteHandler := middleware.RequireAuth(middleware.Logger(http.HandlerFunc(handlers.DeleteAccount)))
	requestEmailHandler := middleware.RequireAuth(middleware.Logger(http.HandlerFunc(handlers.RequestEmail)))
	updatePassswordHandler := middleware.RequireAuth(middleware.Logger(http.HandlerFunc(handlers.UpatePassword)))

	mux.Handle("/signup", signupHandler)
	mux.Handle("/login", loginHandler)
	mux.Handle("/logout", logoutHandler)
	mux.Handle("/delete", deleteHandler)
	mux.Handle("/request", requestEmailHandler)
	mux.Handle("/password/update", updatePassswordHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
