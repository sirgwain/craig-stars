package ai

import (
	"github.com/sirgwain/craig-stars/cs"
)

func (ai *aiPlayer) research() {
	if ai.Race.Spec.ResearchSplashDamage > 0 {
		switch ai.Race.PRT {
		case cs.PP:
			ai.Researching = cs.Energy
			ai.NextResearchField = cs.NextResearchFieldSameField
		case cs.AR:
			if ai.Player.TechLevels.Get(cs.Energy) < 10 {
				// get them resources!
				ai.Researching = cs.Energy
				ai.NextResearchField = cs.NextResearchFieldSameField
			} else {
				// get them starbases
				ai.Researching = cs.Construction
				ai.NextResearchField = cs.NextResearchFieldSameField
			}
		case cs.CA:
			// get that terraforming
			if ai.Player.TechLevels.Get(cs.Biotechnology) < 13 {
				ai.Researching = cs.Biotechnology
				ai.NextResearchField = cs.NextResearchFieldSameField
			} else {
				ai.Researching = cs.Weapons
				ai.NextResearchField = cs.NextResearchFieldSameField
			}
		default:
			// everyone else with GR, get weapons and let the rest come natuarually
			ai.Researching = cs.Weapons
			ai.NextResearchField = cs.NextResearchFieldSameField
		}
		// we have generalized research, just pick weapons and run with it
	} else {
		switch ai.Race.PRT {
		case cs.AR:
			if ai.Player.TechLevels.Get(cs.Energy) < 10 {
				// get them resources!
				ai.Researching = cs.Energy
				ai.NextResearchField = cs.NextResearchFieldSameField
			} else {
				// get them starbases
				ai.NextResearchField = cs.NextResearchFieldLowestField
			}
		case cs.CA:
			if ai.Race.HasLRT(cs.TT) && ai.TechLevels.Get(cs.Biotechnology) < 12 {
				// get that terraforming
				ai.Researching = cs.Biotechnology
				ai.NextResearchField = cs.NextResearchFieldSameField
			} else {
				ai.NextResearchField = cs.NextResearchFieldLowestField
			}
		default:
			// everyone else just go with the lowest
			ai.NextResearchField = cs.NextResearchFieldLowestField
		}
	}

	// if we maxed our goal tech, move to the next one
	if ai.TechLevels.Get(ai.Researching) == ai.game.Rules.MaxTechLevel {
		ai.NextResearchField = cs.NextResearchFieldLowestField
		// find a field that isn't maxed
		for _, field := range cs.TechFields {
			if ai.TechLevels.Get(field) < ai.game.Rules.MaxTechLevel {
				ai.Researching = field
				break
			}
		}
	}
}
