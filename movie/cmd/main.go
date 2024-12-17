package main

import (
	"log"

	"github.com/TylerAldrich814/MetaMovies/common"
	"github.com/TylerAldrich814/MetaMovies/common/endpoint"
	"github.com/TylerAldrich814/MetaMovies/movie/internal/controller/movie"
	metadatagateway "github.com/TylerAldrich814/MetaMovies/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/TylerAldrich814/MetaMovies/movie/internal/gateway/rating/http"
	httphandler "github.com/TylerAldrich814/MetaMovies/movie/internal/handler/http"

	_ "github.com/joho/godotenv/autoload"
)

var (
  metadataAddr = common.EnvString("METADATA_ADDR", "localhost:8081")
  ratingAddr   = common.EnvString("RATING_ADDR", "localhost:8082")
  movieAddr    = common.EnvString("MOVIE_ADDR", "localhost:8083")
)

func main(){
  log.Printf(" ->> Starting the Movie Service @ %s..\n", movieAddr)

  metadataGateway := metadatagateway.New(metadataAddr)
  ratingGateway   := ratinggateway.New(ratingAddr)
  ctrl := movie.New(ratingGateway, metadataGateway)
  h := httphandler.New(ctrl)

  if err := h.Handle(
    endpoint.MovieEndpoint,
    h.GetMoviewDetails,
  ); err != nil {
    panic(err)
  }
  if err := h.ListenAndServe(
    movieAddr,
    nil,
  ); err != nil {
    panic(err)
  }
}
