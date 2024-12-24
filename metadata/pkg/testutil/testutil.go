package testutil

import (
	"github.com/TylerAldrich814/MetaReviews/common/gen"
	"github.com/TylerAldrich814/MetaReviews/metadata/internal/controller/metadata"
	grpchandler "github.com/TylerAldrich814/MetaReviews/metadata/internal/handler/grpc"
	"github.com/TylerAldrich814/MetaReviews/metadata/internal/repository/memory"
)

// NewTestMetadataGRCServer creates a new metadata gRPC server to be used in tests.
func NewTestMetadataGRCServer() gen.MetadataServiceServer {
  r := memory.New()
  ctrl := metadata.New(r)
  return grpchandler.New(ctrl)
}
