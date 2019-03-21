package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	articleHttp "github.com/naveenpatilm/go-clean-arch/article/delivery/http"
	"github.com/naveenpatilm/go-clean-arch/article/mocks"
	"github.com/naveenpatilm/go-clean-arch/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/bxcodec/faker"
)

func TestFetch(t *testing.T) {
	var mockArticle models.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	mockListArticle := make([]*models.Article, 0)
	mockListArticle = append(mockListArticle, &mockArticle)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListArticle, nil)

	req, err := http.NewRequest(http.MethodGet, "/articles?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	handler := articleHttp.HttpArticleHandler{
		AUsecase: mockUCase,
	}
	handler.FetchArticle(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.Usecase)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(nil, models.ErrInternalServerError)

	req, err := http.NewRequest(http.MethodGet, "/articles?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	handler := articleHttp.HttpArticleHandler{
		AUsecase: mockUCase,
	}
	handler.FetchArticle(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockArticle models.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockArticle.ID)

	mockUCase.On("GetByID", mock.Anything, int64(num)).Return(&mockArticle, nil)

	req, err := http.NewRequest(http.MethodGet, "/article/"+strconv.Itoa(int(num)), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	handler := articleHttp.HttpArticleHandler{
		AUsecase: mockUCase,
	}
	router := mux.NewRouter()
	mux.SetURLVars(req, map[string]string{
		"id": strconv.Itoa(num),
	})
	router.HandleFunc("/article/{id}", handler.GetByID)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockArticle := models.Article{
		Title:     "Title",
		Content:   "Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockArticle := mockArticle
	tempMockArticle.ID = 0
	mockUCase := new(mocks.Usecase)

	j, err := json.Marshal(tempMockArticle)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*models.Article")).Return(nil)

	req, err := http.NewRequest(http.MethodPost, "/articles", strings.NewReader(string(j)))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	handler := articleHttp.HttpArticleHandler{
		AUsecase: mockUCase,
	}
	handler.Store(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockArticle models.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockArticle.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	req, err := http.NewRequest(http.MethodDelete, "/article/"+strconv.Itoa(int(num)), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	router := mux.NewRouter()
	mux.SetURLVars(req, map[string]string{
		"id": strconv.Itoa(num),
	})
	handler := articleHttp.HttpArticleHandler{
		AUsecase: mockUCase,
	}
	router.HandleFunc("/article/{id}", handler.Delete)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)

}
