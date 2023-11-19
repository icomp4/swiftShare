package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"swiftShare/controllers"
	"swiftShare/handlers/messages"
	"swiftShare/handlers/middleware"
	"swiftShare/handlers/validators"
	"swiftShare/models"
	"swiftShare/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)
var AuthCodeStore = make(map[string]utils.AuthCode)

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
	var userID int
	var err error
	if r.Method != "POST" {
		http.Error(w, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode request body",http.StatusBadRequest)
		return
	}
	if userID, err = controllers.SignIn(user); err != nil {
		http.Error(w, "Incorrect login information", http.StatusBadRequest)
		return
	}
	user.ID = uint(userID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprint(userID),
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		http.Error(w, "Failed to created token", http.StatusBadGateway)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Path:     "/",
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "user " + user.Username + " has successfully logged in.")
}
func LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().AddDate(0, 0, -1),
		Path:     "/",
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "You have successfully logged out")
}
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	user, ok := r.Context().Value(middleware.UserKey).(models.User)
	if !ok {
		http.Error(w, "Context does not include user information.", http.StatusInternalServerError)
		return
	}
	userid := fmt.Sprint(user.ID)
	if err := controllers.DeleteAccount(userid); err != nil {
		http.Error(w, "Could not delete account", http.StatusInternalServerError)
		return
	}
	fmt.Println("UserID: ",userid)
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().AddDate(0, 0, -1),
		Path:     "/",
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "You have successfully deleted your account and will now be logged out.")
}
func RequestEmail(w http.ResponseWriter, r *http.Request){
	var AuthCode utils.AuthCode
	var err error
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	user, ok := r.Context().Value(middleware.UserKey).(models.User)
	if !ok {
		http.Error(w, "Context does not include user information.", http.StatusInternalServerError)
		return
	}
	userid := fmt.Sprint(user.ID)
	if AuthCode, err = utils.SendConformationLink(userid, user.Email); err != nil{
		http.Error(w, "Could not send conformation email.", http.StatusInternalServerError)
		return
	}
	AuthCodeStore[userid] = AuthCode
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Please check email for conformation code.")
}
func UpatePassword(w http.ResponseWriter, r *http.Request){
	var userInfo utils.UpdatePasswordStruct
	if r.Method != "PATCH" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode request body",http.StatusBadRequest)
		return
	}
	user, ok := r.Context().Value(middleware.UserKey).(models.User)
	if !ok {
		http.Error(w, "Context does not include user information.", http.StatusInternalServerError)
		return
	}
	userid := fmt.Sprint(user.ID)
	authCode, ok := AuthCodeStore[userid]
	if !ok{
		http.Error(w, "Could not get authentication code, try generating a new one.", http.StatusInternalServerError)
		return
	}
	if err := utils.VerifyEmailCode(userInfo.Code, userid, &authCode); !err{
		http.Error(w, "Error verifying auth code.", http.StatusBadRequest)
		return
	}
	if err := controllers.UpatePassword(userid, userInfo.NewPassword); err != nil{
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "You have successfully updated your password")
}
func UpatePfp(w http.ResponseWriter, r *http.Request){
	var url map[string]string
	if r.Method != "PATCH" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "Failed to decode request body",http.StatusBadRequest)
		return
	}
	user, ok := r.Context().Value(middleware.UserKey).(models.User)
	if !ok {
		http.Error(w, "Context does not include user information.", http.StatusInternalServerError)
		return
	}
	userid := fmt.Sprint(user.ID)
	if err := controllers.UpatePfp(userid, url["newUrl"]); err != nil{
		http.Error(w, "Error updating profile picture", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "You have successfully updated your password")
}