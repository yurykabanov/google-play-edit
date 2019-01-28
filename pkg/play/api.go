package play

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiBaseUrl = "https://www.googleapis.com/androidpublisher/v3/applications/%s"

type ApiError struct {
	ErrorDefinition struct {
		Code    int                      `json:"code"`
		Errors  []map[string]interface{} `json:"errors"`
		Message string                   `json:"message"`
	} `json:"error"`
}

func decodeApiErrorResponse(decoder *json.Decoder) error {
	var apiError ApiError

	err := decoder.Decode(&apiError)
	if err != nil {
		return err
	}

	return apiError
}

func (err ApiError) Error() string {
	return fmt.Sprintf("api error %d: %s", err.ErrorDefinition.Code, err.ErrorDefinition.Message)
}

type Api struct {
	client *http.Client

	Edits    EditsApi
	Listings EditListingsApi
	Images   EditImagesApi
}

type ApiClientOption func(c *Api)

func WithApiHttpClient(client *http.Client) ApiClientOption {
	return func(c *Api) {
		c.client = client
	}
}

func NewApi(opts ...ApiClientOption) *Api {
	api := &Api{}

	for _, opt :=range opts {
		opt(api)
	}

	if api.client == nil {
		api.client = defaultHttpClient()
	}

	return &Api{
		Edits: &editsApi{
			client: api.client,
		},
		Listings: &editListingsApi{
			client: api.client,
		},
		Images: &editImagesApi{
			client: api.client,
		},
	}
}

type EditsApi interface {
	Commit(ctx context.Context, token *AccessToken, packageName string, editId string) (*Edit, error)
	Delete(ctx context.Context, token *AccessToken, packageName string, editId string) error
	Get(ctx context.Context, token *AccessToken, packageName string, editId string) (*Edit, error)
	Insert(ctx context.Context, token *AccessToken, packageName string) (*Edit, error)
	Validate(ctx context.Context, token *AccessToken, packageName string, editId string) (*Edit, error)
}

type EditListingsApi interface {
	Delete(ctx context.Context, token *AccessToken, packageName string, editId string, lang string) error
	DeleteAll(ctx context.Context, token *AccessToken, packageName string, editId string) error
	Get(ctx context.Context, token *AccessToken, packageName string, editId string, lang string) (*Listing, error)
	List(ctx context.Context, token *AccessToken, packageName string, editId string) ([]Listing, error)
	Update(ctx context.Context, token *AccessToken, packageName string, editId string, listing *Listing) (*Listing, error)
}

type EditImagesApi interface {
	Delete(ctx context.Context, token *AccessToken, packageName string, editId string, lang string, imageType EditImageType, imageId string) error
	DeleteAll(ctx context.Context, token *AccessToken, packageName string, editId string, lang string, imageType EditImageType) ([]Image, error)
	List(ctx context.Context, token *AccessToken, packageName string, editId string, lang string, imageType EditImageType) ([]Image, error)
	Upload(ctx context.Context, token *AccessToken, packageName string, editId string, lang string, imageType EditImageType, imageReader io.ReadSeeker) (*Image, error)
}
