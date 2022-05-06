package repository

import (
	"errors"

	"github.com/Snikimonkd/tinyUrl/internal/pkg/models"
	"github.com/jmoiron/sqlx"
)

type TinyUrlSQLRepository struct {
	DB *sqlx.DB
}

type TinyUrlInMemoryRepository struct {
	DB map[string]string
}

type TinyUrlRepositoryInterface interface {
	Create(tinyUrl, fullUrl string) error
	Get(tinyUrl string) (string, error)
}

func (r *TinyUrlSQLRepository) Create(fullUrl, tinyUrl string) error {
	_, err := r.DB.Exec(`INSERT INTO urls (fullurl, tinyurl) VALUES ($1, $2)`, tinyUrl, fullUrl)

	return err
}

func (r *TinyUrlSQLRepository) Get(tinyUrl string) (models.URL, error) {
	var fullUrl models.URL
	err := r.DB.Select(&fullUrl, `SELECT fullurl FROM urls WHERE tinyurl = $1`, tinyUrl)
	if err != nil {
		return models.URL{}, err
	}

	return fullUrl, nil
}

func (r *TinyUrlInMemoryRepository) Create(fullUrl, tinyUrl string) error {
	r.DB[tinyUrl] = fullUrl

	return nil
}

func (r *TinyUrlInMemoryRepository) Get(tinyUrl string) (models.URL, error) {
	var url models.URL
	var ok bool
	if url.Val, ok = r.DB[tinyUrl]; !ok {
		return models.URL{}, errors.New("Url does not exist")
	}

	return url, nil
}
