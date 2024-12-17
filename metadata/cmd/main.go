package main

import (
	"fmt"
	"log"

	"github.com/TylerAldrich/MetaMovies/common"
	"github.com/TylerAldrich/MetaMovies/common/endpoint"
	"github.com/TylerAldrich/MetaMovies/metadata/internal/controller/metadata"
	httphandler "github.com/TylerAldrich/MetaMovies/metadata/internal/handler/http"
	"github.com/TylerAldrich/MetaMovies/metadata/internal/repository/memory"

	_ "github.com/joho/godotenv/autoload"
)

var (
  httpPort = common.EnvString("HTTP_PORT", ":8081")
)

func main(){
  log.Println(" ->> Starting the movie metadata service @ %s..", httpPort)
  repo := memory.New()
  ctrl := metadata.New(repo)
  h    := httphandler.New(ctrl)

  if err := h.Handle(
    endpoint.MetadataEndpoint,
    h.GetMetadata,
  ); err != nil {
    log.Fatalf(fmt.Sprintf(
      "Failed to handle HTTP endpoint \"%s\" - %v\n",
      endpoint.MetadataEndpoint.String(),
      err,
    ))
  }
  if err := h.ListenAndServe(
    httpPort,
    nil,
  ); err != nil {
    panic(err)
  }
}
