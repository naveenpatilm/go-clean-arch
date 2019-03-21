package repository

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/sirupsen/logrus"

	"github.com/naveenpatilm/go-clean-arch/article"
	"github.com/naveenpatilm/go-clean-arch/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mysqlArticleRepository struct {
	DB *gorm.DB
}

// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewMysqlArticleRepository(DB *gorm.DB) article.Repository {

	return &mysqlArticleRepository{DB}
}

func (m *mysqlArticleRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Article, error) {

	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, models.ErrBadParamInput
	}
	var articles []*models.Article
	err = m.DB.Where("created_at > ?", decodedCursor).Order("created_at", true).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	if len(articles) > 0 {
		return articles, nil
	} else {
		return nil, models.ErrNotFound
	}
}

func (m *mysqlArticleRepository) GetByID(ctx context.Context, id int64) (*models.Article, error) {
	var article *models.Article
	err := m.DB.First(&article, id).Error
	if err != nil {
		return nil, err
	}
	if article != nil {
		return article, nil
	} else {
		return nil, models.ErrNotFound
	}
}

func (m *mysqlArticleRepository) GetByTitle(ctx context.Context, title string) (*models.Article, error) {
	var article *models.Article
	err := m.DB.Where("title = ?", title).First(&article).Error
	if err != nil {
		return nil, err
	}

	if article != nil {
		return article, nil
	} else {
		return nil, models.ErrNotFound
	}
}

func (m *mysqlArticleRepository) Store(ctx context.Context, a *models.Article) error {
	err := m.DB.Create(&a).Error
	if err != nil {
		return err
	}
	logrus.Debug("Created At: ", a.CreatedAt)
	return nil
}

func (m *mysqlArticleRepository) Delete(ctx context.Context, id int64) error {
	res := m.DB.Where("id = ?", id).Delete(models.Article{})
	err := res.Error
	if err != nil {
		return err
	}
	rowsAffected := res.RowsAffected
	if rowsAffected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAffected)
		return err
	}

	return nil
}

func (m *mysqlArticleRepository) Update(ctx context.Context, ar *models.Article) error {
	res := m.DB.Save(&ar)

	err := res.Error
	if err != nil {
		return err
	}

	affected := res.RowsAffected
	if affected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affected)
		return err
	}

	return nil
}

func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(byt)
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}

func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
