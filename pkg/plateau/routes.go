package plateau

import (
	"context"
	"encoding/base64"
	"fmt"
	log2 "log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/sn/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"google.golang.org/api/option"
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

const (
	// PlateauProjectIDEnv provides string of PLATEAU_PROJECT_ID environement variable
	// used to specify project-id of the tammany service
	PlateauProjectIDEnv = "PLATEAU_PROJECT_ID"
	// PlateauDSURLEnv provides string of PLATEAU_DS_URL environement variable
	// used to specify the datastore URL of the tammany service
	PlateauDSURLEnv = "PLATEAU_DS_URL"
	// PlateauHostURLEnv provides string of PLATEAU_HOST_URL environement variable
	// used to specify the host URL of the tammany service
	PlateauHostURLEnv = "PLATEAU_HOST_URL"

	// UserProjectIDEnv provides string of USER_PROJECT_ID environement variable
	// used to specify the project-id of the user service
	UserProjectIDEnv = "USER_PROJECT_ID"
	// UserDSURLEnv provides string of USER_DS_URL environement variable
	// used to specify the datastore URL of the user service
	UserDSURLEnv = "USER_DS_URL"
	// UserHostURLEnv provides string of USER_HOST_URL environement variable
	// used to specify the host URL of the user service
	UserHostURLEnv = "USER_HOST_URL"

	// PlateauCreds provides string of PLATEAU_CREDS environement variable
	// used to specify the credentials for connecting to the client
	PlateauCreds = "PLATEAU_CREDS"

	sessionName              = "sng-oauth"
	invitationPath           = "invitation"
	invitationsPath          = "invitations"
	gamePath                 = "game"
	gamesPath                = "games"
	gamesIndexPath           = ":status"
	showPath                 = "show/:id"
	selectWardPath           = "selectWard/:id"
	placePiecesPath          = "place/:id"
	removeImmigrantPath      = "remove/:id"
	moveImmigrantFromPath    = "move-from/:id"
	moveImmigrantToPath      = "move-to/:id"
	lockUpPath               = "lock-up/:id"
	slanderPath              = "slander/:id"
	takeChipPath             = "takeChip/:id"
	deputyTakeChipPath       = "deputyTakeChip/:id"
	assignOfficesPath        = "assignOffices/:id"
	bidPath                  = "bid/:id"
	actionsFinishPath        = "actionsFinish/:id"
	takeChipFinishPath       = "takeChipFinish/:id"
	assignOfficesFinishPath  = "assignOfficesFinish/:id"
	electionsFinishPath      = "electionsFinish/:id"
	placeImmigrantFinishPath = "placeImmigrantFinish/:id"
	resetPath                = "/reset/:id"
	undoPath                 = "/undo/:id"
	redoPath                 = "/redo/:id"
	rollbackPath             = "/rollback/:id"
	rollforwardPath          = "/rollforward/:id"
	newPath                  = "/new"
	dropPath                 = "/drop/:id"
	acceptPath               = "/accept/:id"
	detailsPath              = "/details/:id"
	subPath                  = "/subscribe/:id"
)

func projectID() string {
	return os.Getenv(PlateauProjectIDEnv)
}

func dsURL() string {
	return os.Getenv(PlateauDSURLEnv)
}

func getUserProjectID() string {
	return os.Getenv(UserProjectIDEnv)
}

func getUserDSURL() string {
	return os.Getenv(UserDSURLEnv)
}

func getUserHostURL() string {
	return os.Getenv(UserHostURLEnv)
}

// Client provide client structure of the tammany service
type Client struct {
	*sn.Client
	User      *sn.UserClient
	MLog      *sn.MLogClient
	Elo       *sn.EloClient
	Messaging *messaging.Client
}

// NewClient returns a new Client for the plateau service
func NewClient(ctx context.Context) *Client {
	logClient := newLogClient()
	snClient := sn.NewClient(ctx, sn.Options{
		ProjectID: projectID(),
		DSURL:     dsURL(),
		Logger:    logClient.Logger("plateau"),
		Cache:     cache.New(30*time.Minute, 10*time.Minute),
		Router:    gin.Default(),
	})

	uClient := sn.NewUserClient(sn.NewClient(ctx, sn.Options{
		ProjectID: getUserProjectID(),
		DSURL:     getUserDSURL(),
		Logger:    snClient.Log,
		Cache:     snClient.Cache,
		Router:    snClient.Router,
	}))

	store, err := sn.NewCookieClient(uClient.Client).NewStore(ctx)
	if err != nil {
		snClient.Log.Panicf("unable create cookie store: %v", err)
	}
	snClient.Router.Use(sessions.Sessions(sessionName, store))

	if !sn.IsProduction() {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{"https://plateau.fake-slothninja.com:8092"}
		config.AllowCredentials = true
		// config.AllowOrigins = []string{"http://google.com", "http://facebook.com"}
		// config.AllowAllOrigins = true
		snClient.Router.Use(cors.New(config))
	}

	nClient := &Client{
		Client:    snClient,
		User:      uClient,
		MLog:      sn.NewMLogClient(snClient, uClient),
		Elo:       sn.NewEloClient(snClient, "elo"),
		Messaging: newMsgClient(ctx),
	}
	return nClient.addRoutes("")
}

// CloseErrors provides struct of errors returned by the multiple clients that make up the tammany service Client
type CloseErrors struct {
	Client     error
	UserClient error
}

// Error provides error interface for CloseErrors
func (ce CloseErrors) Error() string {
	return fmt.Sprintf("error closing clients: client: %q userClient: %q", ce.Client, ce.UserClient)
}

// Close closes tammany service Client
func (cl *Client) Close() error {
	var ce CloseErrors

	ce.Client = cl.Client.Close()
	ce.UserClient = cl.User.Client.Close()

	if ce.Client != nil || ce.UserClient != nil {
		return ce
	}
	return nil
}

func newMsgClient(ctx context.Context) *messaging.Client {
	if sn.IsProduction() {
		log.Debugf("production")
		app, err := firebase.NewApp(ctx, nil)
		if err != nil {
			log2.Panicf("unable to create messaging client: %v", err)
			return nil
		}
		cl, err := app.Messaging(ctx)
		if err != nil {
			log2.Panicf("unable to create messaging client: %v", err)
			return nil
		}
		return cl
	}
	log.Debugf("development")
	app, err := firebase.NewApp(
		ctx,
		nil,
		option.WithGRPCConnectionPool(50),
		option.WithCredentialsFile(os.Getenv(PlateauCreds)),
	)
	if err != nil {
		log2.Panicf("unable to create messaging client: %v", err)
		return nil
	}
	cl, err := app.Messaging(ctx)
	if err != nil {
		log2.Panicf("unable to create messaging client: %v", err)
		return nil
	}
	return cl
}

func newLogClient() *sn.LogClient {
	client, err := sn.NewLogClient(projectID())
	if err != nil {
		log.Panicf("unable to create logging client: %v", err)
	}
	return client
}

func (cl *Client) loginHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	referer := c.Request.Referer()
	encodedReferer := base64.StdEncoding.EncodeToString([]byte(referer))

	path := getUserHostURL() + "/login?redirect=" + encodedReferer
	cl.Log.Debugf("path: %q", path)
	c.Redirect(http.StatusSeeOther, path)
}

func (cl *Client) logoutHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	referer := c.Request.Referer()
	sn.Logout(c)
	c.Redirect(http.StatusSeeOther, referer)
}

// AddRoutes addes routing for game.
func (cl *Client) addRoutes(prefix string) *Client {
	////////////////////////////////////////////
	// Home
	cl.Router.GET("sn/home", cl.homeHandler)

	// 	////////////////////////////////////////////
	// 	// Invitation Group
	// 	inv := cl.Router.Group(invitationPath)
	//
	// 	// New
	// 	inv.GET(newPath, cl.newInvitationHandler)
	//
	// 	// Create
	// 	inv.PUT(newPath, cl.createHandler)
	//
	// 	// Drop
	// 	inv.PUT(dropPath, cl.dropHandler)
	//
	// 	// Accept
	// 	inv.PUT(acceptPath, cl.acceptHandler)
	//
	// 	// Details
	// 	inv.GET(detailsPath, cl.detailsHandler)
	//
	// 	/////////////////////////////////////////////
	// 	// Invitations Group
	// 	invs := cl.Router.Group(invitationsPath)
	//
	// 	// Index
	// 	invs.POST("", cl.invitationsIndexHandler)
	//
	// 	/////////////////////////////////////////////
	// 	// Games Group
	// 	gs := cl.Router.Group(gamesPath)
	//
	// 	// JSON Data for Index
	// 	gs.POST(gamesIndexPath, cl.gamesIndex)
	//
	// 	/////////////////////////////////////////////
	// 	// Game Group
	// 	g := cl.Router.Group(gamePath)
	//
	// 	// Show
	// 	g.GET(showPath, cl.showHandler)
	//
	// 	// // Select Ward
	// 	// g.PUT(selectWardPath, cl.selectWard)
	//
	// 	// Place Pieces
	// 	g.PUT(placePiecesPath, cl.placePiecesHandler)
	//
	// 	// Remove Immigrant
	// 	g.PUT(removeImmigrantPath, cl.removeImmigrantHandler)
	//
	// 	// Move Immigrant From
	// 	g.PUT(moveImmigrantFromPath, cl.moveImmigrantFromHandler)
	//
	// 	// Move Immigrant To
	// 	g.PUT(moveImmigrantToPath, cl.moveImmigrantToHandler)
	//
	// 	// Lock Up
	// 	g.PUT(lockUpPath, cl.lockUpHandler)
	//
	// 	// Slander Boss
	// 	g.PUT(slanderPath, cl.slanderHandler)
	//
	// 	// Take Chip
	// 	g.PUT(takeChipPath, cl.takeChipHandler)
	//
	// 	// Deputy Take Chip
	// 	g.PUT(deputyTakeChipPath, cl.deputyTakeChipHandler)
	//
	// 	// Assign Offices
	// 	g.PUT(assignOfficesPath, cl.assignOfficesHandler)
	//
	// 	// Place Bid
	// 	g.PUT(bidPath, cl.bidHandler)
	//
	// 	// Actions Finish
	// 	g.PUT(actionsFinishPath, cl.actionsFinishTurnHandler)
	//
	// 	// Take Chip Finish
	// 	g.PUT(takeChipFinishPath, cl.takeChipFinishTurnHandler)
	//
	// 	// Elections Finish
	// 	g.PUT(electionsFinishPath, cl.electionsFinishTurnHandler)
	//
	// 	// Place Immigrant Finish
	// 	g.PUT(placeImmigrantFinishPath, cl.placeImmigrantFinishTurnHandler)
	//
	// 	// Assign Offices Finish
	// 	g.PUT(assignOfficesFinishPath, cl.assignOfficesFinishTurnHandler)
	//
	// 	// Reset
	// 	g.PUT(resetPath, cl.resetHandler)
	//
	// 	// Undo
	// 	g.PUT(undoPath, cl.undoHandler)
	//
	// 	// Redo
	// 	g.PUT(redoPath, cl.redoHandler)
	//
	// 	// Rollback
	// 	g.PUT(rollbackPath, cl.rollbackHandler)
	//
	// 	// Rollforward
	// 	g.PUT(rollforwardPath, cl.rollforwardHandler)
	//
	// 	// Sub
	// 	g.PUT(subPath, cl.subHandler)
	//
	// login
	cl.Router.GET("sn/login", cl.loginHandler)
	//
	// logout
	cl.Router.GET("sn/logout", cl.logoutHandler)
	//
	// 	////////////////////////////////////////////
	// 	// Message Log
	// 	msg := cl.Router.Group("mlog")
	//
	// 	// Get
	// 	msg.GET("/:id", cl.mlogHandler)
	//
	// 	// Add
	// 	msg.PUT("/:id/add", cl.mlogAddHandler)
	//
	return cl
}
