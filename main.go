package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/dasdachs/ffmpeg-stream/controllers"
	_ "github.com/dasdachs/ffmpeg-stream/docs"
	"github.com/dasdachs/ffmpeg-stream/utils"
	// "github.com/giorgisio/goav/avformat"
)

const PORT = "8080"

// @title FFMPEG Server
// @version 1.0
// @description FFMPEG rest api server
// @contact.name Jani Šumak
// @contact.email jani.sumak@gmila.com
// @license.name MIT
// @license.url https://mit-license.org/
// @host localhost:8080
// @BasePath /api/v1
func main() {
	utils.ParseEnv()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// /swagger/index.html
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", port)),
	))

	// API
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/convert", controllers.UploadController)
	})

	fmt.Printf("Server started on %s\n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatal(err)
	}
}
