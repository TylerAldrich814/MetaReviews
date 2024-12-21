package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/TylerAldrich814/MetaReviews/common/endpoint"
	"github.com/TylerAldrich814/MetaReviews/metadata/pkg/model"
	"github.com/TylerAldrich814/MetaReviews/movie/internal/gateway"
	"github.com/TylerAldrich814/MetaReviews/pkg/discovery"
)

// Gateway defines a movie metadata HTTP gateway.
type Gateway struct {
  registry discovery.Registry
}

// New creates a new HTTP gateway for a movie metadata
func New(registry discovery.Registry) *Gateway {
  return &Gateway{ registry  }
}

func(g *Gateway) getURL(
  ctx context.Context,
)( string,error ){
  addrs, err := g.registry.ServiceAddresses(
    ctx,
    endpoint.MetadataEndpoint.String(),
  )
  if err != nil {
    return "", err
  }
  if len(addrs) == 0 {
    return "", discovery.ErrNotFound
  }
  return fmt.Sprintf(
    "http://%s/metadata",
    addrs[rand.Intn(len(addrs))],
  ), nil
}

// Get receives movie metadata by a movie's ID
func(g *Gateway) Get(
  ctx  context.Context,
  id   string,
)( *model.Metadata, error ){
  url, err := g.getURL(ctx)
  if err != nil {
    return nil, err
  }

  log.Printf("Calling metadata service. Request: GET %s\n", url)
  req, err := http.NewRequest(http.MethodGet, url, nil)
  if err != nil {
    return nil, err
  }

  req = req.WithContext(ctx)
  values := req.URL.Query()
  values.Add("id", id)
  req.URL.RawQuery = values.Encode()
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  if resp.StatusCode == http.StatusNotFound {
    return nil, gateway.ErrNotFound
  }
  if resp.StatusCode/100 != 2 {
    return nil, fmt.Errorf("non-2xx response: %v", resp)
  }

  var v *model.Metadata
  if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
    return nil, err
  }

  return v, nil
}
