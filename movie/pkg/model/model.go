package model

import "github.com/TylerAldrich814/MetaMovies/metadata/pkg/model"

// MovieDetails includes movie metadata and it's respective aggregated rating.
type MovieDetails struct {
  Rating   *float64       `json:"rating"`
  Metadata model.Metadata `json:"metadata"`
}
