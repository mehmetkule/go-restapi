package dto

import (
	"fmt"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)


type FilesResponse struct {
	ID []uuid.UUID
}


type FileResponse struct {
	ID       uuid.UUID
	ParentID string
	Name     string
	Data     []byte
	Created  time.Time
}

func (f *FileResponse) ValidateFile(req *http.Request) (int, error) {

	data := req.MultipartForm.File["file"]

	if data == nil {
		return http.StatusBadRequest, fmt.Errorf("resim alanı boş bırakılamaz. %v", f.Data)
	}

	if len(data) > 5 {
		return http.StatusBadRequest, fmt.Errorf("5 taneden fazla gönderilemez. %v", len(f.Data))
	}

	return http.StatusOK, nil
}

