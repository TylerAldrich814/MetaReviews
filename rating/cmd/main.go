package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/TylerAldrich814/MetaMovies/common"
	"github.com/TylerAldrich814/MetaMovies/common/endpoint"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery/consul"
	"github.com/TylerAldrich814/MetaMovies/rating/internal/controller/rating"
	httphandler "github.com/TylerAldrich814/MetaMovies/rating/internal/handler/http"
	"github.com/TylerAldrich814/MetaMovies/rating/internal/repository/memory"
  controller "github.com/TylerAldrich814/MetaMovies/rating/internal/controller/rating"

	_ "github.com/joho/godotenv/autoload"
)

var (
  serviceName = "rating"
  consulAddr  = "localhost:8500"
  httpPort = common.EnvString("HTTP_PORT", ":8082")
)

func main(){
  var port int

  flag.IntVar(&port, "port", 8082, "API Handler Port")
  flag.Parse()
  log.Printf("Starting the Rating Service on port %d\n", port)

  addr := fmt.Sprint("localhost:%d", port)

  registry, err := consul.NewRegistry(consulAddr)
  if err != nil {
    panic(fmt.Sprintf(
      "->> Failed to create a new Rating Service Consul Registry:: %v\n",
      err,
    ))
  }
  ctx := context.Background()
  instanceID := discovery.GenerateInstanceID(serviceName)
  if err := registry.Register(
    ctx,
    instanceID,
    serviceName,
    addr,
  ); err != nil {
    panic(fmt.Sprintf(
      "->> Failed to register Rating Service with Consul Registery: %v\n",
      err,
    ))
  }
  // ->> Rating Service Health-Check Go Routine 
  go func(){
    for {
      if err := registry.ReportHealthyState(
        instanceID,
        serviceName,
      ); err != nil {
        log.Printf("Failed to report healthy State: %s\n", err.Error())
      }
      time.Sleep(1 * time.Second)
    }
  }()
  defer registry.Deregister(ctx, instanceID, serviceName)

  repo := memory.New()
  svc  := controller.New(repo)
  h    := httphandler.New(svc)

  if err := h.Handle(
    endpoint.RatingEndpoint,
    h.HandleReq,
  ); err != nil {
    panic(fmt.Sprintf(
      "Failed to handle HTTP endpoint \"%s\"",
      endpoint.MovieEndpoint.String(),
      err,
    ))
  }

  if err := h.ListenAndServe(
    fmt.Sprintf(":%d", port),
    nil,
  ); err != nil {
    panic(fmt.Sprintf(
      "->> Failed to Create Move HTTP Server: %v\n",
      err,
    ))
  }
}

func oldmain(){
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
