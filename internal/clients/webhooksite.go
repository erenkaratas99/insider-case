package clients

import (
	"errors"
	"github.com/valyala/fasthttp"
)

type WebHookSiteClient struct {
	client *BaseClient
}

func NewWebHookSiteClient() *WebHookSiteClient {
	return &WebHookSiteClient{client: NewBaseClient("https://webhook.site/1ee3c9d3-4434-46bf-8236-6dc92994d88a")}
}

func (wc *WebHookSiteClient) SendMessage(body interface{}) error {
	opts := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := wc.client.POST("", body, opts)
	if err != nil {
		return err
	}

	if res == nil {
		return errors.New("couldn't send the data to the webhook site")
	}

	if res.StatusCode() != fasthttp.StatusOK {
		return errors.New("couldn't send the data to the webhook site")
	}

	fasthttp.ReleaseResponse(res)

	return nil
}
