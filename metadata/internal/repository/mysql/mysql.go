package mysql

import (
	"context"
	"database/sql"

	"github.com/TylerAldrich814/MetaReviews/metadata/internal/repository"
	"github.com/TylerAldrich814/MetaReviews/metadata/pkg/model"
	_ "github.com/go-sql-driver/mysql"
)

// Repository defines a MySql Movie Database
type Repository struct {
  db *sql.DB
}

// New creates a new MySQL Repository
func New()( *Repository,error ){
  db, err := sql.Open(
    "mysql",
    "root:root_passwrod@/MetaMovieDB",
  )
  if err != nil {
    return nil, err
  }

  return &Repository{ db }, nil
}

// Get retreives movie metadata by querying the movie ID
func(r *Repository) Get(
  ctx context.Context,
  id  string,
)( *model.Metadata, error){
  var title, description, director string
  row := r.db.QueryRowContext(
    ctx, 
    "SELECT title, description, director FROM movies WHERE id = ?",
    id,
  )
  if err := row.Scan(
    &title, 
    &description, 
    &director,
  ); err != nil {
    if err == sql.ErrNoRows {
      return nil, repository.ErrNotFound
    }
    return nil, err
  }

  return &model.Metadata{
    ID          : id,
    Title       : title,
    Description : description,
    Director    : director,
  },nil
}

// Put stores a movie metadata for a given movie ID.
func(r *Repository) Put(
  ctx      context.Context,
  id       string,
  metadata *model.Metadata,
) error {
  _, err := r.db.ExecContext(
    ctx,
    `
    INSERT INTO movies (id, title, description, director)
    VALUES (?, ?, ?, ?)
    `,
    id, 
    metadata.Title,
    metadata.Description,
    metadata.Director,
  )
  return err
}
