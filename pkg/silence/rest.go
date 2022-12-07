package silence

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func post[T any](ctx context.Context, client HttpClient, url string, data T) (T, error) {
	var model T

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return model, err
	}

	byteReader := bytes.NewReader(dataBytes)

	req, err := http.NewRequestWithContext(ctx, "POST", url, byteReader)
	if err != nil {
		return model, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return model, err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return model, err
	}

	return parseJSON[T](body)
}

func remove(ctx context.Context, client HttpClient, url string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("received non-200 Status Code from [DELETE] Request")
	}

	return nil
}

func parseJSON[T any](src []byte) (T, error) {
	var model T
	if err := json.Unmarshal(src, &model); err != nil {
		return model, err
	}
	return model, nil
}
