package main

import (
	"github.com/gorilla/mux"
	"github.com/mehmetkule/go-restapi/internal/dto"
	"github.com/mehmetkule/go-restapi/internal/store"
	"github.com/mehmetkule/go-restapi/logger"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const maxuploadsize = 32 << 20

// AddFile adds file to database and creates a json response of the data
func (app *App) AddFile(writer http.ResponseWriter, req *http.Request) {
	// read the request body and request param
	logger.Logger().Debug("Adding File application")
	params := mux.Vars(req)
	parentID := params["parent_id"]
	if parentID == "" {
		app.RenderErrorResponse(writer, http.StatusBadRequest, nil, "Parent id is not defined")
		return
	}

	//req.Body = http.MaxBytesReader(writer, req.Body, maxuploadsize)
	if err := req.ParseMultipartForm(maxuploadsize); err != nil {
		app.RenderErrorResponse(writer, http.StatusBadRequest, err, "The uploaded file is too big. Please choose an file that's less than 1MB in size")
		return
	}
	data := req.MultipartForm.File["file"]

	var request dto.FileResponse

	_, err := request.ValidateFile(req)
	if err != nil {
		app.RenderErrorResponse(writer, http.StatusBadGateway, err, err.Error())
		return
	}

	// read data from multipartform to byte array
	files, errResponse := app.readData(data)
	if errResponse != nil {
		app.RenderErrorResponse(writer, errResponse.Status, errResponse.Error, errResponse.Message)
	}

	// database process
	response, errResponse := app.filesRepo.InsertFiles(files, parentID)
	if errResponse != nil {
		app.RenderErrorResponse(writer, errResponse.Status, errResponse.Error, errResponse.Message)
		return
	}

	// render output
	app.RenderJSON(writer, http.StatusOK, response)
}

// FindFile finds file from database with id and creates a json response of the data
func (app *App) FindFile(writer http.ResponseWriter, req *http.Request) {
	// read the request param
	logger.Logger().Debug("Finding file")

	params := mux.Vars(req)
	id := params["id"]

	// database process
	response, errResponse := app.filesRepo.FindFile(id)
	if errResponse != nil {
		app.RenderErrorResponse(writer, errResponse.Status, errResponse.Error, errResponse.Message)
		return
	}

	// render output
	app.RenderJSON(writer, http.StatusOK, response)

}

// FindFiles finds all files from database with parent id and creates a json response of the data
func (app *App) FindFiles(writer http.ResponseWriter, req *http.Request) {
	// read the request param
	logger.Logger().Debug("Finding files")

	params := mux.Vars(req)
	parentID := params["parent_id"]

	// database process
	response, errResponse := app.filesRepo.FindFiles(parentID)
	if errResponse != nil {
		app.RenderErrorResponse(writer, errResponse.Status, errResponse.Error, errResponse.Message)
		return
	}

	// render output
	app.RenderJSON(writer, http.StatusOK, response)

}

// DeleteFile deletes project item from database and creates a json response of the data
func (app *App) DeleteFile(writer http.ResponseWriter, req *http.Request) {
	// read the request param
	params := mux.Vars(req)
	id := params["id"]

	// database process
	errResponse := app.filesRepo.DeleteFileWithID(id)
	if errResponse != nil {
		app.RenderErrorResponse(writer, errResponse.Status, errResponse.Error, errResponse.Message)
		return
	}

	// render output
	app.RenderJSON(writer, http.StatusOK, "Delete File successful")
}

// DeleteFiles deletes project item from database and creates a json response of the data
func (app *App) DeleteFiles(writer http.ResponseWriter, req *http.Request) {
	// read the request param
	params := mux.Vars(req)
	parentID := params["parent_id"]

	// database process
	errResponse := app.filesRepo.DeleteFilesWithParent(parentID)
	if errResponse != nil {
		app.RenderErrorResponse(writer, errResponse.Status, errResponse.Error, errResponse.Message)
		return
	}

	// render output
	app.RenderJSON(writer, http.StatusOK, "Delete File successful")
}

// readData reading data from multipart and creates a array byte response of the file
func (app *App) readData(files []*multipart.FileHeader) ([]store.Document, *dto.ErrorResponse) {
	var filesData []store.Document
	for _, f := range files {

		// open file
		file, err := f.Open()
		if err != nil {
			return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed to open file"}
		}

		// read file
		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			if file != nil {
				file.Close()
			}
			return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed to read file"}
		}

		filesData = append(filesData, store.Document{Name: f.Filename, Data: fileData})
		if file != nil {
			file.Close()
		}
	}
	return filesData, nil
}

