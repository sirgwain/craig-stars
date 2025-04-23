package cs

import "math"

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

// fleetDescription gets the name of a fleet involved in an invasion,
// or an empty string if multiple fleets were involved.
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
	// figure out how many attackers are stopped by defenses
	attacker := i.attacker
	defender := i.defender
	attackers := int(float64(i.attackers) * (1 - i.planet.Spec.DefenseCoverage*rules.InvasionDefenseCoverageFactor))
	defenders := i.planet.GetPopulation()

	// determine bonuses for warmongers and inner strength
	attackBonus := attacker.Race.Spec.InvasionAttackBonus
	defenseBonus := defender.Race.Spec.InvasionDefendBonus

	remainingDefenders := 0
	remainingAttackers := 0
	attackersKilled := 0
	defendersKilled := 0
	successful := false

	// TODO: Check how invasion pop loss interacts with partial pop once I have enough sanity
	if float64(attackers)*attackBonus > float64(defenders)*defenseBonus {
		// attackers won

		remainingDefenders = 0
		// if we have a last-person-standing, they instantly repopulate. :)
		remainingAttackers = max(100,
			roundTo100(float64(attackers)-float64(defenders)*defenseBonus/attackBonus, math.Round))

		attackersKilled = i.attackers - remainingAttackers
		defendersKilled = defenders
		successful = true
	} else {
		// defenders won
		remainingAttackers = 0
		// if we have a last-person-standing, they instantly repopulate. :)
		remainingDefenders = max(100,
			roundTo100(float64(defenders)-(float64(attackers)*attackBonus)/defenseBonus, math.Round))
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
