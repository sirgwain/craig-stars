package cs

import "github.com/rs/zerolog/log"

// invade a planet with a colonist drop
func invadePlanet(rules *Rules, planet *Planet, fleet *Fleet, defender *Player, attacker *Player, colonistsDropped int) {
	invasionDefenseCoverageFactor := rules.InvasionDefenseCoverageFactor

	// figure out how many attackers are stopped by defenses
	attackers := int(float64(colonistsDropped) * (1 - planet.Spec.DefenseCoverage*invasionDefenseCoverageFactor))
	defenders := planet.population()

	// determine bonuses for warmongers and inner strength
	attackBonus := attacker.Race.Spec.InvasionAttackBonus
	defenseBonus := defender.Race.Spec.InvasionDefendBonus

	remainingAttackers := 0
	remainingDefenders := 0

	if float64(attackers)*attackBonus > float64(defenders)*defenseBonus {
		remainingDefenders = 0
		remainingAttackers = roundToNearest100f(float64(attackers) - float64(defenders)*defenseBonus/attackBonus)

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
		planet.Starbase = nil
		planet.Scanner = false
		planet.Defenses = 0 // defenses are destroyed during invasion
		planet.ProductionQueue = []ProductionQueueItem{}
		planet.setPopulation(remainingAttackers)

		// make sure the defender knows about this new planet
		// the last dying colonist sends a report to their compatriots
		discover := newDiscoverer(defender)
		discover.discoverPlanet(rules, defender, planet, true)

		// apply a production plan
		if len(attacker.ProductionPlans) > 0 {
			plan := attacker.ProductionPlans[0]
			plan.Apply(planet)
		}

		// check for tech trade
		if !attacker.techLevelGained {
			techTrader := newTechTrader()
			field := techTrader.techLevelGained(rules, attacker.TechLevels, defender.TechLevels)
			if field != TechFieldNone {
				// sweet, we gained a tech level
				attacker.techLevelGained = true
				attacker.TechLevels.Set(field, attacker.TechLevels.Get(field)+1)
				attacker.Messages = append(attacker.Messages, newPlanetMessage(PlayerMessageTechLevelGainedInvasion, planet).
					withSpec(PlayerMessageSpec{Field: field}))
				log.Debug().
					Int64("GameID", planet.GameID).
					Int("Attacker", attacker.Num).
					Int("Defender", defender.Num).
					Str("Planet", planet.Name).
					Str("field", string(field)).
					Msgf("invader gained tech level")

			}
		}
	} else {
		remainingAttackers = 0
		remainingDefenders = roundToNearest100f(float64(defenders) - (float64(attackers)*attackBonus)/defenseBonus)

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
		Int("RemainingAttackers", remainingAttackers).
		Int("RemainingDefenders", remainingDefenders).
		Bool("AttackerWon", planet.PlayerNum == attacker.Num).
		Int("PlanetPlayerNum", planet.PlayerNum).
		Msgf("planet invaded")

}
