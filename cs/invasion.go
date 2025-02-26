package cs

import "github.com/rs/zerolog"

type invasion struct {
	planet           *Planet
	defender         *Player
	attacker         *Player
	colonistsDropped int
	fleetName        string // empty for multiple fleets
}

type invader struct {
	invasionsByPlanet map[int][]invasion
}

func (i *invader) addInvasion(inv invasion) {
	invasions := i.invasionsByPlanet[inv.planet.Num]
	if len(invasions) == 0 {
		i.invasionsByPlanet[inv.planet.Num] = []invasion{inv}
		return
	}

	for j := range invasions {
		existingInvasion := &invasions[j]
		if existingInvasion.attacker.Num == inv.attacker.Num {
			// add to this existing invasion and remove the fleet name since
			// we are invading with multiple fleets
			existingInvasion.colonistsDropped += inv.colonistsDropped
			existingInvasion.fleetName = ""
			return
		}
	}

	// we didn't find an existing invasion for this player, add to the other invasions
	i.invasionsByPlanet[inv.planet.Num] = append(i.invasionsByPlanet[inv.planet.Num], inv)
}

// resolveInvasions resolves all invasions for a planet
// TODO: add any logic for N-way invasions
func (i *invader) resolveInvasions(log zerolog.Logger, rules *Rules) {
	for _, invasions := range i.invasionsByPlanet {
		for _, invasion := range invasions {
			invadePlanet(log, rules, invasion.planet, invasion.fleetName, invasion.defender, invasion.attacker, invasion.colonistsDropped)
		}
	}
}

// invade a planet with a colonist drop
func invadePlanet(log zerolog.Logger, rules *Rules, planet *Planet, fleetName string, defender *Player, attacker *Player, colonistsDropped int) {
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
		remainingAttackers = roundToNearest100(float64(attackers) - float64(defenders)*defenseBonus/attackBonus)

		// if we have a last-person-standing, they instantly repopulate. :)
		if remainingAttackers == 0 {
			remainingAttackers = 100
		}

		var attackersKilled = colonistsDropped - remainingAttackers

		// notify each player of the invasion
		messager.planetInvaded(defender, planet, fleetName, attacker, defender, attackersKilled, planet.population(), true)
		messager.planetInvaded(attacker, planet, fleetName, attacker, defender, attackersKilled, planet.population(), true)

		// empty this planet
		planet.emptyPlanet()

		// take over the planet.
		planet.PlayerNum = attacker.Num
		planet.setPopulation(remainingAttackers)

		// make sure the defender knows about this new planet
		// the last dying colonist sends a report to their compatriots
		defender.discoverer.discoverPlanet(rules, planet, true)

		// apply a production plan
		if len(attacker.ProductionPlans) > 0 {
			plan := attacker.ProductionPlans[0]
			plan.Apply(planet)
		}

		// check for tech trades
		if !attacker.techLevelGained {
			tt := newTechTrader()
			field := tt.checkInvasionTechTrade(rules, attacker, defender.TechLevels)
			if field != TechFieldNone {
				// sweet, we gained a tech level
				attacker.techLevelGained = true
				attacker.TechLevels.Set(field, attacker.TechLevels.Get(field)+1) // add 1 to corresponding lvl

				messager.playerTechGainedInvasion(attacker, planet, field)
				attacker.updateTechsJustGained(rules.techs, field)

				log.Debug().
					Int("Attacker", attacker.Num).
					Int("Defender", defender.Num).
					Str("Planet", planet.Name).
					Str("field", string(field)).
					Msgf("invader gained tech level")
			}
		}
	} else {
		remainingAttackers = 0
		remainingDefenders = roundToNearest100(float64(defenders) - (float64(attackers)*attackBonus)/defenseBonus)

		// if we have a last-person-standing, they instantly repopulate. :)
		if remainingDefenders == 0 {
			remainingDefenders = 100
		}
		defendersKilled := planet.population() - remainingDefenders

		// notify each player of the invasion
		messager.planetInvaded(defender, planet, fleetName, attacker, defender, colonistsDropped, defendersKilled, false)
		messager.planetInvaded(attacker, planet, fleetName, attacker, defender, colonistsDropped, defendersKilled, false)

		// reduce the population to however many colonists remain
		planet.setPopulation(remainingDefenders)
	}

	log.Debug().
		Int("Defender", defender.Num).
		Int("Attacker", attacker.Num).
		Str("Fleet", fleetName).
		Str("Planet", planet.Name).
		Int("Attackers", attackers).
		Int("Defenders", defenders).
		Int("RemainingAttackers", remainingAttackers).
		Int("RemainingDefenders", remainingDefenders).
		Bool("AttackerWon", planet.PlayerNum == attacker.Num).
		Int("PlanetPlayerNum", planet.PlayerNum).
		Msgf("planet invaded")

}
