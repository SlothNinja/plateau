package client

import (
	"context"

	"github.com/SlothNinja/sn/v3"
)

const (
	msgEnter = "Entering"
	msgExit  = "Exiting"
)

// Client provide client structure of the Le Plateau service
type Client struct {
	*sn.GameClient[game, *game]
}

// New returns a new Client for the plateau service
func New(ctx context.Context, opts ...sn.Option) *Client {
	return (&Client{sn.NewGameClient[game, *game](ctx, opts...)}).addRoutes()
}

// Close closes plateau service Client
func (cl *Client) Close() error {
	return cl.Client.Close()
}
