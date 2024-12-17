package repository

import "errors"

var (
  // ErrNotFound is returned when a requested record is not found
  ErrNotFound = errors.New("not found")
)
