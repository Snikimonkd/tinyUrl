package delivery

import (
	"errors"
	"io/ioutil"
	"testing"

	"context"

	server_proto "github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/delivery/server"
	mocks "github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/repository/mocks"
	"github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/usecase"
	"github.com/Snikimonkd/tinyUrl/internal/tinyUrl/utils"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	utils.MainLogger = &utils.Logger{Logger: logrus.NewEntry(logrus.StandardLogger())}
	utils.MainLogger.Logger.Logger.SetOutput(ioutil.Discard)
}

var tinyStr string = "0123456789"

func GenerateMock() string {
	return tinyStr
}

func TestCreate_NotValidInput(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := usecase.TinyUrlUseCase{Repository: dbMock}

	handler := TinyUrlHandler{
		Usecase: &usecase,
	}

	fullUrl := server_proto.FullUrl{Val: "not url string"}
	ctx := context.Background()

	res, err := handler.Create(ctx, &fullUrl)

	var expected *server_proto.TinyUrl

	assert.Equal(t, status.Error(codes.InvalidArgument, "URL is not valid"), err)
	assert.Equal(t, expected, res)
}

func TestCreate_CheckIfFullUrlExist_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := usecase.TinyUrlUseCase{Repository: dbMock}

	handler := TinyUrlHandler{
		Usecase: &usecase,
	}

	fullUrl := server_proto.FullUrl{Val: "http://google.com/asdf"}
	ctx := context.Background()

	dbMock.EXPECT().CheckIfFullUrlExist(fullUrl.Val).Return("", errors.New("Some error"))

	res, err := handler.Create(ctx, &fullUrl)

	var expected *server_proto.TinyUrl

	assert.Equal(t, status.Error(codes.Internal, "Server error:Some error"), err)
	assert.Equal(t, expected, res)
}

func TestCreate_CheckIfFullUrlExist(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := usecase.TinyUrlUseCase{Repository: dbMock}

	handler := TinyUrlHandler{
		Usecase: &usecase,
	}

	fullUrl := server_proto.FullUrl{Val: "http://google.com/asdf"}
	ctx := context.Background()

	dbMock.EXPECT().CheckIfFullUrlExist(fullUrl.Val).Return(tinyStr, nil)

	res, err := handler.Create(ctx, &fullUrl)

	expected := server_proto.TinyUrl{Val: "http://my.domain.com/" + tinyStr}

	assert.Equal(t, nil, err)
	assert.Equal(t, expected.Val, res.Val)
}

func TestCreate_CheckIfTinyUrlExist_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := usecase.TinyUrlUseCase{Repository: dbMock}
	usecase.Gen = GenerateMock

	handler := TinyUrlHandler{
		Usecase: &usecase,
	}

	fullUrl := server_proto.FullUrl{Val: "http://google.com/asdf"}
	ctx := context.Background()

	dbMock.EXPECT().CheckIfFullUrlExist(fullUrl.Val).Return("", nil)
	dbMock.EXPECT().CheckIfTinyUrlExist(tinyStr).Return(false, errors.New("Some error"))

	var expected server_proto.TinyUrl

	res, err := handler.Create(ctx, &fullUrl)

	assert.Equal(t, nil, err)
	assert.Equal(t, &expected, res)
}

func TestCreate_RepositoryCreate_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := usecase.TinyUrlUseCase{Repository: dbMock}
	usecase.Gen = GenerateMock

	handler := TinyUrlHandler{
		Usecase: &usecase,
	}

	fullUrl := server_proto.FullUrl{Val: "http://google.com/asdf"}
	ctx := context.Background()

	dbMock.EXPECT().CheckIfFullUrlExist(fullUrl.Val).Return("", nil)
	dbMock.EXPECT().CheckIfTinyUrlExist(tinyStr).Return(false, nil)
	dbMock.EXPECT().Create(fullUrl.Val, tinyStr).Return(errors.New("Some error"))

	var expected *server_proto.TinyUrl

	res, err := handler.Create(ctx, &fullUrl)

	assert.Equal(t, status.Error(codes.Internal, "Server error:Some error"), err)
	assert.Equal(t, expected, res)
}

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := usecase.TinyUrlUseCase{Repository: dbMock}
	usecase.Gen = GenerateMock

	handler := TinyUrlHandler{
		Usecase: &usecase,
	}

	fullUrl := server_proto.FullUrl{Val: "http://google.com/asdf"}
	ctx := context.Background()

	dbMock.EXPECT().CheckIfFullUrlExist(fullUrl.Val).Return("", nil)
	dbMock.EXPECT().CheckIfTinyUrlExist(tinyStr).Return(false, nil)
	dbMock.EXPECT().Create(fullUrl.Val, tinyStr).Return(nil)

	res, err := handler.Create(ctx, &fullUrl)

	assert.Equal(t, nil, err)
	assert.Equal(t, "http://my.domain.com/"+tinyStr, res.Val)
}

func TestGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := usecase.TinyUrlUseCase{Repository: dbMock}

	handler := TinyUrlHandler{
		Usecase: &usecase,
	}

	tinyUrl := server_proto.TinyUrl{Val: "http://my.domain.com/" + tinyStr}
	ctx := context.Background()

	dbMock.EXPECT().Get(tinyStr).Return("http://google.com/asdf", nil)

	res, err := handler.Get(ctx, &tinyUrl)

	assert.Equal(t, nil, err)
	assert.Equal(t, "http://google.com/asdf", res.Val)
}

func TestGet_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := usecase.TinyUrlUseCase{Repository: dbMock}

	handler := TinyUrlHandler{
		Usecase: &usecase,
	}

	tinyUrl := server_proto.TinyUrl{Val: "http://my.domain.com/" + tinyStr}
	ctx := context.Background()

	dbMock.EXPECT().Get(tinyStr).Return("", errors.New("Some error"))

	var expected *server_proto.FullUrl

	res, err := handler.Get(ctx, &tinyUrl)

	assert.Equal(t, status.Error(codes.Internal, "Server error:Some error"), err)
	assert.Equal(t, expected, res)
}

func TestGet_UrlNotExist(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := usecase.TinyUrlUseCase{Repository: dbMock}

	handler := TinyUrlHandler{
		Usecase: &usecase,
	}

	tinyUrl := server_proto.TinyUrl{Val: "http://my.domain.com/" + tinyStr}
	ctx := context.Background()

	dbMock.EXPECT().Get(tinyStr).Return("", nil)

	var expected *server_proto.FullUrl

	res, err := handler.Get(ctx, &tinyUrl)

	assert.Equal(t, status.Error(codes.NotFound, "Can`t find URL:"+tinyUrl.Val), err)
	assert.Equal(t, expected, res)
}

func TestGet_NotValidInput(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := usecase.TinyUrlUseCase{Repository: dbMock}

	handler := TinyUrlHandler{
		Usecase: &usecase,
	}

	tinyUrl := server_proto.TinyUrl{Val: "not url string"}
	ctx := context.Background()

	var expected *server_proto.FullUrl

	res, err := handler.Get(ctx, &tinyUrl)

	assert.Equal(t, status.Error(codes.InvalidArgument, "URL is not valid"), err)
	assert.Equal(t, expected, res)
}

func TestGenerate(t *testing.T) {
	str1 := usecase.Generate()
	str2 := usecase.Generate()

	assert.NotEqual(t, str1, str2)
}
