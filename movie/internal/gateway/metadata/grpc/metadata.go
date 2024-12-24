package grpc

import (
	"context"
	"time"

	"github.com/TylerAldrich814/MetaReviews/common/gen"
	"github.com/TylerAldrich814/MetaReviews/common/grpcutil"
	"github.com/TylerAldrich814/MetaReviews/metadata/pkg/model"
	"github.com/TylerAldrich814/MetaReviews/pkg/discovery"
)

// Gateway defines a movie metadata gRPC Gateway.
type Gateway struct {
  registry discovery.Registry
}

// New creates a new gRPC Gateway for a Movie Metadata Service.
func New(registry discovery.Registry) *Gateway{
  return &Gateway{ registry }
}

// Put puts a movie metadata into the metadata repository.
func(g *Gateway) Put(
  ctx context.Context,
  m   *model.Metadata,
)( error ){
  conn, err := grpcutil.ServiceConnection(
    ctx,
    "metadata",
    g.registry,
  )
  if err != nil {
    return err
  }
  defer conn.Close()

  client := gen.NewMetadataServiceClient(conn)

  maxRetries := 5
  _, err = grpcutil.DoRequestWithBackoff(
    maxRetries,
    time.Duration(100 * time.Microsecond),
    func()( *gen.PutMetadataResponse,error ){
      return client.PutMetadata(ctx, &gen.PutMetadataRequest{
        Metadata: model.MetadataToProto(m),
      })
    },
    grpcutil.ShouldRetry,
  )

  return err
}

// Get returns Movie Metadata, queried via the movie ID.
func(g *Gateway) Get(
  ctx context.Context,
  id  string,
)( *model.Metadata,error ){
  conn, err := grpcutil.ServiceConnection(
    ctx,
    "metadata",
    g.registry,
  )
  if err != nil {
    return nil, err
  }
  defer conn.Close()

  client := gen.NewMetadataServiceClient(conn)

  maxRetries := 5
  resp, err := grpcutil.DoRequestWithBackoff(
    maxRetries,
    time.Duration(100 * time.Millisecond),
    func()( *gen.GetMetadataResponse,error ){
      return client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
    },
    grpcutil.ShouldRetry,
  )
  if err != nil {
    return nil, err
  }

  return model.MetadataFromProto(resp.Metadata), nil
}
