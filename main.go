package main

import (
	"html/template"
	"log"
	"net/http"
	"swiftShare/database"
	"swiftShare/handlers"
	"swiftShare/handlers/middleware"
	"swiftShare/handlers/validators"
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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if validate := validators.ExtractToken(r); validate == nil{
			http.Redirect(w,r, "/main", http.StatusSeeOther)
		}
		templates := template.Must(template.ParseFiles("static/login.html"))
		if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		templates := template.Must(template.ParseFiles("static/signup.html"))
		if err := templates.ExecuteTemplate(w, "signup.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		templates := template.Must(template.ParseFiles("static/main.html"))
		if err := templates.ExecuteTemplate(w, "main.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	signupHandler := middleware.Logger(http.HandlerFunc(handlers.SignUp))
	loginHandler := middleware.Logger(http.HandlerFunc(handlers.SignIn))
	logoutHandler := middleware.Logger(http.HandlerFunc(handlers.LogOut))
	deleteHandler := middleware.RequireAuth(middleware.Logger(http.HandlerFunc(handlers.DeleteAccount)))
	requestEmailHandler := middleware.RequireAuth(middleware.Logger(http.HandlerFunc(handlers.RequestEmail)))
	updatePassswordHandler := middleware.RequireAuth(middleware.Logger(http.HandlerFunc(handlers.UpatePassword)))

	mux.Handle("/api/v1/signup", signupHandler)
	mux.Handle("/api/v1/login", loginHandler)
	mux.Handle("/api/v1/logout", logoutHandler)
	mux.Handle("/api/v1/delete", deleteHandler)
	mux.Handle("/api/v1/request", requestEmailHandler)
	mux.Handle("/api/v1/password/update", updatePassswordHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
