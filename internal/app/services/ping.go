package services

import "context"

type PingService interface {
	Pinging(ctx context.Context) (string, error)
}

type PingSvc struct{}

func (PingSvc) Pinging(ctx context.Context) (string, error) {

	return "Ping - HashApp", nil
}
