package vcrstandings

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ianhaycox/vcrlive/connectors/api"
	"github.com/ianhaycox/vcrlive/model"
)

func (v *VcrStandingsService) Post(ctx context.Context, livePositions *model.LivePositions) error {
	if v.client == nil {
		fmt.Println(livePositions)
		return nil
	}

	r, err := v.client.PrepareRequest(ctx, "", http.MethodPost, nil, livePositions)
	if err != nil {
		return err
	}

	response, err := v.client.CallAPI(r) //nolint:bodyclose // ok
	if err != nil || response == nil {
		return err
	}

	defer api.BodyClose(response)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return v.client.ReportError(response, body)
	}

	return nil
}
