package grpc

import (
	"context"
	"errors"

	"github.com/TylerAldrich814/MetaReviews/common/gen"
	"github.com/TylerAldrich814/MetaReviews/rating/internal/controller/rating"
	"github.com/TylerAldrich814/MetaReviews/rating/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a gRPC rating API Handler.
type Handler struct {
  gen.UnimplementedRatingServiceServer
  svc *rating.Controller
}

// New creates a new Rating gRPC Handler. 
func New(svc *rating.Controller) *Handler {
  return &Handler{ svc:svc }
}

// GetAggregatedRating returns the aggregated rating for a record.
func (h *Handler) GetAggregatedRating(
  ctx context.Context,
  req *gen.GetAggregatedRatingRequest,
)( *gen.GetAggregatedRatingResponse,error ){
  if req == nil {
    return nil, status.Errorf(
      codes.InvalidArgument,
      "nil is req",
    )
  }
  if req.RecordId == "" {
    return nil, status.Errorf(
      codes.InvalidArgument,
      "missing record id",
    )
  }

  v, err := h.svc.GetAggregatedRating(
    ctx,
    model.RecordID(req.RecordId),
    model.RecordType(req.RecordType),
  )
  if err != nil {
    if errors.Is(err, rating.ErrNotFound) {
      return nil, status.Errorf(
        codes.NotFound,
        err.Error(),
      )
    }
    return nil, status.Errorf(
      codes.Internal,
      err.Error(),
    )
  }

  return &gen.GetAggregatedRatingResponse{
    RatingValue: v,
  },nil
}

// PutRating writes a rating for a given record.
func(h *Handler) PutRating(
  ctx context.Context,
  req *gen.PutRatingRequest,
)( *gen.PutRatingResponse,error ){
  errMsg := ""
  if req == nil {
    return nil, status.Errorf(
      codes.InvalidArgument,
      "req is nil",
    )
  }
  if req.RecordId == "" {
    errMsg = "missing record id,"
  }
  if req.UserId == "" {
    errMsg += "missing user id,"
  }
  if errMsg != "" {
    return nil, status.Errorf(
      codes.InvalidArgument,
      errMsg,
    )
  }

  if err := h.svc.PutRating(
    ctx,
    model.RecordID(req.RecordId),
    model.RecordType(req.RecordType),
    &model.Rating{
      UserID : model.UserID(req.UserId),
      Value  : model.RatingValue(req.RatingValue),
    },
  ); err != nil {
    return nil, status.Errorf(
      codes.Internal,
      err.Error(),
    )
  }

  return &gen.PutRatingResponse{}, nil
}
