package server

import "context"

type Currency struct {
	log hclog.Logger
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateResponse) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase())
	return &protos.RateResponse()
}
