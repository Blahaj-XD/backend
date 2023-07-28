package platform

import (
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
)

func NewHeimdallHTTPClient() *httpclient.Client {
	backoffInterval := 2 * time.Millisecond
	maximumJitterInterval := 5 * time.Millisecond
	backoff := heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)
	retrier := heimdall.NewRetrier(backoff)

	timeout := 10000 * time.Millisecond
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(4),
	)

	return client
}
