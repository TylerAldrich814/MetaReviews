package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/TylerAldrich814/MetaMovies/common"
	"github.com/TylerAldrich814/MetaMovies/common/endpoint"
	"github.com/TylerAldrich814/MetaMovies/metadata/internal/controller/metadata"
	httphandler "github.com/TylerAldrich814/MetaMovies/metadata/internal/handler/http"
	"github.com/TylerAldrich814/MetaMovies/metadata/internal/repository/memory"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery/consul"

	_ "github.com/joho/godotenv/autoload"
)

var (
  serviceName = "metadata"
  consulAddr  = "localhost:8500"
  httpPort    = common.EnvString("HTTP_PORT", ":8081")
)

func main(){
  var port int
  
  flag.IntVar(&port, "port", 8081, "API Habdler Port")
  flag.Parse()

  log.Printf("Starting the metadata service on port %d\n", port)

  registry, err := consul.NewRegistry(consulAddr)
  if err != nil {
    panic(fmt.Sprintf(
      "->> Failed to create a new Metadata Service Consul Registry:: %v\n",
      err,
    ))
  }
  ctx := context.Background()

  instanceID := discovery.GenerateInstanceID(serviceName)
  if err := registry.Register(
    ctx, 
    instanceID,
    serviceName,
    fmt.Sprintf(
      "localhost:%s",
      port,
    ),
  ); err != nil {
    panic(fmt.Sprintf(
      "->> Failed to register Metadata Service with Consul Registery: %v\n",
      err,
    ))
  }

  // ->> Metadata Service Health-Check Go Routine 
  go func(){
    for {
      if err := registry.ReportHealthyState(
        instanceID,
        serviceName,
      ); err != nil {
        log.Printf("Failed to report healthy state: %s\n", err.Error())
      }
      time.Sleep(1 * time.Second)
    }
  }()
  defer registry.Deregister(
    ctx,
    instanceID,
    serviceName,
  )

  repo := memory.New()
  svc  := metadata.New(repo)
  h    := httphandler.New(svc)

  if err := h.Handle(
    endpoint.MetadataEndpoint,
    h.GetMetadata,
  ); err != nil {
    log.Fatalf(
      "Failed to handle HTTP endpoint \"%s\" - %v\n",
      endpoint.MetadataEndpoint.String(),
      err,
    )
  }
  if err := h.ListenAndServe(
    httpPort,
    nil,
  ); err != nil {
    panic(err)
  }
}
