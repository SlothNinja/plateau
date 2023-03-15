package main

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/SlothNinja/plateau/pkg/plateau"
	"github.com/SlothNinja/sn/v2"
	"github.com/gin-gonic/gin"
)

const (
	gameKey   = "Game"
	homePath  = "/home"
	jsonKey   = "JSON"
	statusKey = "Status"
	hParam    = "hid"
	msgEnter  = "Entering"
	msgExit   = "Exiting"
)

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()

	if sn.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		cl := plateau.NewClient(ctx)
		defer cl.Close()
		cl.Router.TrustedPlatform = gin.PlatformGoogleAppEngine
		cl.Router.Run()
	} else {
		gin.SetMode(gin.DebugMode)
		cl := staticRoutes(plateau.NewClient(ctx))
		defer cl.Close()
		cl.Router.SetTrustedProxies(nil)
		cl.Router.RunTLS(getPort(), "cert.pem", "key.pem")
	}
}

func getPort() string {
	return ":" + os.Getenv("PORT")
}

func staticRoutes(cl *plateau.Client) *plateau.Client {
	if sn.IsProduction() {
		return cl
	}
	cl.Router.StaticFile("/", "dist/index.html")
	cl.Router.StaticFile("/index.html", "dist/index.html")
	cl.Router.StaticFile("/firebase-messaging-sw.js", "dist/firebase-messaging-sw.js")
	cl.Router.StaticFile("/manifest.json", "dist/manifest.json")
	cl.Router.StaticFile("/robots.txt", "dist/robots.txt")
	cl.Router.StaticFile("/precache-manifest.c0be88927a8120cb7373cf7df05f5688.js", "dist/precache-manifest.c0be88927a8120cb7373cf7df05f5688.js")
	cl.Router.StaticFile("/app.js", "dist/app.js")
	cl.Router.StaticFile("/favicon.ico", "dist/favicon.ico")
	cl.Router.Static("/img", "dist/img")
	cl.Router.Static("/js", "dist/js")
	cl.Router.Static("/css", "dist/css")
	return cl
}

// func (cl *Client) homeHandler(c *gin.Context) {
// 	cl.Log.Debugf(msgEnter)
// 	defer cl.Log.Debugf(msgExit)
//
// 	cu, err := cl.User.Current(c)
// 	if err != nil {
// 		cl.Log.Warningf(err.Error())
// 	}
//
// 	cl.Log.Debugf("cu: %#v", cu)
//
// 	c.JSON(http.StatusOK, gin.H{"cu": cu})
// }