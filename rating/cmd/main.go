package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/TylerAldrich814/MetaMovies/common"
	"github.com/TylerAldrich814/MetaMovies/common/gen"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery/consul"
	"github.com/TylerAldrich814/MetaMovies/rating/internal/controller/rating"
	grpcHandler "github.com/TylerAldrich814/MetaMovies/rating/internal/handler/grpc"
	"github.com/TylerAldrich814/MetaMovies/rating/internal/repository/mysql"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

var (
  serviceName = "rating"
  consulAddr  = "localhost:8500"
  httpPort = common.EnvString("HTTP_PORT", ":8082")
)

func main(){
  log.Printf(" ->> RATING SERVICE <<- ")

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

  repo, err := mysql.New()
  if err != nil {
    panic(err)
  }
  ctrl := rating.New(repo, nil)
  h    := grpcHandler.New(ctrl)
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
        err,
      )
      return
    }
    srv := grpc.NewServer()
    reflection.Register(srv)

    gen.RegisterRatingServiceServer(srv, h)

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
