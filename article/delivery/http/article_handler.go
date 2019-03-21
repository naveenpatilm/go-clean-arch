package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/naveenpatilm/go-clean-arch/models"

	"github.com/naveenpatilm/go-clean-arch/article"

	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// HttpArticleHandler  represent the httphandler for article
type HttpArticleHandler struct {
	AUsecase article.Usecase
}

func NewArticleHttpHandler(r *mux.Router, us article.Usecase) {
	handler := &HttpArticleHandler{
		AUsecase: us,
	}
	r.HandleFunc("/articles", handler.FetchArticle).Methods("GET")
	r.HandleFunc("/articles", handler.Store).Methods("POST")
	r.HandleFunc("/article/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/article/{id}", handler.Delete).Methods("DELETE")

}

func (a *HttpArticleHandler) FetchArticle(w http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()
	fmt.Println(params)
	num, err := strconv.Atoi(params.Get("num"))
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cursor := params.Get("cursor")
	ctx := req.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listAr, err := a.AUsecase.Fetch(ctx, cursor, int64(num))

	if err != nil {
		w.WriteHeader(getStatusCode(err))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(getStatusCode(err))
	json.NewEncoder(w).Encode(listAr)
}

func (a *HttpArticleHandler) GetByID(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)
	fmt.Print("params-------------------------->>>>>>")
	fmt.Println(params)
	fmt.Println(req.URL.Query())
	idP, err := strconv.Atoi(params["id"])
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := int64(idP)

	ctx := req.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := a.AUsecase.GetByID(ctx, id)

	if err != nil {
		w.WriteHeader(getStatusCode(err))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(getStatusCode(err))
	json.NewEncoder(w).Encode(art)
}

func isRequestValid(m *models.Article) (bool, error) {

	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *HttpArticleHandler) Store(w http.ResponseWriter, req *http.Request) {
	var article models.Article
	err := json.NewDecoder(req.Body).Decode(&article)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if ok, err := isRequestValid(&article); !ok {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.AUsecase.Store(ctx, &article)

	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *HttpArticleHandler) Delete(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idP, err := strconv.Atoi(params["id"])
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := int64(idP)
	ctx := req.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.AUsecase.Delete(ctx, id)

	if err != nil {
		w.WriteHeader(getStatusCode(err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	case models.ErrUnprocessableEntity:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
