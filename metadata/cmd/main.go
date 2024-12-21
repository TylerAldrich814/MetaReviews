package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/TylerAldrich814/MetaReviews/common/gen"
	"github.com/TylerAldrich814/MetaReviews/metadata/internal/controller/metadata"
	grpcHandler "github.com/TylerAldrich814/MetaReviews/metadata/internal/handler/grpc"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	"github.com/TylerAldrich814/MetaReviews/metadata/internal/repository/memory"
	"github.com/TylerAldrich814/MetaReviews/pkg/discovery"
	"github.com/TylerAldrich814/MetaReviews/pkg/discovery/consul"

	_ "github.com/joho/godotenv/autoload"
)

var (
  serviceName = "metadata"
)

func main(){
  log.Printf(" ->> METADATA SERVICE <<- ")

  f, err := os.Open("configs/base.yaml")
  if err != nil {
    panic(err)
  }
  defer f.Close()

  var cfg serviceConfig
  if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
    panic(err)
  }
  port := cfg.APIConfig.Port

  ctx, cancel := signal.NotifyContext(
    context.Background(),
    os.Interrupt,
  )
  defer cancel()

  log.Printf("Starting the metadata service on port %d\n", port)

  // registry, err := consul.NewRegistry("localhost:8500")
  registry, err := consul.NewRegistry("host.docker.internal:8500")
  if err != nil {
    panic(fmt.Sprintf(
      "->> Failed to create a new Metadata Service Consul Registry:: %v\n",
      err,
    ))
  }

  instanceID := discovery.GenerateInstanceID(serviceName)
  if err := registry.Register(
    ctx, 
    instanceID, 
    serviceName, 
    fmt.Sprintf("localhost:%d", port),
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

  repo := memory.New()
  svc  := metadata.New(repo)
  h    := grpcHandler.New(svc)
  ch   := make(chan error, 1)

  go func(){
    lis, err := net.Listen(
      "tcp", 
      fmt.Sprintf(
        "localhost:%d",
        port,
      ),
    )
    if err != nil {
      ch <- fmt.Errorf(
        "failed to start tcp listener: %v\n",
        err,
      )
      return
    }
    svc := grpc.NewServer()
    gen.RegisterMetadataServiceServer(svc, h)
    if err := svc.Serve(lis); err != nil {
      ch <- fmt.Errorf(
        "failed to serve grpc server: %v\n",
        err,
      )
    }
  }()

  select {
  case err := <-ch:
    registry.Deregister(
      ctx,
      instanceID,
      serviceName,
    )
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
