package play

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	editsApiBaseUrl  = apiBaseUrl + "/edits"
	editsApiCommit   = editsApiBaseUrl + "/%s:commit"
	editsApiDelete   = editsApiBaseUrl + "/%s"
	editsApiGet      = editsApiBaseUrl + "/%s"
	editsApiInsert   = editsApiBaseUrl
	editsApiValidate = editsApiBaseUrl + "/%s:validate"
)

type Edit struct {
	Id                string `json:"id"`
	ExpiryTimeSeconds string `json:"expiryTimeSeconds"`
}

type editsApi struct {
	client *http.Client
}

func decodeEditResponse(decoder *json.Decoder) (*Edit, error) {
	var edit Edit

	err := decoder.Decode(&edit)
	if err != nil {
		return nil, err
	}
	return &edit, nil
}

func (api *editsApi) Commit(ctx context.Context, token *AccessToken, packageName string, editId string) (*Edit, error) {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf(editsApiCommit, packageName, editId), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	req = req.WithContext(ctx)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, decodeApiErrorResponse(dec)
	}

	return decodeEditResponse(dec)
}

func (api *editsApi) Delete(ctx context.Context, token *AccessToken, packageName string, editId string) error {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(editsApiDelete, packageName, editId), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	req = req.WithContext(ctx)

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return decodeApiErrorResponse(dec)
	}

	return nil
}

func (api *editsApi) Get(ctx context.Context, token *AccessToken, packageName string, editId string) (*Edit, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(editsApiGet, packageName, editId), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	req = req.WithContext(ctx)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, decodeApiErrorResponse(dec)
	}

	return decodeEditResponse(dec)
}

func (api *editsApi) Insert(ctx context.Context, token *AccessToken, packageName string) (*Edit, error) {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf(editsApiInsert, packageName), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	req = req.WithContext(ctx)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, decodeApiErrorResponse(dec)
	}

	return decodeEditResponse(dec)
}

func (api *editsApi) Validate(ctx context.Context, token *AccessToken, packageName string, editId string) (*Edit, error) {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf(editsApiValidate, packageName, editId), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	req = req.WithContext(ctx)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, decodeApiErrorResponse(dec)
	}

	return decodeEditResponse(dec)
}
