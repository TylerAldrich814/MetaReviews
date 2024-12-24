package grpc

import (
	"context"
	"time"

	"github.com/TylerAldrich814/MetaReviews/common/gen"
	"github.com/TylerAldrich814/MetaReviews/common/grpcutil"
	"github.com/TylerAldrich814/MetaReviews/pkg/discovery"
	"github.com/TylerAldrich814/MetaReviews/rating/pkg/model"
)

// Gateway defines a gRPC Gateway for a Rating Service.
type Gateway struct {
  registry discovery.Registry
}

// New creates a new gRPC Gateway for a Rating Service.
func New(registry discovery.Registry) *Gateway {
  return &Gateway{ registry }
}

// GetAggregatedRating - If successful, returns the aggregated rating for a record.
// Otherwise returns an ErrNotFound if no ratings are found.
func(g *Gateway) GetAggregatedRating(
  ctx        context.Context,
  recordID   model.RecordID,
  recordType model.RecordType,
)( float64,error ){
  conn, err := grpcutil.ServiceConnection(
    ctx,
    "rating",
    g.registry,
  )
  if err != nil {
    return 0, err
  }
  defer conn.Close()

  client := gen.NewRatingServiceClient(conn)

  maxRetries := 5
  resp, err := grpcutil.DoRequestWithBackoff(
    maxRetries,
    time.Duration(100 * time.Millisecond),
    func()( *gen.GetAggregatedRatingResponse,error ){
      return client.GetAggregatedRating(
        ctx,
        &gen.GetAggregatedRatingRequest{
          RecordId   : string(recordID),
          RecordType : string(recordType),
        },
      )
    },
    grpcutil.ShouldRetry,
  )
  if err != nil {
    return 0, err
  }

  return resp.RatingValue, nil
}

// PutRating adds a new rating for a given record into the Rating service repository.
func(g *Gateway) PutRating(
  ctx        context.Context,
  recordID   model.RecordID,
  recordType model.RecordType,
  rating     *model.Rating,
) error {
  conn, err := grpcutil.ServiceConnection(
    ctx,
    "rating",
    g.registry,
  )
  if err != nil {
    return err
  }
  defer conn.Close()

  client := gen.NewRatingServiceClient(conn)

  maxRetries := 5
  _, err = grpcutil.DoRequestWithBackoff(
    maxRetries,
    time.Duration(100 * time.Millisecond),
    func()( *gen.PutRatingResponse,error ){
      return client.PutRating(
        ctx,
        &gen.PutRatingRequest{
          UserId      : string(rating.UserID),
          RecordId    : string(recordID),
          RecordType  : string(recordType),
          RatingValue : int32(rating.Value),
        },
      )
    },
    grpcutil.ShouldRetry,
  )

  return err
}
