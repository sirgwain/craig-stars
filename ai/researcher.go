package ai

import (
	"github.com/sirgwain/craig-stars/cs"
)

func (ai *aiPlayer) research() {

	var previousReseearchOrder cs.TechLevel
	ai.NextResearchField = cs.NextResearchFieldSameField

	for _, order := range ai.config.researchOrder {
		order := order.Max(previousReseearchOrder)
		previousReseearchOrder = order

		if !ai.TechLevels.HasRequiredLevels(order) {
			// find the lowest level to research next
			ai.Researching = ai.TechLevels.LowestMissingLevel(order)
			break
		}
	}

	// if we didn't pick one, just let the below logic figure it out
	if ai.Researching == cs.TechFieldNone {
		ai.Researching = cs.Energy
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
