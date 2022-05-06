package usecase

import (
	"github.com/Snikimonkd/tinyUrl/internal/pkg/models"
	"github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/repository"
)

type TinyUrlUseCase struct {
	Usecase repository.TinyUrlRepositoryInterface
}

type TinyUrlUsecaseInterface interface {
	Generate() string
	Create(url models.URL) (tinyUrl models.URL, err error)
	Get(tinyUrl models.URL) (url models.URL, err error)
}
