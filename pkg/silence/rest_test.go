package silence

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dacamposol/monitoring-operator/pkg/silence/mocks"
)

func TestPost(t *testing.T) {
	client := new(mocks.HttpClient)

	response := `{"hello": "world"}`
	body := bytes.NewReader([]byte(response))

	client.EXPECT().Do(mock.Anything).Return(&http.Response{Body: io.NopCloser(body)}, nil)

	result, err := post(context.Background(), client, "https://localhost", map[string]string{"hello": "world"})

	assert.Equal(t, map[string]string{"hello": "world"}, result)
	assert.Nil(t, err)
}
