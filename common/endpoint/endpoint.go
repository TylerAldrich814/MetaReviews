package endpoint

type Endpoint uint

const (
  MetadataEndpoint Endpoint = iota
  RatingEndpoint
  MovieEndpoint
)
var epToString = map[Endpoint]string{
  MetadataEndpoint : "/metadata",
  RatingEndpoint   : "/rating",
  MovieEndpoint    : "/movie",
}
var epFromString = map[string]Endpoint {
  "/metadata" : MetadataEndpoint,
  "/rating"   : RatingEndpoint,
  "/movie"    : MovieEndpoint,
}

// Converts an Endpoint into it's string equivalent
func(ep Endpoint) String() string{
  epstring, ok := epToString[ep]
  if !ok {
    return "unknown"
  }
  return epstring
}

// Converts a string into a valid Endpoint
func FromString(epstring string) *Endpoint {
  endpoint, ok := epFromString[epstring]
  if !ok {
    return nil
  }
  return &endpoint
}
