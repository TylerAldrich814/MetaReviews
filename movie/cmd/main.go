package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/TylerAldrich814/MetaMovies/common/gen"
	"github.com/TylerAldrich814/MetaMovies/movie/internal/controller/movie"
	metadatagateway "github.com/TylerAldrich814/MetaMovies/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/TylerAldrich814/MetaMovies/movie/internal/gateway/rating/http"
	grpchandler "github.com/TylerAldrich814/MetaMovies/movie/internal/handler/grpc"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery/consul"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	_ "github.com/joho/godotenv/autoload"
)

var (
  serviceName  = "movie"
)

func main(){
  log.Printf(" ->> MOVIE SERVICE <<- ")

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

  addr := fmt.Sprintf("localhost:%d", port)

  registry, err := consul.NewRegistry("host.docker.internal:8500")
  if err != nil {
    panic(fmt.Sprintf(
      "->> Failed to create a new Movie Service Consul Registry:: %v\n",
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
      "->> Failed to register Movie Service with Consul Registery: %v\n",
      err,
    ))
  }
  // ->> Movie Service Health-Check Go Routine 
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

  ratingGateway   := ratinggateway.New(registry)
  metadataGateway := metadatagateway.New(registry)

  ctrl := movie.New(ratingGateway, metadataGateway)
  h := grpchandler.New(ctrl)

  grpcAddr := fmt.Sprintf(
    "localhost:%d",
    port,
  )

  ch := make(chan error, 1)
  go func(){
    lis, err := net.Listen("tcp", grpcAddr)
    if err != nil {
      ch <- fmt.Errorf(
        "failed to start grpc listener: %v\n",
      )
      return
    }
    srv := grpc.NewServer()
    gen.RegisterMovieServiceServer(srv, h)

    if err := srv.Serve(lis); err != nil {
      ch <- fmt.Errorf(
        "failed to listen on rating grpc service server: %v\n",
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
