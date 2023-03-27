package main

import "github.com/elliotchance/pie/v2"

type deck []card

// shuffled deck for number players
func deckFor(numPlayers int) deck {
	var d deck
	switch numPlayers {
	case 2:
		d = deck{
			card{rank: oneRank, suit: trumps},
			card{rank: twoRank, suit: trumps},
			card{rank: threeRank, suit: trumps},
			card{rank: fourRank, suit: trumps},
			card{rank: fiveRank, suit: trumps},
			card{rank: sixRank, suit: trumps},
			card{rank: sevenRank, suit: trumps},
			card{rank: tenRank, suit: hearts}, card{rank: tenRank, suit: clubs}, card{rank: tenRank, suit: spades}, card{rank: tenRank, suit: diamonds},
			card{rank: valetRank, suit: hearts}, card{rank: valetRank, suit: clubs}, card{rank: valetRank, suit: spades}, card{rank: valetRank, suit: diamonds},
			card{rank: cavalierRank, suit: hearts}, card{rank: cavalierRank, suit: clubs}, card{rank: cavalierRank, suit: spades}, card{rank: cavalierRank, suit: diamonds},
			card{rank: dameRank, suit: hearts}, card{rank: dameRank, suit: clubs}, card{rank: dameRank, suit: spades}, card{rank: dameRank, suit: diamonds},
			card{rank: roiRank, suit: hearts}, card{rank: roiRank, suit: clubs}, card{rank: roiRank, suit: spades}, card{rank: roiRank, suit: diamonds},
			card{rank: sixteenRank, suit: trumps},
			card{rank: seventeenRank, suit: trumps},
			card{rank: eighteenRank, suit: trumps},
			card{rank: nineteenRank, suit: trumps},
			card{rank: twentyRank, suit: trumps},
			card{rank: twentyoneRank, suit: trumps},
			card{rank: excuseRank, suit: trumps},
		}
	case 3:
		d = deck{
			card{rank: oneRank, suit: trumps},
			card{rank: twoRank, suit: trumps},
			card{rank: threeRank, suit: trumps},
			card{rank: fourRank, suit: trumps},
			card{rank: fiveRank, suit: trumps},
			card{rank: sixRank, suit: trumps},
			card{rank: sevenRank, suit: trumps},
			card{rank: eightRank, suit: trumps},
			card{rank: nineRank, suit: trumps},
			card{rank: tenRank, suit: trumps}, card{rank: tenRank, suit: hearts}, card{rank: tenRank, suit: clubs}, card{rank: tenRank, suit: spades}, card{rank: tenRank, suit: diamonds},
			card{rank: elevenRank, suit: trumps}, card{rank: valetRank, suit: hearts}, card{rank: valetRank, suit: clubs}, card{rank: valetRank, suit: spades}, card{rank: valetRank, suit: diamonds},
			card{rank: twelveRank, suit: trumps}, card{rank: cavalierRank, suit: hearts}, card{rank: cavalierRank, suit: clubs}, card{rank: cavalierRank, suit: spades}, card{rank: cavalierRank, suit: diamonds},
			card{rank: thirteenRank, suit: trumps}, card{rank: dameRank, suit: hearts}, card{rank: dameRank, suit: clubs}, card{rank: dameRank, suit: spades}, card{rank: dameRank, suit: diamonds},
			card{rank: fourteenRank, suit: trumps}, card{rank: roiRank, suit: hearts}, card{rank: roiRank, suit: clubs}, card{rank: roiRank, suit: spades}, card{rank: roiRank, suit: diamonds},
			card{rank: fifteenRank, suit: trumps},
			card{rank: sixteenRank, suit: trumps},
			card{rank: seventeenRank, suit: trumps},
			card{rank: eighteenRank, suit: trumps},
			card{rank: nineteenRank, suit: trumps},
			card{rank: twentyRank, suit: trumps},
			card{rank: twentyoneRank, suit: trumps},
			card{rank: excuseRank, suit: trumps},
		}
	case 4:
		d = deck{
			card{rank: oneRank, suit: trumps},
			card{rank: twoRank, suit: trumps},
			card{rank: threeRank, suit: trumps},
			card{rank: fourRank, suit: trumps},
			card{rank: fiveRank, suit: trumps},
			card{rank: sixRank, suit: trumps}, card{rank: sixRank, suit: spades},
			card{rank: sevenRank, suit: trumps}, card{rank: sevenRank, suit: hearts}, card{rank: sevenRank, suit: clubs}, card{rank: sevenRank, suit: spades}, card{rank: sevenRank, suit: diamonds},
			card{rank: eightRank, suit: trumps}, card{rank: eightRank, suit: hearts}, card{rank: eightRank, suit: clubs}, card{rank: eightRank, suit: spades}, card{rank: eightRank, suit: diamonds},
			card{rank: nineRank, suit: trumps}, card{rank: nineRank, suit: hearts}, card{rank: nineRank, suit: clubs}, card{rank: nineRank, suit: spades}, card{rank: nineRank, suit: diamonds},
			card{rank: tenRank, suit: trumps}, card{rank: tenRank, suit: hearts}, card{rank: tenRank, suit: clubs}, card{rank: tenRank, suit: spades}, card{rank: tenRank, suit: diamonds},
			card{rank: elevenRank, suit: trumps}, card{rank: valetRank, suit: hearts}, card{rank: valetRank, suit: clubs}, card{rank: valetRank, suit: spades}, card{rank: valetRank, suit: diamonds},
			card{rank: twelveRank, suit: trumps}, card{rank: cavalierRank, suit: hearts}, card{rank: cavalierRank, suit: clubs}, card{rank: cavalierRank, suit: spades}, card{rank: cavalierRank, suit: diamonds},
			card{rank: thirteenRank, suit: trumps}, card{rank: dameRank, suit: hearts}, card{rank: dameRank, suit: clubs}, card{rank: dameRank, suit: spades}, card{rank: dameRank, suit: diamonds},
			card{rank: fourteenRank, suit: trumps}, card{rank: roiRank, suit: hearts}, card{rank: roiRank, suit: clubs}, card{rank: roiRank, suit: spades}, card{rank: roiRank, suit: diamonds},
			card{rank: fifteenRank, suit: trumps},
			card{rank: sixteenRank, suit: trumps},
			card{rank: seventeenRank, suit: trumps},
			card{rank: eighteenRank, suit: trumps},
			card{rank: nineteenRank, suit: trumps},
			card{rank: twentyRank, suit: trumps},
			card{rank: twentyoneRank, suit: trumps},
			card{rank: excuseRank, suit: trumps},
		}
	case 5:
		d = deck{
			card{rank: oneRank, suit: trumps},
			card{rank: twoRank, suit: trumps},
			card{rank: threeRank, suit: trumps}, card{rank: threeRank, suit: clubs}, card{rank: threeRank, suit: spades},
			card{rank: fourRank, suit: trumps}, card{rank: fourRank, suit: hearts}, card{rank: fourRank, suit: clubs}, card{rank: fourRank, suit: spades}, card{rank: fourRank, suit: diamonds},
			card{rank: fiveRank, suit: trumps}, card{rank: fiveRank, suit: hearts}, card{rank: fiveRank, suit: clubs}, card{rank: fiveRank, suit: spades}, card{rank: fiveRank, suit: diamonds},
			card{rank: sixRank, suit: trumps}, card{rank: sixRank, suit: hearts}, card{rank: sixRank, suit: clubs}, card{rank: sixRank, suit: spades}, card{rank: sixRank, suit: diamonds},
			card{rank: sevenRank, suit: trumps}, card{rank: sevenRank, suit: hearts}, card{rank: sevenRank, suit: clubs}, card{rank: sevenRank, suit: spades}, card{rank: sevenRank, suit: diamonds},
			card{rank: eightRank, suit: trumps}, card{rank: eightRank, suit: hearts}, card{rank: eightRank, suit: clubs}, card{rank: eightRank, suit: spades}, card{rank: eightRank, suit: diamonds},
			card{rank: nineRank, suit: trumps}, card{rank: nineRank, suit: hearts}, card{rank: nineRank, suit: clubs}, card{rank: nineRank, suit: spades}, card{rank: nineRank, suit: diamonds},
			card{rank: tenRank, suit: trumps}, card{rank: tenRank, suit: hearts}, card{rank: tenRank, suit: clubs}, card{rank: tenRank, suit: spades}, card{rank: tenRank, suit: diamonds},
			card{rank: elevenRank, suit: trumps}, card{rank: valetRank, suit: hearts}, card{rank: valetRank, suit: clubs}, card{rank: valetRank, suit: spades}, card{rank: valetRank, suit: diamonds},
			card{rank: twelveRank, suit: trumps}, card{rank: cavalierRank, suit: hearts}, card{rank: cavalierRank, suit: clubs}, card{rank: cavalierRank, suit: spades}, card{rank: cavalierRank, suit: diamonds},
			card{rank: thirteenRank, suit: trumps}, card{rank: dameRank, suit: hearts}, card{rank: dameRank, suit: clubs}, card{rank: dameRank, suit: spades}, card{rank: dameRank, suit: diamonds},
			card{rank: fourteenRank, suit: trumps}, card{rank: roiRank, suit: hearts}, card{rank: roiRank, suit: clubs}, card{rank: roiRank, suit: spades}, card{rank: roiRank, suit: diamonds},
			card{rank: fifteenRank, suit: trumps},
			card{rank: sixteenRank, suit: trumps},
			card{rank: seventeenRank, suit: trumps},
			card{rank: eighteenRank, suit: trumps},
			card{rank: nineteenRank, suit: trumps},
			card{rank: twentyRank, suit: trumps},
			card{rank: twentyoneRank, suit: trumps},
			card{rank: excuseRank, suit: trumps},
		}
	case 6:
		d = deck{
			card{rank: oneRank, suit: trumps}, card{rank: oneRank, suit: hearts}, card{rank: oneRank, suit: clubs}, card{rank: oneRank, suit: spades}, card{rank: oneRank, suit: diamonds},
			card{rank: twoRank, suit: trumps}, card{rank: twoRank, suit: hearts}, card{rank: twoRank, suit: clubs}, card{rank: twoRank, suit: spades}, card{rank: twoRank, suit: diamonds},
			card{rank: threeRank, suit: trumps}, card{rank: threeRank, suit: hearts}, card{rank: threeRank, suit: clubs}, card{rank: threeRank, suit: spades}, card{rank: threeRank, suit: diamonds},
			card{rank: fourRank, suit: trumps}, card{rank: fourRank, suit: hearts}, card{rank: fourRank, suit: clubs}, card{rank: fourRank, suit: spades}, card{rank: fourRank, suit: diamonds},
			card{rank: fiveRank, suit: trumps}, card{rank: fiveRank, suit: hearts}, card{rank: fiveRank, suit: clubs}, card{rank: fiveRank, suit: spades}, card{rank: fiveRank, suit: diamonds},
			card{rank: sixRank, suit: trumps}, card{rank: sixRank, suit: hearts}, card{rank: sixRank, suit: clubs}, card{rank: sixRank, suit: spades}, card{rank: sixRank, suit: diamonds},
			card{rank: sevenRank, suit: trumps}, card{rank: sevenRank, suit: hearts}, card{rank: sevenRank, suit: clubs}, card{rank: sevenRank, suit: spades}, card{rank: sevenRank, suit: diamonds},
			card{rank: eightRank, suit: trumps}, card{rank: eightRank, suit: hearts}, card{rank: eightRank, suit: clubs}, card{rank: eightRank, suit: spades}, card{rank: eightRank, suit: diamonds},
			card{rank: nineRank, suit: trumps}, card{rank: nineRank, suit: hearts}, card{rank: nineRank, suit: clubs}, card{rank: nineRank, suit: spades}, card{rank: nineRank, suit: diamonds},
			card{rank: tenRank, suit: trumps}, card{rank: tenRank, suit: hearts}, card{rank: tenRank, suit: clubs}, card{rank: tenRank, suit: spades}, card{rank: tenRank, suit: diamonds},
			card{rank: elevenRank, suit: trumps}, card{rank: valetRank, suit: hearts}, card{rank: valetRank, suit: clubs}, card{rank: valetRank, suit: spades}, card{rank: valetRank, suit: diamonds},
			card{rank: twelveRank, suit: trumps}, card{rank: cavalierRank, suit: hearts}, card{rank: cavalierRank, suit: clubs}, card{rank: cavalierRank, suit: spades}, card{rank: cavalierRank, suit: diamonds},
			card{rank: thirteenRank, suit: trumps}, card{rank: dameRank, suit: hearts}, card{rank: dameRank, suit: clubs}, card{rank: dameRank, suit: spades}, card{rank: dameRank, suit: diamonds},
			card{rank: fourteenRank, suit: trumps}, card{rank: roiRank, suit: hearts}, card{rank: roiRank, suit: clubs}, card{rank: roiRank, suit: spades}, card{rank: roiRank, suit: diamonds},
			card{rank: fifteenRank, suit: trumps},
			card{rank: sixteenRank, suit: trumps},
			card{rank: seventeenRank, suit: trumps},
			card{rank: eighteenRank, suit: trumps},
			card{rank: nineteenRank, suit: trumps},
			card{rank: twentyRank, suit: trumps},
			card{rank: twentyoneRank, suit: trumps},
			card{rank: excuseRank, suit: trumps},
		}
	default:
		d = nil
	}

	return pie.Shuffle(d, myRandomSource)
}
