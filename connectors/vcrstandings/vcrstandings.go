// Package vcrstandings send data
package vcrstandings

import (
	"context"

	"github.com/ianhaycox/vcrlive/connectors/api"
	"github.com/ianhaycox/vcrlive/model"
)

type VcrStandingsService struct {
	client api.APIClientInterface
	auth   api.Authenticator
}

func NewVcrStandingsService(client api.APIClientInterface, auth api.Authenticator) *VcrStandingsService {
	return &VcrStandingsService{
		client: client,
		auth:   auth,
	}
}

//go:generate mockgen -package vcrstandings -destination vcrstandings_mock.go -source vcrstandings.go
type VcrStandingsAPI interface {
	Post(ctx context.Context, livePositions *model.LivePositions) error
}
