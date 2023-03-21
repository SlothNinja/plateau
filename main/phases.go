package main

type phase string

const (
	noPhase                 phase = ""
	actions                 phase = "Actions"
	placeImmigrant          phase = "Place Immigrant"
	elections               phase = "Elections"
	takeFavorChip           phase = "Take Favor Chip"
	endGameScoring          phase = "End Game Scoring"
	scoreVictoryPoints      phase = "Score Victory Points"
	assignCityOffices       phase = "Assign City Offices"
	awardFavorChips         phase = "Award Favor Chips"
	setup                   phase = "Setup"
	castleGarden            phase = "Castle Garden"
	announceWinners         phase = "Announce Winners"
	gameOver                phase = "Game Over"
	assignDeputyMayor       phase = "Assign Deputy Mayor"
	deputyMayorAssignOffice phase = "Deputy Mayor Assigns Office"
)
