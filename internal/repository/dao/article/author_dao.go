package article

import (
	"context"
)

type AuthorDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
}

func NewAuthorDAO() AuthorDAO {}
