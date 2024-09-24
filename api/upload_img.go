package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func UploadImages(w http.ResponseWriter, r *http.Request, filePath string) (string, error) {
	err := r.ParseMultipartForm(20 << 20) // 20 MB max
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return "", err
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		if file == nil {
			http.Error(w, "No file selected", http.StatusBadRequest)
			return "", errors.New("no file selected")
		}
		http.Error(w, "Error retrieving file from form", http.StatusBadRequest)
		return "", err
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file data", http.StatusInternalServerError)
		return "", err
	}

	if len(fileData) == 0 {
		http.Error(w, "No file data received", http.StatusBadRequest)
		return "", errors.New("no file data received")
	}

	fileSize := int64(len(fileData))
	maxFileSize := int64(20 << 20)
	if fileSize > maxFileSize {
		return "", errors.New("taille de l'image trop grande")
	}

	f, err := os.OpenFile(filepath.Join(filePath, handler.Filename), os.O_WRONLY|os.O_CREATE, 0o666)
	fmt.Println(filePath)
	if err != nil {
		http.Error(w, "Error creating file on server", http.StatusInternalServerError)
		return "", err
	}
	defer f.Close()

	_, err = f.Write(fileData)
	if err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return "", err
	}

	return path.Join(".", filePath, handler.Filename), nil
}

