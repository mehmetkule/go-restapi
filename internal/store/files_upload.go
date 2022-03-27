package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/mehmetkule/go-restapi/internal/dto"
	"github.com/mehmetkule/go-restapi/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// FilesRepo Struct
type FilesRepo struct {
	DB *sql.DB
}

type Document struct {
	Name string
	Data []byte
}


func (r *FilesRepo) InsertFiles(data []Document, parentID string) (*dto.FilesResponse, *dto.ErrorResponse) {
	var insertedID []uuid.UUID
	for _, file := range data {
		sql := "INSERT INTO document(parent_id,name,data,created) VALUES($1,$2,$3,$4) returning id;"
		row := r.DB.QueryRowContext(context.Background(), sql, parentID, file.Name, file.Data, time.Now())
		var lastInsertID uuid.UUID
		if row.Scan(&lastInsertID) != nil {
			return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: nil, Message: "Some went wrong"}
		}
		insertedID = append(insertedID, lastInsertID)
	}
	return &dto.FilesResponse{
		ID: insertedID,
	}, nil
}

// FindFile finds file with id from database
func (r *FilesRepo) FindFile(ID string) (*dto.FileResponse, *dto.ErrorResponse) {
	sqlQuery := "SELECT id,parent_id,name,data,created FROM document WHERE id=$1"

	var rows *sql.Rows
	var err error
	logger.Logger().Debug("Finding file",zap.String("ID",ID))

	// select sql message
	if rows, err = r.DB.QueryContext(context.Background(), sqlQuery, ID); err != nil {
		return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: fmt.Sprintf("Failed find file %s", ID)}
	}

	defer rows.Close()

	response := dto.FileResponse{}
	if rows.Next() {
		err = rows.Scan(&response.ID, &response.ParentID, &response.Name, &response.Data, &response.Created)
		if err != nil {
			return nil, &dto.ErrorResponse{Status: http.StatusNotFound, Error: err, Message: fmt.Sprintf("File [%s] not found", ID)}
		}
		return &response, nil
	} else {
		if err = rows.Err(); err != nil {
			return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed to fetch file"}
		} else {
			return nil, &dto.ErrorResponse{Status: http.StatusNotFound, Error: fmt.Errorf("not found"), Message: "file not found"}
		}
	}
}

// FindFiles finds all files with parent id from database
func (r *FilesRepo) FindFiles(parentID string) (*[]dto.FileResponse, *dto.ErrorResponse) {
	sqlQuery := "SELECT id,parent_id,name,data,created FROM document WHERE parent_id=$1"
	var rows *sql.Rows
	var err error

	logger.Logger().Debug("Finding files",zap.String("parentID",parentID))

	// select sql message
	if rows, err = r.DB.Query(sqlQuery, parentID); err != nil {
		return nil, &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed find files"}
	}

	defer rows.Close()

	response := []dto.FileResponse{}
	for rows.Next() {
		rowResult := dto.FileResponse{}
		err = rows.Scan(&rowResult.ID, &rowResult.ParentID, &rowResult.Name, &rowResult.Data, &rowResult.Created)
		if err != nil {
			return nil, &dto.ErrorResponse{Status: http.StatusNotFound, Error: err, Message: fmt.Sprintf("Files [%s] not found", parentID)}
		}
		response = append(response, rowResult)
	}

	return &response, nil
}

// DeleteFileWithID deletes all files with id
func (r *FilesRepo) DeleteFileWithID(id string) *dto.ErrorResponse {
	_, err := r.DB.Exec("DELETE FROM document WHERE id=$1;", id)
	if err != nil {
		return &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed delete file"}
	}
	return nil
}

// DeleteFilesWithParent deletes all files with parent id
func (r *FilesRepo) DeleteFilesWithParent(parentID string) *dto.ErrorResponse {
	_, err := r.DB.Exec("DELETE FROM document WHERE parent_id=$1;", parentID)
	if err != nil {
		return &dto.ErrorResponse{Status: http.StatusInternalServerError, Error: err, Message: "Failed delete files"}
	}
	return nil
}