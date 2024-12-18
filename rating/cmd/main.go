package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/TylerAldrich814/MetaMovies/common"
	"github.com/TylerAldrich814/MetaMovies/common/endpoint"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery/consul"
	controller "github.com/TylerAldrich814/MetaMovies/rating/internal/controller/rating"
	httphandler "github.com/TylerAldrich814/MetaMovies/rating/internal/handler/http"
	"github.com/TylerAldrich814/MetaMovies/rating/internal/repository/memory"

	_ "github.com/joho/godotenv/autoload"
)

var (
  serviceName = "rating"
  consulAddr  = "localhost:8500"
  httpPort = common.EnvString("HTTP_PORT", ":8082")
)

func main(){
  log.Printf(" ->> RATING SERVICE <<- ")

  defer func(){
    if r := recover(); r != nil {
      log.Printf("Recovered from panic: %v", r)
    }
  }()
  ctx, cancel := signal.NotifyContext(
    context.Background(),
    os.Interrupt,
  )
  defer cancel()

  var port int

  flag.IntVar(&port, "port", 8082, "API Handler Port")
  flag.Parse()
  log.Printf("Starting the Rating Service on port %d\n", port)

  addr := fmt.Sprintf("localhost:%d", port)

  registry, err := consul.NewRegistry(consulAddr)
  if err != nil {
    panic(fmt.Sprintf(
      "->> Failed to create a new Rating Service Consul Registry:: %v\n",
      err,
    ))
  }
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
      select {
      case <-ctx.Done():
      default:
        if err := registry.ReportHealthyState(
          instanceID,
          serviceName,
        ); err != nil {
          log.Printf("Failed to report healthy State: %s\n", err.Error())
        }
        time.Sleep(1 * time.Second)
      }
    }
  }()

  repo := memory.New()
  svc  := controller.New(repo)
  h    := httphandler.New(svc)

  ch := make(chan error, 1)
  go func(){
    if err := h.Handle(
      endpoint.RatingEndpoint,
      h.HandleReq,
    ); err != nil {
      ch<-fmt.Errorf(
        "Failed to handle HTTP endpoint \"%s\"",
        endpoint.MovieEndpoint.String(),
        err,
      )
    }

    if err := h.ListenAndServe(
      fmt.Sprintf(":%d", port),
      nil,
    ); err != nil {
      ch<-fmt.Errorf(
        "->> Failed to Create Movie HTTP Server: %v\n",
        err,
      )
    }
  }()

  select {
  case err := <-ch:
    panic(err)
  case <-ctx.Done():
    log.Printf("->> GRAEFULLY SHUTTING DOWN")
    registry.Deregister(
      ctx,
      instanceID,
      serviceName,
    )
  }
}
