package play

import (
	"net/http"
	"time"
)

func defaultHttpClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}
