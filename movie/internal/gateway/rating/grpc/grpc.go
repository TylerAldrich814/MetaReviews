package grpc

import (
	"context"

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
  resp, err := client.GetAggregatedRating(
    ctx,
    &gen.GetAggregatedRatingRequest{
      RecordId   : string(recordID),
      RecordType : string(recordType),
    },
  )
  if err != nil {
    return 0, err
  }

  return resp.RatingValue, nil
}
