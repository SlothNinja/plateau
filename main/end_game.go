package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
	"github.com/mailjet/mailjet-apiv3-go"
)

// type endGameVPs map[sn.PID]chips
// type crmap map[*datastore.Key]*sn.CurrentRating

func (cl *Client) endGame(c *gin.Context, g *game, uid sn.UID) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	// g.Phase = endGameScoring
	// results := make(endGameVPs, g.NumPlayers)
	// for _, p := range g.players {
	// 	results[p.ID] = make(chips)
	// }
	// g.awardFavorChipPoints(results)
	// g.awardSlanderChipPoints(results)

	// g.newEntry(message{
	// 	"template": "end-game-vp",
	// 	"results":  results,
	// })

	places := g.setFinishOrder()
	cl.Log.Debugf("after setFinishOrder")

	g.Status = sn.Completed

	stats, err := cl.getUStats(c, g.UserIDS...)
	if err != nil {
		sn.JErr(c, err)
		return
	}
	cl.Log.Debugf("after cl.getUStats")

	g.updateUStats(stats)
	cl.Log.Debugf("after g.updateUStats")

	oldElos, newElos, err := cl.Elo.Update(c, g.UserIDS, places)
	if err != nil {
		sn.JErr(c, err)
		return
	}
	cl.Log.Debugf("after cl.Elo.Update")

	g.Undo.Commit()
	g.EndedAt = time.Now()

	_, err = cl.DS.RunInTransaction(c, func(tx *datastore.Transaction) error {
		h := g.Header
		ks := []*datastore.Key{g.headerKey(), g.gameKey(), g.committedKey()}
		es := []interface{}{&h, g, g}

		for _, stat := range stats {
			ks = append(ks, stat.Key)
			es = append(es, stat)
		}

		for _, newElo := range newElos {
			ks = append(ks, newElo.Key, newElo.IncompleteKey())
			es = append(es, newElo, newElo)
		}

		_, err := tx.PutMulti(ks, es)
		if err != nil {
			return err
		}
		return cl.clearCached(c, g, uid)
	})
	if err != nil {
		sn.JErr(c, err)
		return
	}
	cl.Log.Debugf("after cl.DS.RunInTransaction")

	err = cl.sendEndGameNotifications(c, g, oldElos, newElos)
	if err != nil {
		// log but otherwise ignore send errors
		cl.Log.Warningf(err.Error())
	}
	cl.Log.Debugf("after cl.sendEndGameNotifications")
	c.JSON(http.StatusOK, gin.H{"game": g})

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

	sortedByScore(g.players, descending)

	place := 1
	places := make(sn.PlacesMap, len(g.players))
	for i, p1 := range g.players {
		// Update player stats
		p1.stats.Finish = place
		uid1 := g.uidForPID(p1.id)
		places[uid1] = place

		// Update Winners
		if place == 1 {
			g.WinnerIDS = append(g.WinnerIDS, uid1)
		}

		// if next player in order is tied with current player
		// place is not changed.
		// otherwise, place is set to index plus two (one to account for zero index and one to increment place)
		if j := i + 1; j < len(g.players) {
			p2 := g.players[j]
			if p1.compareByScore(p2) != equalTo {
				place = i + 2
			}
		}

	}

	g.newEntry(message{
		"template": "announce-winners",
		"winners":  g.WinnerIDS,
	})
	return places
}

type result struct {
	Place, Rating, Score int
	Name, Inc            string
}

type results []result

func (cl *Client) sendEndGameNotifications(c *gin.Context, g *game, oldElos, newElos []*sn.Elo) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	if g == nil {
		return errors.New("cl.g was nil")
	}

	g.Status = sn.Completed
	rs := make(results, g.NumPlayers)

	for i, p := range g.players {
		rs[i] = result{
			Place:  p.stats.Finish,
			Rating: newElos[i].Rating,
			Score:  p.score,
			Name:   g.NameFor(p.id),
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

	ms := make([]mailjet.InfoMessagesV31, len(g.players))
	subject := fmt.Sprintf("SlothNinja Games: Tammany Hall #%d Has Ended", g.id())
	body := buf.String()
	for i, p := range g.players {
		ms[i] = mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "webmaster@slothninja.com",
				Name:  "Webmaster",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: g.EmailFor(p.id),
					Name:  g.NameFor(p.id),
				},
			},
			Subject:  subject,
			HTMLPart: body,
		}
	}
	_, err = sn.SendMessages(c, ms...)
	return err
}

func (g *game) winnerNames() []string {
	return pie.Map(g.WinnerIDS, func(uid sn.UID) string { return g.NameFor(g.playerByUID(uid).id) })
}
