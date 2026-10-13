package ai

import "github.com/sirgwain/craig-stars/cs"

// A doctrine is how a primary racial trait changes the standard play. A zero
// field leaves the standard behavior alone.
type doctrine struct {
	// colonyShips multiplies how many colony ships are in flight or waiting at once
	colonyShips float64
	// seedHabitability replaces the habitability a world needs before an idle
	// colony ship takes an early minimum load to it
	seedHabitability int
	// valueTerraformedWorlds judges a free world by what it will be once
	// settled, for a race that terraforms the moment it lands
	valueTerraformedWorlds bool
	// attackYears is added to the year the first attack is allowed
	attackYears int
	// warFleets multiplies how many war fleets are under construction at once
	warFleets float64
	// outposts gives young colonies a cheap orbital fort early, for the
	// stargate or mass driver it carries
	outposts bool
	// defensesPerYear are built on every world ahead of terraforming and the
	// rest of its factories and mines
	defensesPerYear int
	// mineLayers keeps a mine layer over every owned world
	mineLayers bool
	// cloakedScouts designs scouts with a cloak as well as a scanner
	cloakedScouts bool
}

var doctrines = map[cs.PRT]doctrine{
	// growth is doubled and worlds are small, so spread out as fast as possible
	cs.HE: {colonyShips: 2, seedHabitability: 40},
	// hidden pickets keep watch on the neighbors
	cs.SS: {cloakedScouts: true},
	// the best attackers in the game should attack sooner and with more
	cs.WM: {attackYears: -5, warFleets: 1.5},
	// every world is at its best the moment it is settled
	cs.CA: {colonyShips: 1.5, valueTerraformedWorlds: true},
	// cheap, strong defenses make every world hard to take
	cs.IS: {defensesPerYear: 5},
	cs.SD: {mineLayers: true},
	// a mass driver on every colony extends packet scouting and catches packets
	cs.PP: {outposts: true},
	// a stargate on every colony lets fleets cross the empire in a year
	cs.IT: {outposts: true},
}

// doctrineFor returns the doctrine for the player's race. A test can switch a
// race's doctrine off with the behavior "doctrine-<PRT>".
func (ai *aiPlayer) doctrineFor() doctrine {
	if !ai.enabled("doctrine-" + string(ai.Race.PRT)) {
		return doctrine{}
	}
	return doctrines[ai.Race.PRT]
}

const (
	outpostDevelopment = .2 // development a colony needs before it gets an outpost
	outpostYears       = 5  // an outpost is built if it takes no longer than this
)

// colonyHabitability is the habitability to judge a free world by
func (ai *aiPlayer) colonyHabitability(intel *cs.Planet) int {
	if ai.doctrine.valueTerraformedWorlds && ai.Race.Spec.Instaforming {
		return max(intel.Spec.Habitability, intel.Spec.TerraformedHabitability)
	}
	return intel.Spec.Habitability
}
