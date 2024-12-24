package grpcutil

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

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

  return grpc.NewClient(
    addrs[rand.Intn(len(addrs))],
    grpc.WithTransportCredentials(insecure.NewCredentials()),
  )
}


// ShouldRetry determines if an error is a gRPC related error.
// If so, it then tests is the gRPC error is that of the following
// gRPC error kinds:
//   - DeadlineExceeded
//   - ResourceExhausted
//   - Unavailable
// If the gRPC error is any of the three listest, we then return true.
func ShouldRetry(err error) bool {
  e, ok := status.FromError(err)
  if !ok {
    return false
  }
  errCode := e.Code()
  return errCode == codes.DeadlineExceeded ||
         errCode == codes.ResourceExhausted ||
         errCode == codes.Unavailable
}
