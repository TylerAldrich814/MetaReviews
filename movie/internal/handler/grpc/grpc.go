package grpc

import (
	"context"
	"errors"

	"github.com/TylerAldrich814/MetaMovies/common/gen"
	"github.com/TylerAldrich814/MetaMovies/metadata/pkg/model"
	"github.com/TylerAldrich814/MetaMovies/movie/internal/controller/movie"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a Movie gRPC Handler.
type Handler struct {
  gen.UnimplementedMovieServiceServer
  ctrl *movie.Controller
}

func New(ctrl *movie.Controller) *Handler {
  return &Handler{ ctrl:ctrl }
}

func(h *Handler) GetMovieDetails(
  ctx  context.Context,
  req  *gen.GetMovieDetailsRequest,
)( *gen.GetMovieDetailsResponse,error ){
  if req == nil {
    return nil, status.Errorf(
      codes.InvalidArgument,
      "movie details request is missing",
    )
  }
  if req.MovieId == "" {
    return nil, status.Errorf(
      codes.InvalidArgument,
      "movie ID is missing",
    )
  }
  m, err := h.ctrl.Get(ctx, req.MovieId)
  if err != nil {
    if errors.Is(err, movie.ErrNotFound){
      return nil, status.Errorf(
        codes.NotFound,
        "movie ID not found",
      )
    }
    return nil, status.Errorf(
      codes.Internal,
      err.Error(),
    )
  }

  return &gen.GetMovieDetailsResponse{
    MovieDetails : &gen.MovieDetails{
      Rating   : *m.Rating,
      Metadata : model.MetadataToProto(&m.Metadata),
    },
  }, nil
}