package play

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const authUrl = "https://www.googleapis.com/oauth2/v4/token"

type InvalidPrivateKeyTypeError struct {
	privateKey interface{}
}

func (err InvalidPrivateKeyTypeError) Error() string {
	return fmt.Sprintf("invalid private key type: %T", err.privateKey)
}

type ServiceAccount struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type AuthError struct {
	ErrorType        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (err AuthError) Error() string {
	return fmt.Sprintf("error '%s': %s", err.ErrorType, err.ErrorDescription)
}

type AuthClient struct {
	client *http.Client
}

type AuthClientOption func(c *AuthClient)

func WithAuthHttpClient(client *http.Client) AuthClientOption {
	return func(c *AuthClient) {
		c.client = client
	}
}

func NewAuthClient(opts ...AuthClientOption) *AuthClient {
	auth := &AuthClient{}

	for _, opt := range opts {
		opt(auth)
	}

	if auth.client == nil {
		auth.client = defaultHttpClient()
	}

	return auth
}

func (auth *AuthClient) Authenticate(ctx context.Context, account *ServiceAccount) (*AccessToken, error) {
	pkey, err := auth.rsaPrivateKey(account)
	if err != nil {
		return nil, err
	}

	signedJwtToken, err := auth.makeJwtToken(account).SignedString(pkey)
	if err != nil {
		return nil, err
	}

	req, err := auth.makeAuthRequest(signedJwtToken)

	req = req.WithContext(ctx)

	resp, err := auth.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, auth.decodeError(dec)
	}

	return auth.decodeAccessToken(dec)
}

func (auth *AuthClient) makeJwtToken(account *ServiceAccount) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   account.ClientEmail,
		"scope": "https://www.googleapis.com/auth/androidpublisher",
		"aud":   authUrl,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	})
}

func (auth *AuthClient) makeAuthRequestValues(signedJwtToken string) *url.Values {
	params := url.Values{}
	params.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	params.Set("assertion", signedJwtToken)

	return &params
}

func (auth *AuthClient) makeAuthRequest(signedJwtToken string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, authUrl,
		bytes.NewBufferString(auth.makeAuthRequestValues(signedJwtToken).Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (auth *AuthClient) rsaPrivateKey(account *ServiceAccount) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(account.PrivateKey))

	pkey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := pkey.(*rsa.PrivateKey)
	if !ok {
		return nil, InvalidPrivateKeyTypeError{privateKey: key}
	}

	return key, nil
}

func (auth *AuthClient) decodeError(decoder *json.Decoder) error {
	var authError AuthError

	err := decoder.Decode(&authError)
	if err != nil {
		return err
	}

	return authError
}

func (auth *AuthClient) decodeAccessToken(decoder *json.Decoder) (*AccessToken, error) {
	var accessToken AccessToken

	err := decoder.Decode(&accessToken)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}
