package delivery

import (
	"context"
	"database/sql"
	"strings"

	server "github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/delivery/server"
	"github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/usecase"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TinyUrlHandler struct {
	Usecase usecase.TinyUrlUsecaseInterface
}

type TinyUrlHandlerInterface interface {
	Create(ctx context.Context, fullUrl *server.FullUrl) (*server.TinyUrl, error)
	Get(ctx context.Context, tinyUrl *server.TinyUrl) (*server.FullUrl, error)
}

func (h *TinyUrlHandler) Create(ctx context.Context, fullUrl *server.FullUrl) (*server.TinyUrl, error) {
	ok := govalidator.IsURL(fullUrl.Val)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "URL is not valid")
	}

	res, err := h.Usecase.Create(fullUrl.Val)
	if err != nil {
		return nil, status.Error(codes.Internal, "Server error:"+err.Error())
	}

	tinyUrl := server.TinyUrl{Val: res}
	return &tinyUrl, nil
}

func (h *TinyUrlHandler) Get(ctx context.Context, tinyUrl *server.TinyUrl) (*server.FullUrl, error) {
	ok := govalidator.IsURL(tinyUrl.Val)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "URL is not valid")
	}

	trimedUrl := strings.TrimLeft(tinyUrl.Val, "http://")

	res, err := h.Usecase.Get(trimedUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Can`t find full URL of this tiny URL")
		}
		return nil, status.Error(codes.Internal, "Server error:"+err.Error())
	}

	return &server.FullUrl{Val: res}, nil
}
