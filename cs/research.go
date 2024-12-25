package cs

type TechField string

const (
	TechFieldNone TechField = ""
	Energy        TechField = "Energy"
	Weapons       TechField = "Weapons"
	Propulsion    TechField = "Propulsion"
	Construction  TechField = "Construction"
	Electronics   TechField = "Electronics"
	Biotechnology TechField = "Biotechnology"
)

var TechFields []TechField = []TechField{
	Energy,
	Weapons,
	Propulsion,
	Construction,
	Electronics,
	Biotechnology,
}

type NextResearchField string

const (
	NextResearchFieldSameField     NextResearchField = "SameField"
	NextResearchFieldEnergy        NextResearchField = "Energy"
	NextResearchFieldWeapons       NextResearchField = "Weapons"
	NextResearchFieldPropulsion    NextResearchField = "Propulsion"
	NextResearchFieldConstruction  NextResearchField = "Construction"
	NextResearchFieldElectronics   NextResearchField = "Electronics"
	NextResearchFieldBiotechnology NextResearchField = "Biotechnology"
	NextResearchFieldLowestField   NextResearchField = "LowestField"
)

// The researcher interface is used during turn generation to research for a player
// and research with splash damage/stolen resources
type researcher interface {
	research(player *Player, resourcesToSpend int, onLevelGained func(player *Player, field TechField)) (spent TechLevel)
	researchField(player *Player, field TechField, resourcesToSpend int, onLevelGained func(player *Player, field TechField))
	getTotalCost(techLevels TechLevel, field TechField, researchCostLevel ResearchCostLevel, level int) int
}

type research struct {
	rules *Rules
}

func NewResearcher(rules *Rules) researcher { return &research{rules} }

// This function will be called repeatedly until no more levels are passed
// From starsfaq
//
//	 The cost of a tech level depends on four things:
//	1) Your research setting for that field (cheap, normal, or expensive)
//	2) The level you are researching (higher level, higher cost)
//	3) The total number of tech levels you have already have in all fields (you can add it up yourself, or look at 'tech levels' on the 'score' screen).
//	4) whether 'slow tech advance' was selected as a game parameter.
//
//	in general,
//
//	totalCost=(baseCost + (totalLevels * 10)) * costFactor
//
//	where  totalLevels=the sum of your current levels in all fields
//	  costFactor =.5 if your setting for the field is '50% less'
//	                     =1 if your setting for the field is 'normal'
//	                     =1.75 if your setting for the field is '75% more expensive'
//
//	If 'slow tech advance' is a game parameter, totalCost should be doubled.
//
//	Below is a table showing the base cost of each level.
//
//	1     50              14    18040
//	2     80              15    22440
//	3     130             16    27050
//	4     210             17    31870
//	5     340             18    36900
//	6     550             19    42140
//	7     890             20    47590
//	8     1440            21    53250
//	9     2330            22    59120
//	10    3770            23    65200
//	11    6100            24    71490
//	12    9870            25    77990
//	13    13850           26    84700
func (r *research) research(player *Player, resourcesToSpend int, onLevelGained func(player *Player, field TechField)) (spent TechLevel) {
	// keep spending resources until we are done
	for {
		field := player.Researching

		if field == TechFieldNone {
			// nothing to research, break out
			// this happens if we reach max level
			break
		}
		levelGained, leftoverResources := r.researchFieldOnce(player, player.Researching, resourcesToSpend)
		spent.Set(field, spent.Get(field)+resourcesToSpend-leftoverResources)

		// we gained a level, switch to a new field
		if levelGained {
			player.Researching = r.getNextResearchField(player)
			if player.NextResearchField != NextResearchFieldLowestField {
				player.NextResearchField = NextResearchFieldSameField
			}
			resourcesToSpend = leftoverResources
			onLevelGained(player, field)
		} else {
			// break out
			break
		}
	}

	return spent
}

// keep researching this field until we run out of resources or it maxes out
func (r *research) researchField(player *Player, field TechField, resourcesToSpend int, onLevelGained func(player *Player, field TechField)) {
	for {
		// keep researching this field until we max it out or run out of resources to spend
		if resourcesToSpend == 0 || r.isAtMaxLevel(player, field) {
			break
		}

		levelGained, leftoverResources := r.researchFieldOnce(player, field, resourcesToSpend)
		if levelGained {
			onLevelGained(player, field)
		}
		// keep spending until we're out of resources
		resourcesToSpend = leftoverResources
	}

}

func (r *research) researchFieldOnce(player *Player, field TechField, resourcesToSpend int) (levelGained bool, resourcesLeftover int) {
	if r.isAtMaxLevel(player, field) {
		return false, resourcesToSpend
	}

	// don't research more than the max on this level
	level := player.TechLevels.Get(field)

	// add the resourcesToSpend to how much we've currently spent
	spent := player.TechLevelsSpent.Get(field)

	totalCost := r.getTotalCost(player.TechLevels, field, player.Race.ResearchCost.Get(field), level)

	if spent+resourcesToSpend >= totalCost {
		// increase a level
		levelGained = true

		// gain a level, set our total cost
		player.TechLevels.Set(field, level+1)
		player.TechLevelsSpent.Set(field, 0)

		spent += resourcesToSpend

		// figure out how many leftover points we have
		resourcesLeftover = spent - totalCost
	} else {
		// didn't gain a level, apply resources to field
		player.TechLevelsSpent.Set(field, spent+resourcesToSpend)
		resourcesLeftover = 0
	}

	return levelGained, resourcesLeftover
}

func (r *research) getTotalCost(techLevels TechLevel, field TechField, researchCostLevel ResearchCostLevel, level int) int {
	maxTechLevel := len(r.rules.TechBaseCost) - 1
	// we can't research more than tech level 26
	if level >= maxTechLevel {
		return 0
	}

	// figure out our total levels
	totalLevels := techLevels.Sum()

	// figure out the cost to advance to the next level
	baseCost := r.rules.TechBaseCost[level+1]
	costFactor := 1.0
	switch researchCostLevel {
	case ResearchCostExtra:
		costFactor = 1.75
	case ResearchCostLess:
		costFactor = .5
	}

	// from starsfaq
	return int(float64(baseCost+(totalLevels*10)) * costFactor)

}

// check if a player is at max research level for their current field
func (r *research) isAtMaxLevel(player *Player, field TechField) bool {
	maxTechLevel := len(r.rules.TechBaseCost) - 1
	return player.TechLevels.Get(field) >= maxTechLevel
}

// get the next TechField to research based on the NextResearchField setting
func (r *research) getNextResearchField(player *Player) (nextField TechField) {

	// find the next field
	nextField = Energy
	switch player.NextResearchField {
	case NextResearchFieldSameField:
		nextField = player.Researching
	case NextResearchFieldEnergy:
		nextField = Energy
	case NextResearchFieldWeapons:
		nextField = Weapons
	case NextResearchFieldPropulsion:
		nextField = Propulsion
	case NextResearchFieldConstruction:
		nextField = Construction
	case NextResearchFieldElectronics:
		nextField = Electronics
	case NextResearchFieldBiotechnology:
		nextField = Biotechnology
	case NextResearchFieldLowestField:
		nextField = player.TechLevels.Lowest()
	}

	// if the player is at the max level for this nextField, pick the lowest.
	// if they are at the maxLevel for the lowest, return none
	if r.isAtMaxLevel(player, nextField) {
		nextField = player.TechLevels.Lowest()
		// determine the next field to research
		if r.isAtMaxLevel(player, nextField) {
			return TechFieldNone
		}
	}

	return nextField
}
