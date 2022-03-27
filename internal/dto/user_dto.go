package dto

import (
	"fmt"
	"net/http"
)

// UserRequest Struct
type UserRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsUsing2FA bool   `json:"is_2fa"`
}

// UserResponse Struct
type UserResponse struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	IsUsing2FA bool   `json:"is_2fa"`
}

// User db user struct
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWTToken jwt token
type JWTToken struct {
	Token string `json:"token"`
}

// ValidateUser validates request
func (request *UserRequest) ValidateUser() (int, error) {

	if request.FirstName == "" {
		return http.StatusBadRequest, fmt.Errorf("First Name is wrong")
	}
	if request.LastName == "" {
		return http.StatusBadRequest, fmt.Errorf("Last Name is wrong")
	}

	if request.Email == "" {
		return http.StatusBadRequest, fmt.Errorf("Email is wrong")
	}
	if request.Password == "" {
		return http.StatusBadRequest, fmt.Errorf("Password is wrong")
	}

	return http.StatusOK, nil
}

