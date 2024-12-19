package grpc

import (
	"context"
	"errors"

	"github.com/TylerAldrich814/MetaMovies/common/gen"
	"github.com/TylerAldrich814/MetaMovies/metadata/internal/controller/metadata"
	"github.com/TylerAldrich814/MetaMovies/metadata/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a Movie Metadata gRPC Handler.
type Handler struct {
  gen.UnimplementedMetadataServiceServer
  ctrl *metadata.Controller
}

// New Creates a new Movie Metadata gRPC Handler.
func New(ctrl *metadata.Controller) *Handler {
  return &Handler{ ctrl:ctrl }
}

// GetMetadataByID returns movie metadata by it's id.
func(h *Handler) GetMetadataByID(
  ctx context.Context,
  req *gen.GetMetadataRequest,
)( *gen.GetMetadataResponse, error){
  if req == nil || req.MovieId == "" {
    return nil, status.Errorf(
      codes.InvalidArgument,
      "nil req or empty id",
    )
  }

  m, err := h.ctrl.Get(
    ctx,
    req.MovieId,
  )
  if err != nil {
    if errors.Is(err, metadata.ErrNotFound){
      return nil, status.Error(
        codes.NotFound, 
        err.Error(),
      )
    }
    return nil, status.Errorf(
      codes.Internal, 
      err.Error(),
    )
  }

  return &gen.GetMetadataResponse{
    Metadata: model.MetadataToProto(m),
  }, nil
}
