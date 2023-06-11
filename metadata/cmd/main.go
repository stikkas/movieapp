package main

import (
	"context"
	"log"
	"net/http"

	"movieexample.com/metadata/internal/controller/metadata"
	httphandler "movieexample.com/metadata/internal/handler/http"
	"movieexample.com/metadata/internal/repository/memory"
	"movieexample.com/metadata/pkg/model"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	repo.Put(context.Background(), "-1", &model.Metadata{"-1", "Default movie", "Default description", "Unknown"})

	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))

	panic(http.ListenAndServe(":8081", nil))
}
