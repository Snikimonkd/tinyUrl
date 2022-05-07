package delivery

import (
	"context"

	_ "github.com/jackc/pgx"

	server "github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/delivery/server"
	"github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/usecase"
	"github.com/Snikimonkd/tinyUrl/internal/tinyUrl/utils"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TinyUrlHandler struct {
	server.TinyUrlServerServer
	Usecase usecase.TinyUrlUsecaseInterface
}

type TinyUrlHandlerInterface interface {
	Create(ctx context.Context, fullUrl *server.FullUrl) (*server.TinyUrl, error)
	Get(ctx context.Context, tinyUrl *server.TinyUrl) (*server.FullUrl, error)
}

func (h *TinyUrlHandler) Create(ctx context.Context, fullUrl *server.FullUrl) (*server.TinyUrl, error) {
	ok := govalidator.IsURL(fullUrl.Val)
	if !ok {
		utils.MainLogger.LogInfo(status.Error(codes.InvalidArgument, "URL is not valid"))
		return nil, status.Error(codes.InvalidArgument, "URL is not valid")
	}

	res, err := h.Usecase.Create(fullUrl.Val)
	if err != nil {
		utils.MainLogger.LogError(status.Error(codes.Internal, "Server error:"+err.Error()))
		return nil, status.Error(codes.Internal, "Server error:"+err.Error())
	}

	utils.MainLogger.LogInfo(res)
	return &server.TinyUrl{Val: res}, nil
}

func (h *TinyUrlHandler) Get(ctx context.Context, tinyUrl *server.TinyUrl) (*server.FullUrl, error) {
	ok := govalidator.IsURL(tinyUrl.Val)
	if !ok {
		utils.MainLogger.Logger.Infoln(status.Error(codes.InvalidArgument, "URL is not valid"))
		return nil, status.Error(codes.InvalidArgument, "URL is not valid")
	}

	res, err := h.Usecase.Get(tinyUrl.Val)
	if err != nil {
		utils.MainLogger.LogError(status.Error(codes.Internal, "Server error:"+err.Error()))
		return nil, status.Error(codes.Internal, "Server error:"+err.Error())
	}
	if res == "" {
		utils.MainLogger.LogInfo(status.Error(codes.NotFound, "Can`t find URL:"+tinyUrl.Val))
		return nil, status.Error(codes.NotFound, "Can`t find URL:"+tinyUrl.Val)
	}

	utils.MainLogger.LogInfo(res)
	return &server.FullUrl{Val: res}, nil
}
