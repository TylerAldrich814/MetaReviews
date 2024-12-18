package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/TylerAldrich814/MetaMovies/common"
	"github.com/TylerAldrich814/MetaMovies/common/endpoint"
	"github.com/TylerAldrich814/MetaMovies/metadata/internal/controller/metadata"
	"github.com/TylerAldrich814/MetaMovies/metadata/internal/handler"
	"github.com/TylerAldrich814/MetaMovies/metadata/internal/repository"
)

// Handler defines a movie metadata HTTP Handler.
type Handler struct {
  ctrl *metadata.Controller
}

func(h *Handler) Handle(
  endpoint    endpoint.Endpoint,
  handlerFunc func(http.ResponseWriter, *http.Request),
) error {
  ep := endpoint.String()
  if ep == "unknown" {
    return handler.ErrUnknownEndpoint
  }

  http.Handle(
    "/"+ep,
    http.HandlerFunc(handlerFunc),
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

// New createa a new movie metadata HTTP Handler.
func New(ctrl *metadata.Controller) *Handler {
  return &Handler{ ctrl }
}

// GetMetadata handles GET /metadata requests.
func(h *Handler) GetMetadata(
  w   http.ResponseWriter, 
  req *http.Request) {
  id := req.FormValue("id")
  if id == "" {
    if err := common.WriteError(
      w, 
      http.StatusNotFound,
      "missing valid id form value",
    ); err != nil {
      log.Printf("Response encode error: %v\n", err)
    }
    return
  }

  ctx := req.Context()
  m, err := h.ctrl.Get(ctx, id)
  if err != nil && errors.Is(err, repository.ErrNotFound){
    if err := common.WriteError(
      w,
      http.StatusNotFound,
      err.Error(),
    ); err != nil {
      log.Printf("Response encode error: %v\n", err)
    }
    return
  }
  if err != nil {
    log.Printf("Repository GET error: %v\n", err)
    if err := common.WriteError(
      w,
      http.StatusInternalServerError,
      "An internal server error occurred",
    ); err != nil {
      log.Printf("Response encode error: %v\n", err)
    }
    return
  }

  if err := common.WriteJSON(
    w,
    http.StatusOK,
    m,
  ); err != nil {
    log.Printf("Response encode error: %v\n", err)
  }
}
