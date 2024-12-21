package grpcutil

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"

  "github.com/TylerAldrich814/MetaReviews/pkg/discovery"
)

// ServiceConnection attempts to select a random service instance.
// When successful, we return a gRPC connection to the requested
// service.
func ServiceConnection(
  ctx         context.Context,
  serviceName string,
  registry    discovery.Registry,
)( *grpc.ClientConn, error){
  addrs, err := registry.ServiceAddresses(
    ctx, serviceName,
  )
  if err != nil {
    return nil, err
  }

  return grpc.Dial(
    addrs[rand.Intn(len(addrs))],
    grpc.WithTransportCredentials(insecure.NewCredentials()),
  )
}
