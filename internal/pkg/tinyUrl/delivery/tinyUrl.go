package delivery

import (
	"context"

	server "github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/delivery/server"
	"github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/usecase"
)

type TinyUrlHandler struct {
	Usecase usecase.TinyUrlUsecaseInterface
}

type TinyUrlHandlerInterface interface {
	Create(ctx context.Context, url *server.FullUrl) (*server.TinyUrl, error)
	Get(ctx context.Context, tinyUrl *server.TinyUrl) (*server.FullUrl, error)
}

func (h *TinyUrlHandler) Create(ctx context.Context, url *server.FullUrl) (*server.TinyUrl, error) {
	return nil, nil
}

func (h *TinyUrlHandler) Get(ctx context.Context, tinyUrl *server.TinyUrl) (*server.FullUrl, error) {
	return nil, nil
}
