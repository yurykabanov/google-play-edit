package task

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"

	"github.com/yurykabanov/google-play-edit/pkg/play"
)

type mockEditApi struct {
	mock.Mock
}

func (mock *mockEditApi) Commit(ctx context.Context, token *play.AccessToken, packageName string, editId string) (*play.Edit, error) {
	args := mock.Called(ctx, token, packageName, editId)

	return args.Get(0).(*play.Edit), args.Error(1)
}

func (mock *mockEditApi) Delete(ctx context.Context, token *play.AccessToken, packageName string, editId string) error {
	args := mock.Called(ctx, token, packageName, editId)

	return args.Error(0)
}

func (mock *mockEditApi) Get(ctx context.Context, token *play.AccessToken, packageName string, editId string) (*play.Edit, error) {
	args := mock.Called(ctx, token, packageName, editId)

	return args.Get(0).(*play.Edit), args.Error(1)
}

func (mock *mockEditApi) Insert(ctx context.Context, token *play.AccessToken, packageName string) (*play.Edit, error) {
	args := mock.Called(ctx, token, packageName)

	return args.Get(0).(*play.Edit), args.Error(1)
}

func (mock *mockEditApi) Validate(ctx context.Context, token *play.AccessToken, packageName string, editId string) (*play.Edit, error) {
	args := mock.Called(ctx, token, packageName, editId)

	return args.Get(0).(*play.Edit), args.Error(1)
}

type mockListingApi struct {
	mock.Mock
}

func (mock *mockListingApi) Delete(ctx context.Context, token *play.AccessToken, packageName string, editId string, lang string) error {
	args := mock.Called(ctx, token, packageName, editId, lang)

	return args.Error(0)
}

func (mock *mockListingApi) DeleteAll(ctx context.Context, token *play.AccessToken, packageName string, editId string) error {
	args := mock.Called(ctx, token, packageName, editId)

	return args.Error(0)
}

func (mock *mockListingApi) Get(ctx context.Context, token *play.AccessToken, packageName string, editId string, lang string) (*play.Listing, error) {
	args := mock.Called(ctx, token, packageName, editId, lang)

	return args.Get(0).(*play.Listing), args.Error(1)
}

func (mock *mockListingApi) List(ctx context.Context, token *play.AccessToken, packageName string, editId string) ([]play.Listing, error) {
	args := mock.Called(ctx, token, packageName, editId)

	return args.Get(0).([]play.Listing), args.Error(1)
}

func (mock *mockListingApi) Update(ctx context.Context, token *play.AccessToken, packageName string, editId string, listing *play.Listing) (*play.Listing, error) {
	args := mock.Called(ctx, token, packageName, editId, listing)

	return args.Get(0).(*play.Listing), args.Error(1)
}

type mockImagesApi struct {
	mock.Mock
}

func (mock *mockImagesApi) Delete(ctx context.Context, token *play.AccessToken, packageName string, editId string, lang string, imageType play.EditImageType, imageId string) error {
	args := mock.Called(ctx, token, packageName, editId, lang, imageType, imageId)

	return args.Error(0)
}

func (mock *mockImagesApi) DeleteAll(ctx context.Context, token *play.AccessToken, packageName string, editId string, lang string, imageType play.EditImageType) ([]play.Image, error) {
	args := mock.Called(ctx, token, packageName, editId, lang, imageType)

	return args.Get(0).([]play.Image), args.Error(1)
}

func (mock *mockImagesApi) List(ctx context.Context, token *play.AccessToken, packageName string, editId string, lang string, imageType play.EditImageType) ([]play.Image, error) {
	args := mock.Called(ctx, token, packageName, editId, lang, imageType)

	return args.Get(0).([]play.Image), args.Error(1)
}

func (mock *mockImagesApi) Upload(ctx context.Context, token *play.AccessToken, packageName string, editId string, lang string, imageType play.EditImageType, imageReader io.ReadSeeker) (*play.Image, error) {
	args := mock.Called(ctx, token, packageName, editId, lang, imageType, imageReader)

	return args.Get(0).(*play.Image), args.Error(1)
}


