package handler

import "errors"

var (
  ErrUnknownEndpoint = errors.New("attempted to access an unknown endpoint")
)
