package controllers

import (
	"log"
	"strings"
	"swiftShare/backend/database"
	"swiftShare/backend/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(user models.User) error {
	user.Username = strings.ToLower(user.Username)
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}
	user.Password = string(password)
	createUser := `INSERT INTO users (username, email, password, pfpurl) VALUES ($1, $2, $3, $4)`
	if _, err := database.DB.Exec(createUser, user.Username, user.Email, user.Password, user.PfpUrl); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
func SignIn(user models.User) (int, error) {
	query := "SELECT * FROM users WHERE username = $1"
	rows, err := database.DB.Query(query, user.Username)
	if err != nil {
		panic(err)
	}
	var id int
	var username string
	var email string
	var password string
	var pfpurl string
	var createdAt time.Time
	var deletedAt time.Time
	for rows.Next() {
		err := rows.Scan(&id, &username, &email, &password, &pfpurl, &createdAt, &deletedAt)
		if err != nil {
			panic(err)
		}
	}
	defer rows.Close()
	err2 := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err2 != nil {
		return -1, err2
	}
	return id, nil
}
func DeleteAccount(userID string) error {
	deleteAcc := `DELETE FROM users WHERE id = $1`
	if _, err := database.DB.Exec(deleteAcc, userID); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
func UpatePassword(userID string, newPassword string) error {
	password, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}
	deleteAcc := `UPDATE users SET password = $1 WHERE id = $2`
	if _, err := database.DB.Exec(deleteAcc, string(password), userID); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
func UpatePfp(userID string, newURL string) error {
	updatePfp := `UPDATE users SET pfpurl = $1 WHERE id = $2`
	if _, err := database.DB.Exec(updatePfp, newURL, userID); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
func FindByID(userID string) (models.User, error) {
	findUser := "SELECT * FROM users WHERE id = $1"
	rows, err := database.DB.Query(findUser, userID)
	if err != nil {
		return models.User{}, err
	}
	var id int
	var username string
	var email string
	var password string
	var pfpurl string
	var createdAt time.Time
	var updatedAt time.Time
	for rows.Next() {
		err := rows.Scan(&id, &username, &email, &password, &pfpurl, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return models.User{}, err
		}
	}
	defer rows.Close()
	user := models.User{
		ID:        uint(id),
		Username:  userID,
		Email:     email,
		Password:  password,
		PfpUrl:    pfpurl,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return user, nil
}
