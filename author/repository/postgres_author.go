package repository

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/naveenpatilm/go-clean-arch/author"
	"github.com/naveenpatilm/go-clean-arch/models"
)

type mysqlAuthorRepo struct {
	DB *gorm.DB
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewMysqlAuthorRepository(db *gorm.DB) author.Repository {

	return &mysqlAuthorRepo{
		DB: db,
	}
}

func (m *mysqlAuthorRepo) GetByID(ctx context.Context, id int64) (*models.Author, error) {
	var author *models.Author
	err := m.DB.First(&author, id).Error
	if err != nil {
		return nil, err
	}
	if author != nil {
		return author, nil
	} else {
		return nil, models.ErrNotFound
	}
}
