package proxy

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/shamhub/exercise/pkg/errorlib"
)

func IsKindServiceAlive(URL *url.URL) (bool, *errorlib.HealthCheckError) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", URL.String(), nil)
	if err != nil {
		return false, errorlib.NewHealthCheckError(http.StatusTooManyRequests, err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, errorlib.NewHealthCheckError(resp.StatusCode, err.Error())
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}
