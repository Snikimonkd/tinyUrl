package repository

import (
	"github.com/Snikimonkd/tinyUrl/internal/pkg/models"
	"github.com/jmoiron/sqlx"
)

type TinyUrlRepository struct {
	DB *sqlx.DB
}

type TinyUrlRepositoryInterface interface {
	Create(url models.URL) (tinyUrl models.URL, err error)
	Get(tinyUrl models.URL) (url models.URL, err error)
}
