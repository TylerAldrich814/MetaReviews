package grpc

import (
	"context"

	"github.com/TylerAldrich814/MetaMovies/common/gen"
	"github.com/TylerAldrich814/MetaMovies/common/grpcutil"
	"github.com/TylerAldrich814/MetaMovies/metadata/pkg/model"
	"github.com/TylerAldrich814/MetaMovies/pkg/discovery"
)

// Gateway defines a movie metadata gRPC Gateway.
type Gateway struct {
  registry discovery.Registry
}

// New creates a new gRPC Gateway for a Movie Metadata Service.
func New(registry discovery.Registry) *Gateway{
  return &Gateway{ registry }
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
  resp, err := client.GetMetadata(
    ctx,
    &gen.GetMetadataRequest{
      MovieId: id,
    },
  )
  if err != nil {
    return nil, err
  }

  return model.MetadataFromProto(resp.Metadata), nil
}
