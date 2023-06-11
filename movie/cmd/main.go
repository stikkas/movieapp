package main

import (
	"log"
	"net/http"

	"movieexample.com/movie/internal/controller/movie"
	metadatagw "movieexample.com/movie/internal/gateway/metadata/http"
	ratinggw "movieexample.com/movie/internal/gateway/rating/http"
	httphandler "movieexample.com/movie/internal/handler/http"
)

func main() {
	log.Println("Starting the movie service")
	metadataGateway := metadatagw.New("http://localhost:8081")
	ratingGateway := ratinggw.New("http://localhost:8082")
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	panic(http.ListenAndServe(":8083", nil))
}
