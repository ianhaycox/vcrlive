package vcrstandings

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/ianhaycox/vcrlive/connectors/api"
	"github.com/ianhaycox/vcrlive/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPost(t *testing.T) {
	t.Run("Happy path should not return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		client := api.NewMockAPIClientInterface(ctrl)
		client.EXPECT().PrepareRequest(ctx, "", "POST", nil, &model.LivePositions{})
		client.EXPECT().CallAPI(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("OK"))}, nil)

		v := NewVcrStandingsService(client, nil)

		err := v.Post(ctx, &model.LivePositions{})
		assert.NoError(t, err)
	})
}
