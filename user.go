package main

import (
	"encoding/json"
	"github.com/mehmetkule/go-restapi/internal/dto"
	"github.com/mehmetkule/go-restapi/logger"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

// Register defines api for registering users
func (app *App) Register(writer http.ResponseWriter, req *http.Request) {
	logger.Logger().Debug("Registering user")
	var reqBody []byte
	var err error
	// read the request body
	if reqBody, err = io.ReadAll(req.Body); err != nil {
		app.RenderErrorResponse(writer, http.StatusBadRequest, err, "Failed to read request body")
		return
	}

	// serialize to json
	var request dto.UserRequest
	if err = json.Unmarshal(reqBody, &request); err != nil {
		app.RenderErrorResponse(writer, http.StatusInternalServerError, err, "Failed to convert json code")
		return
	}
	// validate fields
	if status, errValidate := request.ValidateUser(); errValidate != nil {
		app.RenderErrorResponse(writer, status, errValidate, "Validation error")
		return
	}
	response, _ := app.userRepo.FindUserByEmail(request.Email)
	if response != nil {
		app.RenderErrorResponse(writer, http.StatusConflict, err, "Conflict Email")
		return
	}

	// insert message
	response, errResponse := app.userRepo.CreateUser(request)
	if errResponse != nil {
		app.RenderErrorResponse(writer, errResponse.Status, errResponse.Error, errResponse.Message)
		return
	}

	// render output
	app.RenderJSON(writer, http.StatusOK, response)

}

func (app *App) FindUserByID(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	uuid, err := uuid.FromString(id)
	if err != nil {
		app.RenderErrorResponse(writer, http.StatusBadRequest, err, "Invalid id")

	}
	response, errResponse := app.userRepo.FindUserByID(uuid)
	if errResponse != nil {
		app.RenderErrorResponse(writer, errResponse.Status, errResponse.Error, errResponse.Message)
		return
	}

	// render output
	app.RenderJSON(writer, http.StatusOK, response)
}

func (app *App) FindUserByEmail(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	email, ok := params["email"]
	if !ok {
		app.RenderErrorResponse(writer, http.StatusBadRequest, nil, "Invalid email")
		return
	}
	response, errResponse := app.userRepo.FindUserByEmail(email)
	if errResponse != nil {
		app.RenderErrorResponse(writer, errResponse.Status, errResponse.Error, errResponse.Message)
		return
	}
	// render output
	app.RenderJSON(writer, http.StatusOK, response)
}

func (app *App) DeleteUser(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	errResponse := app.userRepo.DeleteUser(id)
	if errResponse != nil {
		app.RenderErrorResponse(writer, errResponse.Status, errResponse.Error, errResponse.Message)
		return
	}
	app.RenderJSON(writer, http.StatusOK, "")
}

