package memory

import (
	"context"
	"sync"

	"github.com/TylerAldrich814/MetaReviews/metadata/internal/repository"
	"github.com/TylerAldrich814/MetaReviews/metadata/pkg/model"
)

// Repository defined a memory movie metadata repository
type Repository struct {
  sync.RWMutex
  data map[string]*model.Metadata
}

func New() *Repository {
  return &Repository{data: map[string]*model.Metadata{}}
}

// Get retrieves a new Memory Repository
func(r *Repository) Get(_ context.Context, id string)( *model.Metadata,error ){
  r.RLock()
  defer r.RUnlock()

  m, ok := r.data[id]
  if !ok {
    return nil, repository.ErrNotFound
  }

  return m,nil
}

// Put adds move metadata for a given move id.
func(r *Repository) Put(
  _        context.Context, 
  id       string, 
  metadata *model.Metadata,
) error {
  r.Lock()
  defer r.Unlock()
  r.data[id] = metadata
  return nil
}

