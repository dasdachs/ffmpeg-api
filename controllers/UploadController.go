package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const MAX_UPLOAD_SIZE = 10 * 1024 * 1024 // 2 MB

// ConvertFile godoc
// @Summary Upload and convert file
// @Description Upload file and convert using
// @Accept mpfd
// @Success 201 {string} string "created"
// @Failure 500 {string} string "fail"
// @Router /api/v1/convet [post]
func UploadController(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			log.Printf("Could not parse multipart form: %v\n", err)
			http.Error(w, "Could not process data", http.StatusInternalServerError)
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "File not correct", http.StatusBadRequest)
			return
		}

		defer file.Close()

		fileSize := fileHeader.Size

		log.Printf("File size (bytes): %v\n", fileSize)

		if fileSize > MAX_UPLOAD_SIZE {
			http.Error(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "INVALID_FILE", http.StatusBadRequest)
		}
		newPath := filepath.Join("temp", fileHeader.Filename)
		log.Printf(newPath)

		newFile, err := os.Create(newPath)
		if err != nil {
			http.Error(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil {
			http.Error(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, writeErr := w.Write([]byte("File saved to disk"))

		if writeErr != nil {
			log.Printf("Could not send response")
		}

	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
