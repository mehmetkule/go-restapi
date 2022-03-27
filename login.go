package main


import (
	"encoding/json"
	"github.com/mehmetkule/go-restapi/internal/dto"
	"github.com/mehmetkule/go-restapi/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

)

// Login logs user to the server
func (app *App) Login(writer http.ResponseWriter, request *http.Request) {
	var req dto.User

	var err error
	if err = json.NewDecoder(request.Body).Decode(&req); err != nil {
		app.RenderErrorResponse(writer, http.StatusForbidden, err, "Failed to login")
		return
	}
	logger.Logger().Info("Login for",zap.String("email",req.Email))
	var user *dto.User
	var errResponse *dto.ErrorResponse
	if user, errResponse = app.userRepo.GetUser(req.Email); err != nil {
		app.RenderErrorResponse(writer, errResponse.Status, errResponse.Error, errResponse.Message)
		return
	}

	if valid := ComparePasswords(user.Password, []byte(req.Password)); !valid {
		app.RenderErrorResponse(writer, http.StatusBadRequest, err, "Invalid Username/Password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"isUsing2FA": false,
		"username":   user.Email,
		"exp":        time.Now().Add(time.Minute * 1000000).Unix(),
	})

	var tokenString string
	if tokenString, err = token.SignedString([]byte(app.conf.PasswordKey)); err != nil {
		app.RenderErrorResponse(writer, http.StatusForbidden, err, "Login Failed")
		return
	}
	app.RenderJSON(writer, http.StatusOK, dto.JWTToken{Token: tokenString})
}

// ComparePasswords compare given passwordws
func ComparePasswords(hashedPassword string, password []byte) bool {
	byteHash := []byte(hashedPassword)
	if err := bcrypt.CompareHashAndPassword(byteHash, password); err != nil {
		return false
	}
	return true
}
