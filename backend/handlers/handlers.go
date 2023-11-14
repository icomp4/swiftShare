package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"swiftShare/backend/controllers"
	"swiftShare/backend/handlers/messages"
	"swiftShare/backend/handlers/validators"
	"swiftShare/backend/models"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode response body.", http.StatusMethodNotAllowed)
		return
	}
	if len(user.Username) < 3 {
		http.Error(w, "Username must be at least 3 characters long.", http.StatusMethodNotAllowed)
		return
	}
	if !validators.PasswordIsValid(user.Password) {
		http.Error(w, "Please choose a stronger password.", http.StatusMethodNotAllowed)
		return
	}
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Fields must not be blank.", http.StatusMethodNotAllowed)
		return
	}
	if err := controllers.SignUp(user); err != nil {
		http.Error(w, "Email or username already in use.", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	resp := messages.BasicSuccessMessage{
		Status:  "Success",
		Message: "User " + user.Username + " has succsessfully signed up.",
	}
	json, _ := json.Marshal(resp)
	w.Write(json)
}
func SignIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if r.Method != "POST" {
		WriteError(w, http.StatusMethodNotAllowed, "Method must be POST")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		WriteError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}
	if err := controllers.SignIn(user); err != nil {
		WriteError(w, http.StatusBadRequest, "Incorrect login information")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		WriteError(w, http.StatusBadGateway, "Failed to created token")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Path:     "/",
		HttpOnly: true,
	})
	resp := messages.BasicSuccessMessage{
		Status:  "Sucess",
		Message: "User " + user.Username + " has successfully logged in.",
	}
	json, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
func LogOut(w http.ResponseWriter, r *http.Request){
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	io.WriteString(w,"You have successfully logged out")
}
func Validate(w http.ResponseWriter, r *http.Request){
	io.WriteString(w, "You are authorized!")
}
func WriteError(w http.ResponseWriter, httpStatus int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	resp := messages.BasicFailMessage{
		Status:  "Failed",
		Message: message,
	}
	json, _ := json.Marshal(resp)
	w.Write(json)
}
