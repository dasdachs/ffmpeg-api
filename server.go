package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const PORT = "8080"

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
	log.SetPrefix("[FFMPEG servie] ")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			log.Printf("POST request")

			r.ParseForm()

			fmt.Fprintf(w, "Got the post request", r.Body)
		} else {
			fmt.Fprintf(w, "Only post methods allowed")
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
