package testutil

import (
	"github.com/TylerAldrich814/MetaReviews/common/gen"
	"github.com/TylerAldrich814/MetaReviews/movie/internal/controller/movie"
	metadatagateway "github.com/TylerAldrich814/MetaReviews/movie/internal/gateway/metadata/grpc"
	ratinggateway "github.com/TylerAldrich814/MetaReviews/movie/internal/gateway/rating/grpc"
	grpchandler "github.com/TylerAldrich814/MetaReviews/movie/internal/handler/grpc"
	"github.com/TylerAldrich814/MetaReviews/pkg/discovery"
)

// NewTestMovieGRPCServer creates a new movie gRPC server to be used in tests.
func NewTestMovieGRPCServer(
  registry discovery.Registry,
) gen.MovieServiceServer {
  ratingGateway := ratinggateway.New(registry)
  metadataGateway := metadatagateway.New(registry)
  ctrl := movie.New(ratingGateway, metadataGateway)

  return grpchandler.New(ctrl)
}
