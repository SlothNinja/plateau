package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
	"github.com/mailjet/mailjet-apiv3-go"
)

func (g game) endGameCheck() bool {
	return g.currentHand() == g.finalHand()
}

func (cl Client) endGame(ctx *gin.Context, g game, cu sn.User) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	places := g.setFinishOrder()
	cl.Log.Debugf("places: %#v", places)

	g.Status = sn.Completed

	stats, err := sn.GetUStats(ctx, cl.FS, maxPlayers, g.UserIDS...)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}
	cl.Log.Debugf("after cl.getUStats")

	g.UpdateUStats(stats, g.playerStats(), g.playerUIDS())

	oldElos, newElos, err := sn.UpdateElo(ctx, cl.FS, g.UserIDS, places)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	g.Undo.Commit()
	g.EndedAt = time.Now()
	err = cl.FS.RunTransaction(ctx, func(c context.Context, tx *firestore.Transaction) error {
		if err := cl.saveGameIn(ctx, tx, g, cu); err != nil {
			return err
		}

		if err := sn.SaveUStatsIn(ctx, cl.FS, tx, stats); err != nil {
			return err
		}

		return sn.SaveElosIn(ctx, cl.FS, tx, newElos)
	})
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = cl.sendEndGameNotifications(ctx, g, oldElos, newElos)
	if err != nil {
		// log but otherwise ignore send errors
		cl.Log.Warningf(err.Error())
	}
	ctx.JSON(http.StatusOK, nil)

}

// func (g *game) awardFavorChipPoints(results endGameVPs) {
// 	for _, n := range nationalities() {
// 		for _, p := range g.chipLeaders(n) {
// 			p.Score += 2
// 			results[p.ID][n] = 2
// 		}
// 	}
// }
//
// func (g *game) chipLeaders(n nationality) []*player {
// 	max := -1
// 	var leaders []*player
// 	for _, p := range g.players {
// 		switch chips := p.chipsFor(n); {
// 		case chips > max:
// 			max = chips
// 			leaders = []*player{p}
// 		case chips == max:
// 			leaders = append(leaders, p)
// 		}
// 	}
// 	return leaders
// }
//
// func (g *game) awardSlanderChipPoints(results endGameVPs) {
// 	for _, p := range g.players {
// 		var slanderVP int
// 		for _, chip := range p.SlanderChips {
// 			if chip {
// 				slanderVP++
// 			}
// 		}
// 		p.Score += slanderVP
// 		results[p.ID]["slander"] = slanderVP
// 	}
// }

func (g *game) setFinishOrder() sn.PlacesMap {
	g.Phase = announceWinners
	g.Status = sn.Completed

	// Set to no current player
	g.setCurrentPlayers()

	sortedByScore(g.Players, descending)

	place := 1
	places := make(sn.PlacesMap, len(g.Players))
	for i, p1 := range g.Players {
		// Update player stats
		p1.Stats.Finish = place
		uid1 := g.uidForPID(p1.ID)
		places[uid1] = place

		// Update Winners
		if place == 1 {
			g.WinnerIDS = append(g.WinnerIDS, uid1)
		}

		// if next player in order is tied with current player
		// place is not changed.
		// otherwise, place is set to index plus two (one to account for zero index and one to increment place)
		if j := i + 1; j < len(g.Players) {
			p2 := g.Players[j]
			if p1.compareByScore(p2) != equalTo {
				place = i + 2
			}
		}

	}

	// g.newEntry(message{
	// 	"template": "announce-winners",
	// 	"winners":  g.WinnerIDS,
	// })
	return places
}

type result struct {
	Place  int
	Rating int
	Score  int64
	Name   string
	Inc    string
}

type results []result

func (cl Client) sendEndGameNotifications(ctx *gin.Context, g game, oldElos, newElos []sn.Elo) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id := getID(ctx)
	g.Status = sn.Completed
	rs := make(results, g.NumPlayers)

	for i, p := range g.Players {
		rs[i] = result{
			Place:  p.Stats.Finish,
			Rating: newElos[i].Rating,
			Score:  p.Score,
			Name:   g.NameFor(p.ID),
			Inc:    fmt.Sprintf("%+d", newElos[i].Rating-oldElos[i].Rating),
		}
	}

	buf := new(bytes.Buffer)
	tmpl := template.New("end_game_notification")
	tmpl, err := tmpl.Parse(`
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html>
        <head>
                <meta http-equiv="content-type" content="text/html; charset=ISO-8859-1">
        </head>
        <body bgcolor="#ffffff" text="#000000">
                {{range $i, $r := $.Results}}
                <div style="height:3em">
                        <div style="height:3em;float:left;padding-right:1em">{{$r.Place}}.</div>
                        <div style="height:1em">{{$r.Name}} scored {{$r.Score}} points.</div>
                        <div style="height:1em">Elo {{$r.Inc}} (-> {{$r.Rating}})</div>
                </div>
                {{end}}
                <p></p>
                <p>Congratulations: {{$.Winners}}.</p>
        </body>
</html>`)
	if err != nil {
		return err
	}

	err = tmpl.Execute(buf, gin.H{
		"Results": rs,
		"Winners": sn.ToSentence(g.winnerNames()),
	})
	if err != nil {
		return err
	}

	ms := make([]mailjet.InfoMessagesV31, len(g.Players))
	subject := fmt.Sprintf("SlothNinja Games: Tammany Hall (%s) Has Ended", id)
	body := buf.String()
	for i, p := range g.Players {
		ms[i] = mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "webmaster@slothninja.com",
				Name:  "Webmaster",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: g.EmailFor(p.ID),
					Name:  g.NameFor(p.ID),
				},
			},
			Subject:  subject,
			HTMLPart: body,
		}
	}
	_, err = sn.SendMessages(ctx, ms...)
	return err
}

func (g game) winnerNames() []string {
	return pie.Map(g.WinnerIDS, func(uid sn.UID) string { return g.NameFor(g.playerByUID(uid).ID) })
}
