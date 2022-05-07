package repository

import (
	"github.com/jmoiron/sqlx"
)

type TinyUrlSQLRepository struct {
	DB *sqlx.DB
}

type TinyUrlInMemoryRepository struct {
	DB map[string]string
}

type TinyUrlRepositoryInterface interface {
	Create(fullUrl, tinyUrl string) error
	Get(tinyUrl string) (string, error)
	CheckIfFullUrlExist(fullUrl string) (string, error)
	CheckIfTinyUrlExist(tinyUrl string) (bool, error)
}

func (r *TinyUrlSQLRepository) Create(fullUrl, tinyUrl string) error {
	_, err := r.DB.Exec(`INSERT INTO urls (fullurl, tinyurl) VALUES ($1, $2)`, fullUrl, tinyUrl)

	return err
}

func (r *TinyUrlSQLRepository) Get(tinyUrl string) (string, error) {
	var fullUrl []string
	err := r.DB.Select(&fullUrl, `SELECT fullurl FROM urls WHERE tinyurl = $1`, tinyUrl)
	if err != nil {
		return "", err
	}
	if len(fullUrl) == 0 {
		return "", nil
	}

	return fullUrl[0], nil
}

func (r *TinyUrlSQLRepository) CheckIfTinyUrlExist(tinyUrl string) (bool, error) {
	var buf []string
	err := r.DB.Select(&buf, `SELECT tinyUrl FROM urls WHERE tinyurl = $1`, tinyUrl)
	if err != nil {
		return false, err
	}
	if len(buf) == 0 {
		return false, nil
	}

	return true, nil
}

func (r *TinyUrlSQLRepository) CheckIfFullUrlExist(fullUrl string) (string, error) {
	var buf []string
	err := r.DB.Select(&buf, `SELECT tinyurl FROM urls WHERE fullurl = $1`, fullUrl)
	if err != nil {
		return "", err
	}
	if len(buf) == 0 {
		return "", nil
	}

	return buf[0], nil
}

func (r *TinyUrlInMemoryRepository) Create(fullUrl, tinyUrl string) error {
	r.DB[tinyUrl] = fullUrl

	return nil
}

func (r *TinyUrlInMemoryRepository) Get(tinyUrl string) (string, error) {
	var url string
	var ok bool
	if url, ok = r.DB[tinyUrl]; !ok {
		return "", nil
	}

	return url, nil
}

func (r *TinyUrlInMemoryRepository) CheckIfTinyUrlExist(tinyUrl string) (bool, error) {
	if _, ok := r.DB[tinyUrl]; !ok {
		return false, nil
	}

	return true, nil
}

func (r *TinyUrlInMemoryRepository) CheckIfFullUrlExist(fullUrl string) (string, error) {
	for key, val := range r.DB {
		if val == fullUrl {
			return key, nil
		}
	}

	return "", nil
}
