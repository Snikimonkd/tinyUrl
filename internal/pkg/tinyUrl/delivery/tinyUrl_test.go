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

type UsecaseMock struct {
	usecase.TinyUrlUseCase
}

var tinyStr string = "0123456789"

func (u *UsecaseMock) generate() string {
	return tinyStr
}

func TestCreate_CheckIfFullUrlExist_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := UsecaseMock{TinyUrlUseCase: usecase.TinyUrlUseCase{Repository: dbMock}}

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
	usecase := UsecaseMock{TinyUrlUseCase: usecase.TinyUrlUseCase{Repository: dbMock}}

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

func TestCreate_NotValidInput(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
	usecase := UsecaseMock{TinyUrlUseCase: usecase.TinyUrlUseCase{Repository: dbMock}}

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

// func TestCreate(t *testing.T) {
// 	utils.MainLogger = &utils.Logger{Logger: logrus.NewEntry(logrus.StandardLogger())}
// 	mockCtrl := gomock.NewController(t)
// 	dbMock := mocks.NewMockTinyUrlRepositoryInterface(mockCtrl)
// 	usecase := UsecaseMock{TinyUrlUseCase: usecase.TinyUrlUseCase{Repository: dbMock}}

// 	handler := TinyUrlHandler{
// 		Usecase: &usecase,
// 	}

// 	fmt.Println(usecase.generate())

// 	fullUrl := server_proto.FullUrl{Val: "http://google.com/asdf"}
// 	ctx := context.Background()

// 	dbMock.EXPECT().CheckIfFullUrlExist(fullUrl.Val).Return("", nil)
// 	dbMock.EXPECT().CheckIfTinyUrlExist(tinyStr).Return(false, nil)
// 	dbMock.EXPECT().Create(fullUrl.Val, tinyStr).Return(nil)

// 	res, err := handler.Create(ctx, &fullUrl)

// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, tinyStr, res.Val)
// }
