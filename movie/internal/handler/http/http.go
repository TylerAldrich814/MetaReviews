package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/TylerAldrich814/MetaMovies/common"
	"github.com/TylerAldrich814/MetaMovies/common/endpoint"
	"github.com/TylerAldrich814/MetaMovies/movie/internal/controller/movie"
	"github.com/TylerAldrich814/MetaMovies/movie/internal/handler"
)

// Handler defines a movie handler
type Handler struct {
  ctrl *movie.Controller
}

func(h *Handler) Handle(
  endpoint    endpoint.Endpoint,
  handlerfunc func(http.ResponseWriter, *http.Request),
) error {
  ep := endpoint.String()
  if ep == "unknown" {
    return handler.ErrUnknownEndpoint
  }
  http.Handle(
    "/"+ep,
    http.HandlerFunc(handlerfunc),
  )
  return nil
}

func(h *Handler) ListenAndServe(
  addr    string,
  handler http.Handler,
) error {
  return http.ListenAndServe(
    addr,
    handler,
  )
}

// New creates a new movie HTTP handler.
func New(ctrl *movie.Controller) *Handler {
  return &Handler{ ctrl }
}

// GetMoviewDetails handles new movie GET handler.
func(h *Handler) GetMoviewDetails(
  w   http.ResponseWriter,
  req *http.Request,
) {
  id := req.FormValue("id")
  log.Printf("GetMoviewDetails: id= %s", id)
  details, err := h.ctrl.Get(req.Context(), id)
  if err != nil {
    if errors.Is(err, movie.ErrNotFound){
      w.WriteHeader(http.StatusNotFound)
      return
    }
    log.Printf("Repository GET Error: %v\n", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  if err := common.WriteJSON(
    w,
    http.StatusOK,
    details,
  ); err != nil {
    log.Printf("Resposne encoding error: %v\n", err)
  }
}
