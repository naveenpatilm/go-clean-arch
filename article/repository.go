package article

import (
	"context"

	"github.com/naveenpatilm/go-clean-arch/models"
)

// Repository represent the article's repository contract
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Article, err error)
	GetByID(ctx context.Context, id int64) (*models.Article, error)
	GetByTitle(ctx context.Context, title string) (*models.Article, error)
	Update(ctx context.Context, ar *models.Article) error
	Store(ctx context.Context, a *models.Article) error
	Delete(ctx context.Context, id int64) error
}
