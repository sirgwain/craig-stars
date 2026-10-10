package cs

import "sort"

type invasion struct {
	planet       *Planet
	defender     *Player
	attacker     *Player
	attackers    int
	fleets       []*Fleet // fleets involved in the invasion
	colonization bool
}

type invasionResult struct {
	invasion
	defenders                 int
	attackersKilled           int
	attackersKilledByDefenses int
	defendersKilled           int
	remainingAttackers        int
	remainingDefenders        int
	successful                bool
	uninhabited               bool
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
			existingInvasion.colonization = existingInvasion.colonization || inv.colonization
			existingInvasion.fleets = append(existingInvasion.fleets, inv.fleets...)
			return
		}
	}

	// we didn't find an existing invasion for this player, add to the other invasions
	i.invasionsByPlanet[inv.planet.Num] = append(i.invasionsByPlanet[inv.planet.Num], inv)
}

// resolveInvasions resolves all drops at each planet together, in kT, returning
// one result per attacker, grouped by planet.
func (i *invader) resolveInvasions(rules *Rules) []invasionResult {
	var results []invasionResult
	planets := make([]int, 0, len(i.invasionsByPlanet))
	for num := range i.invasionsByPlanet {
		planets = append(planets, num)
	}
	sort.Ints(planets)
	for _, num := range planets {
		results = append(results, resolvePlanetInvasions(rules, i.invasionsByPlanet[num])...)
	}
	return results
}

// resolvePlanetInvasions resolves every attacker dropping colonists on a planet in the same phase.
// Attackers combine their strength against the defenders. If they win, the
// strongest attacker takes the planet, losing colonists to the defenders and
// to the next strongest rival. Attackers tied for strongest destroy each other.
func resolvePlanetInvasions(rules *Rules, invasions []invasion) []invasionResult {
	if len(invasions) == 0 {
		return nil
	}
	sort.Slice(invasions, func(a, b int) bool { return invasions[a].attacker.Num < invasions[b].attacker.Num })
	planet, defender := invasions[0].planet, invasions[0].defender
	defenders := planet.Cargo.Colonists
	defensePower := 0
	if defender != nil {
		defensePower = int(float64(defenders) * defender.Race.Spec.InvasionDefendBonus)
	}

	// work in kT, like the colonists on the planet
	survival := 1 - planet.Spec.DefenseCoverage*rules.InvasionDefenseCoverageFactor
	total, strongest, second, winner := 0, 0, 0, -1
	tied := false
	for n, inv := range invasions {
		power := 0
		if !inv.attacker.Race.Spec.LivesOnStarbases {
			power = int(float64(inv.attackers/100) * inv.attacker.Race.Spec.InvasionAttackBonus * survival)
		}
		total += power
		switch {
		case winner == -1 || power > strongest:
			second, strongest, winner, tied = strongest, power, n, false
		case power == strongest:
			second, tied = power, true
		case power > second:
			second = power
		}
	}

	remainingDefenders, remainingAttackers := 0, 0
	defenderWins := defensePower > total
	if defenderWins {
		remainingDefenders = defenders - defenders*total/defensePower
	} else if !tied {
		remainingAttackers = invasions[winner].attackers / 100
		if total > 0 {
			remainingAttackers = remainingAttackers * (total - defensePower) / total
		}
		if strongest > 0 {
			remainingAttackers = remainingAttackers * (strongest - second) / strongest
		}
		// the last colonist standing repopulates
		remainingAttackers = max(1, remainingAttackers)
	}

	results := make([]invasionResult, len(invasions))
	for n, inv := range invasions {
		remaining := 0
		success := !defenderWins && !tied && n == winner
		if success {
			remaining = remainingAttackers * 100
		}
		results[n] = invasionResult{
			invasion:                  inv,
			defenders:                 defenders * 100,
			attackersKilled:           inv.attackers - remaining,
			attackersKilledByDefenses: inv.attackers - int(float64(inv.attackers)*survival),
			defendersKilled:           (defenders - remainingDefenders) * 100,
			remainingAttackers:        remaining,
			remainingDefenders:        remainingDefenders * 100,
			successful:                success,
			uninhabited:               !defenderWins && tied,
		}
	}
	return results
}
