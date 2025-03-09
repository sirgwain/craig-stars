package cs

type invasion struct {
	planet    *Planet
	defender  *Player
	attacker  *Player
	attackers int
	fleets    []*Fleet // fleets involved in the invasion
}

type invasionResult struct {
	invasion
	defenders          int
	attackersKilled    int
	defendersKilled    int
	remainingAttackers int
	remainingDefenders int
	successful         bool
}

type invader struct {
	invasionsByPlanet map[int][]invasion
}

func newInvader() invader {
	return invader{invasionsByPlanet: make(map[int][]invasion)}
}

// fleetDescription gets the fleet name or an empty string if there were many fleets involved
func (i *invasion) fleetDescription() string {
	if len(i.fleets) == 1 {
		return i.fleets[0].Name
	}

	return ""
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
			existingInvasion.attackers += inv.attackers
			existingInvasion.fleets = append(existingInvasion.fleets, inv.fleets...)
			return
		}
	}

	// we didn't find an existing invasion for this player, add to the other invasions
	i.invasionsByPlanet[inv.planet.Num] = append(i.invasionsByPlanet[inv.planet.Num], inv)
}

// resolveInvasions resolves all invasions for a planet
// TODO: add any logic for N-way invasions
func (i *invader) resolveInvasions(rules *Rules) []invasionResult {
	var results []invasionResult
	for _, invasions := range i.invasionsByPlanet {
		for _, invasion := range invasions {
			results = append(results, invasion.resolve(rules))
		}
	}
	return results
}

// invade a planet with a colonist drop
func (i invasion) resolve(rules *Rules) invasionResult {
	invasionDefenseCoverageFactor := rules.InvasionDefenseCoverageFactor

	// figure out how many attackers are stopped by defenses
	attacker := i.attacker
	defender := i.defender
	attackersAfterDefense := int(float64(i.attackers) * (1 - i.planet.Spec.DefenseCoverage*invasionDefenseCoverageFactor))
	defenders := i.planet.population()

	// determine bonuses for warmongers and inner strength
	attackBonus := attacker.Race.Spec.InvasionAttackBonus
	defenseBonus := defender.Race.Spec.InvasionDefendBonus

	remainingDefenders := 0
	remainingAttackers := 0
	attackersKilled := 0
	defendersKilled := 0
	successful := false

	if float64(attackersAfterDefense)*attackBonus > float64(defenders)*defenseBonus {
		remainingDefenders = 0
		remainingAttackers = roundToNearest100(float64(attackersAfterDefense) - float64(defenders)*defenseBonus/attackBonus)

		// if we have a last-person-standing, they instantly repopulate. :)
		if remainingAttackers == 0 {
			remainingAttackers = 100
		}

		attackersKilled = i.attackers - remainingAttackers
		defendersKilled = defenders
		successful = true
	} else {
		remainingAttackers = 0
		remainingDefenders = roundToNearest100(float64(defenders) - (float64(attackersAfterDefense)*attackBonus)/defenseBonus)

		// if we have a last-person-standing, they instantly repopulate. :)
		if remainingDefenders == 0 {
			remainingDefenders = 100
		}
		attackersKilled = i.attackers
		defendersKilled = defenders - remainingDefenders
	}

	return invasionResult{
		invasion:           i,
		defenders:          defenders,
		attackersKilled:    attackersKilled,
		defendersKilled:    defendersKilled,
		remainingAttackers: remainingAttackers,
		remainingDefenders: remainingDefenders,
		successful:         successful,
	}
}
