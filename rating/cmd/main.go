package main

import (
	"log"

	"github.com/TylerAldrich/MetaMovies/common"
	"github.com/TylerAldrich/MetaMovies/common/endpoint"
	"github.com/TylerAldrich/MetaMovies/rating/internal/controller/rating"
	httphandler "github.com/TylerAldrich/MetaMovies/rating/internal/handler/http"
	"github.com/TylerAldrich/MetaMovies/rating/internal/repository/memory"

	_ "github.com/joho/godotenv/autoload"
)

var (
  httpPort = common.EnvString("HTTP_PORT", ":8082")
)

func main(){
  log.Println(" ->> Starting the rating service @ %s..", httpPort)

  repo := memory.New()
  ctrl := rating.New(repo)
  h    := httphandler.New(ctrl)

  if err := h.Handle(
    endpoint.RatingEndpoint,
    h.HandleReq,
  ); err != nil {
    log.Fatal("Failed to start new rating HTTP handler: %v", err)
  }
  if err := h.ListenAndServe(
    httpPort,
    nil,
  ); err != nil {
    panic(err)
  }
}
