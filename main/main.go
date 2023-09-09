package main

import (
	"context"

	"github.com/SlothNinja/plateau/main/client"
	"github.com/SlothNinja/sn/v3"
)

func main() {
	cl := client.New(
		context.Background(),
		sn.WithLoggerID("user-service"),
	)

	defer func() {
		if err := cl.Close(); err != nil {
			sn.Warningf("error when closing client: %w", err)
		}
	}()

	if sn.IsProduction() {
		cl.Router.Run()
		return
	}

	cl.Router.RunTLS(":"+cl.GetPort(), "cert.pem", "key.pem")
}
