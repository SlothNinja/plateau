package main

import "github.com/SlothNinja/sn/v3"

const (
	noPhase           sn.Phase = ""
	setupPhase        sn.Phase = "setup"
	dealPhase         sn.Phase = "deal"
	bidPhase          sn.Phase = "bid"
	exchangePhase     sn.Phase = "card exchange"
	pickPartnerPhase  sn.Phase = "pick parner"
	incObjectivePhase sn.Phase = "increase objective"
	cardPlayPhase     sn.Phase = "card play"
	endHandPhase      sn.Phase = "end of hand"
	announceWinners   sn.Phase = "announce winners"
)
