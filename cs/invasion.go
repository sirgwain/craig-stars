package cs

import "github.com/rs/zerolog/log"

// invade a planet with a colonist drop
func invadePlanet(planet *Planet, fleet *Fleet, defender *Player, attacker *Player, colonistsDropped int, invasionDefenseCoverageFactor float64) {

	// figure out how many attackers are stopped by defenses
	attackers := int(float64(colonistsDropped) * (1 - planet.Spec.DefenseCoverage*invasionDefenseCoverageFactor))
	defenders := planet.population()

	// determine bonuses for warmongers and inner strength
	attackBonus := attacker.Race.Spec.InvasionAttackBonus
	defenseBonus := defender.Race.Spec.InvasionDefendBonus

	if float64(attackers)*attackBonus > float64(defenders)*defenseBonus {
		remainingAttackers := roundToNearest100f(float64(attackers) - float64(defenders)*defenseBonus/attackBonus)

		// if we have a last-person-standing, they instantly repopulate. :)
		if remainingAttackers == 0 {
			remainingAttackers = 100
		}

		var attackersKilled = colonistsDropped - remainingAttackers

		// notify each player of the invasion
		messager.planetInvaded(defender, planet, fleet, defender.Race.PluralName, attacker.Race.PluralName, attackersKilled, planet.population(), true)
		messager.planetInvaded(attacker, planet, fleet, defender.Race.PluralName, attacker.Race.PluralName, attackersKilled, planet.population(), true)

		// take over the planet.
		// empty this planet
		planet.PlayerNum = attacker.Num
		planet.starbase = nil
		planet.Scanner = false
		planet.Defenses = 0 // defenses are destroyed during invasion
		planet.ProductionQueue = []ProductionQueueItem{}
		planet.setPopulation(remainingAttackers)

		// apply a production plan
		if len(attacker.ProductionPlans) > 0 {
			plan := attacker.ProductionPlans[0]
			plan.Apply(planet)
		}

	} else {
		var remainingDefenders = roundToNearest100f(float64(defenders) - (float64(attackers)*attackBonus)/defenseBonus)

		// if we have a last-person-standing, they instantly repopulate. :)
		if remainingDefenders == 0 {
			remainingDefenders = 100
		}
		defendersKilled := planet.population() - remainingDefenders

		// notify each player of the invasion
		messager.planetInvaded(defender, planet, fleet, defender.Race.PluralName, attacker.Race.PluralName, colonistsDropped, defendersKilled, false)
		messager.planetInvaded(attacker, planet, fleet, defender.Race.PluralName, attacker.Race.PluralName, colonistsDropped, defendersKilled, false)

		// reduce the population to however many colonists remain
		planet.setPopulation(remainingDefenders)
	}

	log.Debug().
		Int64("GameID", planet.GameID).
		Int("Defender", defender.Num).
		Int("Attacker", attacker.Num).
		Str("Fleet", fleet.Name).
		Str("Planet", planet.Name).
		Int("Attackers", attackers).
		Int("Defenders", defenders).
		Int("PlanetPlayerNum", planet.PlayerNum).
		Msgf("planet invaded")

}
