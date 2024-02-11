package main

import (
	"context"

	"github.com/SlothNinja/plateau/main/client"
	"github.com/SlothNinja/sn/v3"
)

func main() {
	cl := client.New(
		context.Background(),
		sn.WithLoggerID("plateau-service"),
	)

	defer func() {
		if err := cl.Close(); err != nil {
			sn.Warningf("error when closing client: %w", err)
		}
	}()

	if sn.IsProduction() {
		sn.Debugf("in production")
		cl.Router.Run()
		return
	}

	sn.Debugf("in debug")
	cl.Router.RunTLS(":"+cl.GetPort(), "cert.pem", "key.pem")
}
