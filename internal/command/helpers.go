package command

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/yurykabanov/google-play-edit/internal/pretty"
	"github.com/yurykabanov/google-play-edit/pkg/play"
)

func mustMakeHttpClient() *http.Client {
	proxyUrl := viper.GetString("proxy")
	insecure := viper.GetBool("proxy-insecure")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if proxyUrl != "" {
		proxy, err := url.Parse(proxyUrl)
		if err != nil {
			pretty.Errorf("Unable to parse proxy URL: %s", err.Error())
			os.Exit(1)
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}

		if insecure {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		client.Transport = transport
	}

	return client
}

func loadServiceAccount(path string) (*play.ServiceAccount, error) {
	var acc play.ServiceAccount

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)
	err = dec.Decode(&acc)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

func mustAuthenticate(client *http.Client) *play.AccessToken {
	var token *play.AccessToken
	accountPath := viper.GetString("account")
	accessToken := viper.GetString("token")

	if accountPath != "" {
		serviceAccount, err := loadServiceAccount(accountPath)
		if err != nil {
			pretty.Errorf("Unable to load service account: %s", err.Error())
			os.Exit(1)
		}

		token, err = play.NewAuthClient(play.WithAuthHttpClient(client)).Authenticate(context.Background(), serviceAccount)
		if err != nil {
			pretty.Errorf("Unable to authenticate: %s", err.Error())
			os.Exit(1)
		}
	} else if accessToken != "" {
		token = &play.AccessToken{AccessToken: accessToken, TokenType: "Bearer", ExpiresIn: 3600}
	} else {
		pretty.Errorf("Neither Service Account nor Access Token was specified")
		os.Exit(1)
	}

	if viper.GetBool("print-token") {
		pretty.Errorf("Access Token: %s\nExpires in: %d\n", token.AccessToken, token.ExpiresIn)
	}

	return token
}
