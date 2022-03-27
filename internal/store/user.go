package store


import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/mehmetkule/go-restapi/internal/dto"
	"github.com/mehmetkule/go-restapi/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// UserRepo Struct
type UserRepo struct {
	DB *sql.DB
}

// CreateUser func
func (r *UserRepo) CreateUser(request dto.UserRequest) (*dto.UserResponse, *dto.ErrorResponse) {
	var lastInsertID uuid.UUID
	var err error
	if lastInsertID, err = r.insertCreateUser(request); err != nil {
		return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed to insert user"}
	}
	return &dto.UserResponse{
		ID:        lastInsertID.String(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
	}, nil
}

// insertCreateUser func
func (r *UserRepo) insertCreateUser(request dto.UserRequest) (uuid.UUID, error) {
	sql := "INSERT INTO users(first_name,last_name,email,password,is_2fa,created) VALUES($1,$2,$3,$4,$5,$6) returning id;"
	password := hashAndSalt([]byte(request.Password))
	row := r.DB.QueryRowContext(context.Background(), sql, request.FirstName, request.LastName, request.Email, password, request.IsUsing2FA, time.Now())
	var lastInsertID uuid.UUID
	return lastInsertID, row.Scan(&lastInsertID)
}

// FindUserByID fetches user by ID
func (r *UserRepo) FindUserByID(id uuid.UUID) (*dto.UserResponse, *dto.ErrorResponse) {
	logger.Logger().Debug("Finding user item",zap.String("email",id.String()))
	sqlQuery := "SELECT id,first_name,last_name,email FROM users WHERE id=$1"
	var rows *sql.Rows
	var err error
	if rows, err = r.DB.QueryContext(context.Background(), sqlQuery, id); err != nil {
		return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed find user"}
	}

	defer rows.Close()

	response := dto.UserResponse{}
	if rows.Next() {
		err = rows.Scan(&response.ID, &response.FirstName, &response.LastName, &response.Email)
		if err != nil {
			return nil, &dto.ErrorResponse{Status: http.StatusNotFound, Error: err, Message: fmt.Sprintf("User [%s] not found", id)}
		}
		return &response, nil
	}

	if err = rows.Err(); err != nil {
		return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed to fetch user"}
	}
	return nil, &dto.ErrorResponse{Status: http.StatusNotFound, Error: fmt.Errorf("not found"), Message: "user not found"}
}

// FindUserByEmail fetches user by email
func (r *UserRepo) FindUserByEmail(email string) (*dto.UserResponse, *dto.ErrorResponse) {
	logger.Logger().Debug("Finding user item",zap.String("email",email))
	sqlQuery := "SELECT id,first_name,last_name,email,is_2fa FROM users WHERE email=$1"
	var rows *sql.Rows
	var err error
	if rows, err = r.DB.QueryContext(context.Background(), sqlQuery, email); err != nil {
		return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed find user"}
	}

	defer rows.Close()

	response := dto.UserResponse{}
	if rows.Next() {
		err = rows.Scan(&response.ID, &response.FirstName, &response.LastName, &response.Email, &response.IsUsing2FA)
		if err != nil {
			return nil, &dto.ErrorResponse{Status: http.StatusNotFound, Error: err, Message: fmt.Sprintf("User [%s] not found", email)}
		}
		return &response, nil
	}
	if err = rows.Err(); err != nil {
		return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed to fetch user"}
	}
	return nil, &dto.ErrorResponse{Status: http.StatusNotFound, Error: fmt.Errorf("not found"), Message: "user not found"}
}

// GetUser gets user from DB for given email
func (r *UserRepo) GetUser(email string) (*dto.User, *dto.ErrorResponse) {
	logger.Logger().Debug("Finding user item",zap.String("email",email))
	sqlQuery := "SELECT email,password FROM users WHERE email=$1"
	var rows *sql.Rows
	var err error
	if rows, err = r.DB.Query(sqlQuery, email); err != nil {
		return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed find user"}
	}

	defer rows.Close()

	response := dto.User{}
	if rows.Next() {
		err = rows.Scan(&response.Email, &response.Password)
		if err != nil {
			return nil, &dto.ErrorResponse{Status: http.StatusNotFound, Error: err, Message: fmt.Sprintf("User [%s] not found", email)}
		}
		return &response, nil
	}
	if err = rows.Err(); err != nil {
		return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed to fetch user"}
	}
	return nil, &dto.ErrorResponse{Status: http.StatusNotFound, Error: fmt.Errorf("not found"), Message: "user not found"}
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

//DeleteUser delete user
func (r *UserRepo) DeleteUser(id string) *dto.ErrorResponse {
	_, err := r.DB.Exec("DELETE FROM users WHERE id=$1;", id)
	if err != nil {
		return &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed delete user"}
	}

	return nil
}
