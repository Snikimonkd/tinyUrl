package usecase

import (
	"math/rand"
	"strings"

	"github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/repository"
)

type TinyUrlUseCase struct {
	Repository repository.TinyUrlRepositoryInterface
}

type TinyUrlUsecaseInterface interface {
	generate() string
	Create(fullUrl string) (string, error)
	Get(tinyUrl string) (string, error)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")

func (u *TinyUrlUseCase) generate() string {
	randStr := make([]rune, 10)
	for i := 0; i < len(randStr); i++ {
		randStr[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(randStr)
}

func (u *TinyUrlUseCase) Create(fullUrl string) (string, error) {
	// Если для этого URL уже создан укороченный, то возвращаем его
	tinyUrl, err := u.Repository.CheckIfFullUrlExist(fullUrl)
	if err != nil {
		return "", err
	}
	if tinyUrl != "" {
		return ("http://my.domain.com/" + tinyUrl), nil
	}

	// Если сгенерировали уже существующий укороченный URL, то перегенерируем его
	var tinyUrlStr string
	for exist := true; exist == true; {
		tinyUrlStr = u.generate()
		exist, err = u.Repository.CheckIfTinyUrlExist(tinyUrlStr)
		if err != nil {
			return "", nil
		}
	}

	err = u.Repository.Create(fullUrl, tinyUrlStr)
	if err != nil {
		return "", err
	}

	return ("http://my.domain.com/" + tinyUrlStr), nil
}

func (u *TinyUrlUseCase) Get(tinyUrl string) (string, error) {
	trimedUrl := strings.TrimLeft(tinyUrl, "http://my.domain.com/")

	fullUrl, err := u.Repository.Get(trimedUrl)
	if err != nil {
		return "", err
	}

	return fullUrl, nil
}