package ai

import "github.com/sirgwain/craig-stars/cs"

// All difficulties follow the same validity rules. Difficulty controls how
// quickly the AI commits to expansion and war, never deliberate broken orders.
func (ai *aiPlayer) configureStrategy() {
	ai.doctrine = ai.doctrineFor()
	ai.config.colonyPipeline = 4
	ai.config.pace = 1
	switch ai.AIDifficulty {
	case cs.AIDifficultyEasy:
		ai.config.colonyPipeline = 2
		ai.config.pace = .5
		ai.config.colonizerPopulationDensity = .33
		ai.config.colonistTransportDensity = .33
		ai.config.startAttackingYear += 10
	case cs.AIDifficultyHard, cs.AIDifficultyCheater:
		ai.config.colonyPipeline = 8
		ai.config.pace = 1.5
		ai.config.researchReserve = 15
		ai.config.startAttackingYear = max(0, ai.config.startAttackingYear-5)
		ai.config.bomberProductionCutoff = .65
		// Unlock economic and racial capabilities before the common military
		// progression. These goals accumulate in research(), as normal.
		goals := []cs.TechLevel{{Propulsion: 2}, {Construction: 4}}
		switch ai.Race.PRT {
		case cs.AR:
			goals = []cs.TechLevel{{Energy: 3}, {Construction: 2, Electronics: 1}, {Construction: 4}, {Energy: 6}}
		case cs.SS:
			goals = append(goals, cs.TechLevel{Electronics: 4})
		case cs.WM:
			goals = append(goals, cs.TechLevel{Weapons: 6})
		case cs.CA:
			goals = append(goals, cs.TechLevel{Energy: 3, Weapons: 5, Biotechnology: 2})
		case cs.IS:
			goals = append(goals, cs.TechLevel{Energy: 3})
		case cs.SD:
			goals = append(goals, cs.TechLevel{Biotechnology: 3})
		case cs.PP:
			goals = append(goals, cs.TechLevel{Energy: 4})
		case cs.IT:
			goals = append(goals, cs.TechLevel{Propulsion: 5, Construction: 5})
		case cs.JoaT:
			goals = append(goals, cs.TechLevel{Electronics: 3})
		}
		ai.config.researchOrder = append(goals, ai.config.researchOrder...)
	}
	ai.config.startAttackingYear = max(0, ai.config.startAttackingYear+ai.doctrine.attackYears)

	// AR needs energy for income and remote miners for a sustainable mineral
	// supply on every difficulty, even when research is otherwise conservative.
	if ai.Race.PRT == cs.AR && ai.AIDifficulty != cs.AIDifficultyHard && ai.AIDifficulty != cs.AIDifficultyCheater {
		ai.config.researchOrder = append([]cs.TechLevel{{Energy: 3}, {Construction: 2, Electronics: 1}}, ai.config.researchOrder...)
	}
	if ai.Race.PRT == cs.AR {
		// AR income scales with sqrt(population). Spreading population across
		// colonies pays sooner, and a base upgrade can greatly increase capacity
		// without changing the population actually available to seed colonies.
		switch ai.AIDifficulty {
		case cs.AIDifficultyEasy:
			ai.config.colonizerPopulationDensity = .25
		case cs.AIDifficultyHard, cs.AIDifficultyCheater:
			ai.config.colonizerPopulationDensity = .10
		default:
			ai.config.colonizerPopulationDensity = .15
		}
	}
	ai.fleetsByPurpose[cs.FleetPurposeMiner] = fleet{purpose: cs.FleetPurposeMiner, ships: []fleetShip{{purpose: cs.ShipDesignPurposeMiner, quantity: 1}}}
	ai.fleetsByPurpose[cs.FleetPurposeMineLayer] = fleet{purpose: cs.FleetPurposeMineLayer, ships: []fleetShip{{purpose: cs.ShipDesignPurposeDamageMineLayer, quantity: 1}}}
}
