package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/TylerAldrich814/MetaMovies/common"
	"github.com/TylerAldrich814/MetaMovies/common/gen"
	"github.com/TylerAldrich814/MetaMovies/movie/internal/controller/movie"
	metadatagateway "github.com/TylerAldrich814/MetaMovies/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/TylerAldrich814/MetaMovies/movie/internal/gateway/rating/http"
	grpchandler "github.com/TylerAldrich814/MetaMovies/movie/internal/handler/grpc"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery/consul"
	"google.golang.org/grpc"

	_ "github.com/joho/godotenv/autoload"
)

var (
  serviceName  = "movie"
  consulAddr  = "localhost:8500"
  metadataAddr = common.EnvString("METADATA_ADDR", "localhost:8081")
  ratingAddr   = common.EnvString("RATING_ADDR", "localhost:8082")
  movieAddr    = common.EnvString("MOVIE_ADDR", "localhost:8083")
)

func main(){
  log.Printf(" ->> MOVIE SERVICE <<- ")
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

  flag.IntVar(&port, "port", 8083, "API Handler Port")
  flag.Parse()
  log.Printf("Starting the Movie Service on port %d\n", port)

  addr := fmt.Sprintf("localhost:%d", port)

  registry, err := consul.NewRegistry(consulAddr)
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
