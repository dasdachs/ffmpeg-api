package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const PORT = "8080"
const MAX_UPLOAD_SIZE = 10 * 1024 * 1024 // 2 MB

func ReadEnvFile() []string {
	_, err := os.Stat(".env")
	if os.IsNotExist(err) {
		return make([]string, 0)
	}

	data, err := ioutil.ReadFile(".env")
	if err != nil {
		log.Panic("Could not parse .env file")
	}

	content := strings.Split(string(data), "\n")

	return content
}

func ParseAnSetEnv(content []string) {
	for _, line := range content {
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		vals := strings.Split(line, "=")

		key := vals[0]
		val := strings.TrimSpace(strings.Split(vals[1], "#")[0])

		err := os.Setenv(key, val)
		if err != nil {
			log.Fatalf("Could not parse %s\n", line)
		}
	}
}

func main() {
	// Parse and set the .env file in development
	envContent := ReadEnvFile()
	ParseAnSetEnv(envContent)

	// Setup simple logger
	log.SetFlags(2 | 3)
	log.SetPrefix("[FFMPEG service] ")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
			w.Write([]byte("File saved to disk"))

		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	// Starting the webserver
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}
