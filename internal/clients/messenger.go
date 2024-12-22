package clients

import (
	"errors"
	"github.com/valyala/fasthttp"
	"insider/configs/appConfigs"
)

type MessengerClient struct {
	client *BaseClient
}

func NewMessengerClient(cfg *appConfigs.Configurations) *MessengerClient {
	return &MessengerClient{client: NewBaseClient(cfg.MessengerApi.Url)}
}

func (jc *MessengerClient) GetTwoPendingMessages(offset string) ([]byte, error) {
	opts := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := jc.client.GET("/messenger/get-two?offset="+offset, opts)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New("couldn't fetch the data from Messenger")
	}

	bodyCopy := append([]byte{}, res.Body()...)
	if res.StatusCode() != fasthttp.StatusOK {

		return nil, errors.New("couldn't fetch the data from Messenger")
	}

	fasthttp.ReleaseResponse(res)

	return bodyCopy, nil
}

func (jc *MessengerClient) CommitMessageAsSent(msgId string) error {
	opts := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := jc.client.PUT("/messenger/commit/"+msgId, nil, opts)
	if err != nil {
		return err
	}

	if res.StatusCode() != fasthttp.StatusOK {

		return errors.New("couldn't commit the message to Messenger")
	}

	fasthttp.ReleaseResponse(res)

	return nil
}
