package client

import (
	"context"

	"github.com/SlothNinja/sn/v3"
)

// Client provide client structure of the Le Plateau service
type Client struct {
	*sn.GameClient[game, *game]
}

// New returns a new Client for the plateau service
func New(ctx context.Context, opts ...sn.Option) *Client {
	return (&Client{sn.NewGameClient[game, *game](ctx, opts...)}).addRoutes()
	// return cl.initUserDatastore(ctx).addRoutes()
	// sn.Debugf("Entering")
	// defer sn.Debugf("Exiting")

	// nClient := &Client{sn.NewGameClient[*game, *player](ctx, // sn.Options{
	// 	sn.WithProjectID(projectID()),
	// 	sn.WithUserProjectID(getUserProjectID()),
	// 	sn.WithUserDSURL(getUserDSURL()),
	// 	sn.WithLoggerID("plateau"),
	// 	sn.WithCORSAllow(
	// 		"https://plateau.fake-slothninja.com:8092/*",
	// 		"https://plateau.fake-slothninja.com:8091/sn/user/current",
	// 		"https://plateau.slothninja.com/*",
	// 	),
	// 	sn.WithPrefix("sn"),
	// )}
	// return nClient.addRoutes("sn")
}

// Close closes plateau service Client
func (cl *Client) Close() error {
	return cl.Client.Close()
}
