package movie

import (
	"context"
	"errors"

	metadatamodel "github.com/TylerAldrich814/MetaMovies/metadata/pkg/model"
	"github.com/TylerAldrich814/MetaMovies/movie/internal/gateway"
	"github.com/TylerAldrich814/MetaMovies/movie/pkg/model"
	ratingmodel "github.com/TylerAldrich814/MetaMovies/rating/pkg/model"
)

var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
  GetAggregatedRating(context.Context, ratingmodel.RecordID, ratingmodel.RecordType)( float64,error )
  PutRating(context.Context, ratingmodel.RecordID, ratingmodel.RecordType, *ratingmodel.Rating) error
}

type metadataGateway interface {
  Get(ctx context.Context, id string)( *metadatamodel.Metadata,error )
}

// Controller defines a movie service controller.
type Controller struct {
  ratingGateway   ratingGateway
  metadataGateway metadataGateway
}

func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
  return &Controller{ ratingGateway, metadataGateway }
}

func(c *Controller) Get(
  ctx context.Context,
  id  string,
)( *model.MovieDetails, error ){
  metadata, err := c.metadataGateway.Get(ctx, id)
  if err != nil {
    if errors.Is(err, gateway.ErrNotFound){
      return nil, ErrNotFound
    }
    return nil, err
  }
  details := &model.MovieDetails{ Metadata: *metadata }
  rating, err := c.ratingGateway.GetAggregatedRating(
    ctx,
    ratingmodel.RecordID(id),
    ratingmodel.RecordTypeMove,
  )
  if err != nil {
    if !errors.Is(err, gateway.ErrNotFound){
      return details, nil
    }
    return nil, err
  }
  details.Rating = &rating

  return details, nil
}
