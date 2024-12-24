package testutil

import (
	"github.com/TylerAldrich814/MetaReviews/common/gen"
	"github.com/TylerAldrich814/MetaReviews/rating/internal/controller/rating"
	"github.com/TylerAldrich814/MetaReviews/rating/internal/repository/memory"
	grpchandler "github.com/TylerAldrich814/MetaReviews/rating/internal/handler/grpc"
)

// NewTestRatingGRPCServer creates a new rating gRPC server to be used in tests.
func NewTestRatingGRPCServer() gen.RatingServiceServer {
  r := memory.New()
  ctrl := rating.New(r, nil)

  return grpchandler.New(ctrl)
}
