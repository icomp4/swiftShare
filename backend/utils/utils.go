package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"time"
)

type AuthCode struct {
	UserID string
	Code    string
	Expires time.Time
	Used bool
}

type UpdatePasswordStruct struct{
	NewInfo string
	NewPassword string
}
func SendConformationLink(userID string, userEmail string) (AuthCode, error) {
	from := "swiftshareauth@gmail.com"
	password := os.Getenv("EMAIL_PW")
	to := []string{
		userEmail,
	}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	randomBigInt, randErr := rand.Int(rand.Reader, big.NewInt(90000))
	if randErr != nil {
		return AuthCode{}, randErr
	}
	verificationCode := AuthCode{
		UserID: userID,
		Code: fmt.Sprint(randomBigInt.Int64() + 10000),
		Expires: time.Now().Add(time.Minute * 5),
		Used: false,
	}
	subject := "SwiftShare password change request."
	body := fmt.Sprintf("Please use the following code to confirm your password change: \n%s", verificationCode.Code)
	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)
	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		fmt.Println(err)
		return AuthCode{}, err
	}
	return verificationCode, nil
}
func VerifyEmailCode(userCode string, userID string, authCode *AuthCode) bool {
	if(fmt.Sprint(userID) != fmt.Sprint(authCode.UserID) || authCode.Used || authCode.Expires.Before(time.Now()) || userCode != authCode.Code){
		return false
	}
	return true
}