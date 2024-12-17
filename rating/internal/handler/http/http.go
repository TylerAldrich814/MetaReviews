package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TylerAldrich/MetaMovies/common"
	"github.com/TylerAldrich/MetaMovies/common/endpoint"
	"github.com/TylerAldrich/MetaMovies/rating/internal/controller/rating"
	"github.com/TylerAldrich/MetaMovies/rating/internal/handler"
	"github.com/TylerAldrich/MetaMovies/rating/pkg/model"
)

// Handler defines a rating service controller
type Handler struct {
  ctrl *rating.Controller
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
    ep,
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

// New creates a new rating service HTTP Handler
func New(ctrl *rating.Controller) *Handler {
  return &Handler{ ctrl }
}

// HandleReq handles PUT and GET endpoint requests.
func(h *Handler) HandleReq(
  w   http.ResponseWriter, 
  req *http.Request,
) {
  recordID := model.RecordID(req.FormValue("id"))
  if recordID == "" {
    if err := common.WriteError(
      w,
      http.StatusBadRequest,
      "'id' is a required form value",
    ); err != nil {
      log.Printf("Response encode error: %v\n", err)
    }
    return
  }
  recordType := model.RecordType(req.FormValue("type"))
  if recordType == "" {
    if err := common.WriteError(
      w,
      http.StatusBadRequest,
      "'type' is a required form value",
    ); err != nil {
      log.Printf("Response encode error: %v\n", err)
    }
    return
  }
  switch req.Method {
  case http.MethodGet:
    v, err := h.ctrl.GetAggregatedRating(
      req.Context(),
      recordID,
      recordType,
    )
    if err != nil {
      if err := common.WriteError(
        w,
        http.StatusNotFound,
        err.Error(),
      ); err != nil {
        log.Printf("Response encode error: %v\n", err)
      }
    }
    if err := common.WriteJSON(
      w,
      http.StatusOK,
      v,
    ); err != nil {
      log.Printf("Response encode error: %v\n", err)
    }
  case http.MethodPut:
    userID := model.UserID(req.FormValue("userId"))
    if userID == "" {
      common.WriteError(
        w,
        http.StatusBadRequest,
        "'userId' is a required form value",
      )
    }

    v, err := strconv.ParseFloat(req.FormValue("value"), 64)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }
    if err := h.ctrl.PutRating(
      req.Context(), 
      recordID, 
      recordType, 
      &model.Rating{
        UserID: userID,
        Value: model.RatingValue(v),
      },
    ); err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }
  default:
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}
