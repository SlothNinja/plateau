package main

import (
	"context"
	log2 "log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/sn/v3"
	"google.golang.org/api/option"
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

	sessionName           = "sng-oauth"
	invitationsPath       = "invitations"
	selectWardPath        = "selectWard/:id"
	placePiecesPath       = "place/:id"
	removeImmigrantPath   = "remove/:id"
	moveImmigrantFromPath = "move-from/:id"
	moveImmigrantToPath   = "move-to/:id"
	lockUpPath            = "lock-up/:id"
	slanderPath           = "slander/:id"
	takeChipPath          = "takeChip/:id"
	deputyTakeChipPath    = "deputyTakeChip/:id"
	assignOfficesPath     = "assignOffices/:id"
	subPath               = "/subscribe/:id"
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

// Client provide client structure of the Le Plateau service
type Client struct {
	sn.Client[*game, *invitation, *player]
}

// NewClient returns a new Client for the plateau service
func NewClient(ctx context.Context) Client {
	nClient := Client{sn.NewClient[*game, *invitation, *player](ctx, sn.Options{
		ProjectID:     projectID(),
		UserProjectID: getUserProjectID(),
		UserDSURL:     getUserDSURL(),
		LoggerID:      "plateau",
		CorsAllow:     []string{"https://plateau.fake-slothninja.com:8092/*"},
		Prefix:        "sn",
	})}
	return nClient.addRoutes("sn")
}

// Close closes plateau service Client
func (cl *Client) Close() error {
	return cl.Client.Close()
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

// func newLogClient() *sn.LogClient {
// 	client, err := sn.NewLogClient(projectID())
// 	if err != nil {
// 		log.Panicf("unable to create logging client: %v", err)
// 	}
// 	return client
// }

// AddRoutes addes routing for game.
func (cl Client) addRoutes(prefix string) Client {
	/////////////////////////////////////////////
	// Game Group
	gGroup := cl.Router.Group(prefix + "/game")

	// Place Bid
	gGroup.PUT("bid/:id", cl.CachedHandler((*game).placeBid))

	// Pass Bid
	gGroup.PUT("passBid/:id", cl.CachedHandler((*game).passBid))

	// Increase Objective
	gGroup.PUT("incObjective/:id", cl.CachedHandler((*game).incObjective))

	// Abdicate
	gGroup.PUT("abdicate/:id", cl.CachedHandler((*game).abdicate))

	// Card Exchange
	gGroup.PUT("exchange/:id", cl.CachedHandler((*game).exchange))

	// Play Card
	gGroup.PUT("play/:id", cl.CachedHandler((*game).playCard))

	// Pick Partner
	gGroup.PUT("finish/pick/:id", cl.FinishTurnHandler((*game).pickPartner))

	// Bid Finish
	gGroup.PUT("finish/bid/:id", cl.FinishTurnHandler((*game).bidFinishTurn))

	// Exchange Finish
	gGroup.PUT("finish/exchange/:id", cl.FinishTurnHandler((*game).exchangeFinishTurn))

	// IncObjective Finish
	gGroup.PUT("finish/objective/:id", cl.FinishTurnHandler((*game).incObjectiveFinishTurn))

	// Play Card Finish
	gGroup.PUT("finish/play/:id", cl.FinishTurnHandler((*game).playCardFinishTurn))

	return cl
}
