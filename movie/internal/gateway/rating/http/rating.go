package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	// "github.com/TylerAldrich814/MetaReviews/common/endpoint"
	"github.com/TylerAldrich814/MetaReviews/movie/internal/gateway"
	"github.com/TylerAldrich814/MetaReviews/pkg/discovery"
	"github.com/TylerAldrich814/MetaReviews/rating/pkg/model"
)

// Gateway defines an HTTP gateway for a rating service
type Gateway struct {
  registry discovery.Registry
}

// New creates a new HTTP gateway for a rating service.
func New(registry discovery.Registry) *Gateway {
  return &Gateway{ registry  }
}

func(g *Gateway) getURL(
  ctx context.Context,
)( string,error ){
  addrs, err := g.registry.ServiceAddresses(
    ctx,
    "review",
  )
  if err != nil {
    return "", err
  }
  if len(addrs) == 0 {
    return "", discovery.ErrNotFound
  }
  return fmt.Sprintf(
    "http://%s/review",
    addrs[rand.Intn(len(addrs))],
  ), nil
}

// GetAggregatedRating returns the aggregated rating for a record
// or ErrNotFound if there are no ratings for it.
func(g *Gateway) GetAggregatedRating(
  ctx        context.Context,
  recordID   model.RecordID,
  recordType model.RecordType,
)( float64, error) {
  addrs, err := g.registry.ServiceAddresses(
    ctx,
    "rating",
  )
  if err != nil {
    return 0, nil
  }
  url := fmt.Sprintf(
    "http://%s/rating",
    addrs[rand.Intn(len(addrs))],
  )
  log.Printf("Calling rating service. Request: GET %s\n", url)
  req, err := http.NewRequest(http.MethodGet, url, nil)
  if err != nil {
    return 0, nil
  }

  req = req.WithContext(ctx)
  values := req.URL.Query()
  values.Add("id", string(recordID))
  values.Add("type", fmt.Sprintf("%v", recordType))

  req.URL.RawQuery = values.Encode()
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return 0, err
  }
  defer resp.Body.Close()

  if resp.StatusCode == http.StatusNotFound {
    return 0, gateway.ErrNotFound
  }
  if resp.StatusCode/100 != 2 {
    return 0, fmt.Errorf("non-2xx response: %v", resp)
  }
  var v float64
  if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
    return 0, err
  }

  return v, nil
}

// PutRating writes a rating.
func(g *Gateway) PutRating(
  ctx        context.Context,
  recordID   model.RecordID,
  recordType model.RecordType,
  rating     *model.Rating,
) error {
  url, err := g.getURL(ctx)
  if err != nil {
    return err
  }

  req, err := http.NewRequest(
    http.MethodPut,
    url,
    nil,
  )
  if err != nil {
    return err
  }
  req = req.WithContext(ctx)
  values := req.URL.Query()
  values.Add("id", string(recordID))
  values.Add("type", fmt.Sprintf("%v", recordType))
  values.Add("userId", string(rating.UserID))
  values.Add("value", fmt.Sprintf("%v", rating.Value))
  
  req.URL.RawQuery = values.Encode()
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  if resp.StatusCode/100 != 2 {
    return fmt.Errorf("non-2xx response: %v", err)
  }

  return nil
}
