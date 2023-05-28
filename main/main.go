package main

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

var myRandomSource = rand.NewSource(time.Now().UnixNano())

func main() {
	ctx := context.Background()

	if sn.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		cl := NewClient(ctx)
		defer cl.Close()
		cl.Router.TrustedPlatform = gin.PlatformGoogleAppEngine
		cl.Router.Run()
	} else {
		gin.SetMode(gin.DebugMode)
		cl := NewClient(ctx)
		defer cl.Close()
		cl.Router.SetTrustedProxies(nil)
		cl.Router.RunTLS(getPort(), "cert.pem", "key.pem")
	}
}

func getPort() string {
	return ":" + os.Getenv("PORT")
}
