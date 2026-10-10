package cs

import (
	"log/slog"
	"math"
	"sort"
)

// Bombers orbiting enemy planets will Bomb planets, killing population and destroying installations.
type Bomb struct {
	Quantity             int     `json:"quantity,omitempty"`
	KillRate             float64 `json:"killRate,omitempty"`
	MinKillRate          int     `json:"minKillRate,omitempty"`
	StructureDestroyRate float64 `json:"structureDestroyRate,omitempty"`
	UnterraformRate      int     `json:"unterraformRate,omitempty"`
}

type BombingResult struct {
	BomberName         string `json:"bomberName,omitempty"`
	NumBombers         int    `json:"numBombers,omitempty"`
	ColonistsKilled    int    `json:"colonistsKilled,omitempty"`
	MinesDestroyed     int    `json:"minesDestroyed,omitempty"`
	FactoriesDestroyed int    `json:"factoriesDestroyed,omitempty"`
	DefensesDestroyed  int    `json:"defensesDestroyed,omitempty"`
	UnterraformAmount  Hab    `json:"unterraformAmount,omitzero"`
	PlanetEmptied      bool   `json:"planetEmptied,omitempty"`
	fleet              *Fleet
}

type bomber struct {
	rules *Rules
	log   *slog.Logger
}

func newBomber(log *slog.Logger, rules *Rules) bomber {
	return bomber{rules: rules, log: log}
}

// bombPlanet bombs this planet with any enemy bombers orbiting it.
//
// Each player bombs in turn, in player number order so results are deterministic.
// A player's normal, smart and retro bombs all strike together (see bombPlayer),
// and each later player bombs what is left after the previous players.
// Bombing stops once the colony is dead.
func (b *bomber) bombPlanet(planet *Planet, planetOwner *Player, enemyBombers []*Fleet, pg playerGetter) {
	// get the players with bombers here, sorted by player number
	players := make(map[int]bool)
	for _, fleet := range enemyBombers {
		players[fleet.PlayerNum] = true
	}
	nums := make([]int, 0, len(players))
	for num := range players {
		nums = append(nums, num)
	}
	sort.Ints(nums)

	// each player bombs the planet, and both sides get a message
	for _, num := range nums {
		if planet.Cargo.Colonists == 0 {
			break
		}
		fleets := b.getBombersForPlayer(enemyBombers, num)
		result := b.bombPlayer(planet, planetOwner, fleets)
		if result.NumBombers == 0 {
			continue
		}
		messager.fleetBombedPlanet(pg.getPlayer(num), result.fleet, planet, result)
		messager.planetBombed(planetOwner, planet, result.fleet, result)
	}
	// the bombing killed everyone, the planet is now empty
	if planet.Cargo.Colonists == 0 {
		planet.emptyPlanet()
		messager.planetDiedOff(planetOwner, planet)
	}
}

// get a slice of all bombers for a player
func (b *bomber) getBombersForPlayer(fleets []*Fleet, playerNum int) []*Fleet {
	result := []*Fleet{}
	for _, fleet := range fleets {
		if fleet.PlayerNum == playerNum {
			result = append(result, fleet)
		}
	}
	return result
}

// bombPlayer evaluates all of one player's bombs against the defenses at the
// start of their strike. Smart casualties precede normal casualties; the minimum is
// a floor on their combined total.
//
// Rates are in thousandths and population is in kT (100s of colonists) until the
// final result.
//
// The flow is:
//  1. collect all the player's normal, smart and retro bombs
//  2. compute the kill and destroy rates, reduced by the planet's defenses
//  3. destroy installations in proportion to how many of each the planet has
//  4. kill colonists with smart bombs, then normal bombs, then apply the minimum kill
//  5. unterraform the planet with retro bombs
//  6. update the planet
func (b *bomber) bombPlayer(planet *Planet, defender *Player, fleets []*Fleet) BombingResult {
	// 1. collect bombs from every fleet with bombs. The first one names the attack
	result := BombingResult{}
	normal, smart, retro := []Bomb{}, []Bomb{}, 0
	smartSurvival := 1.0
	for _, fleet := range fleets {
		if len(fleet.Spec.Bombs)+len(fleet.Spec.SmartBombs)+len(fleet.Spec.RetroBombs) == 0 {
			continue
		}
		if result.fleet == nil {
			result.fleet = fleet
			result.BomberName = fleet.Name
		}
		result.NumBombers++
		normal = append(normal, fleet.Spec.Bombs...)
		smart = append(smart, fleet.Spec.SmartBombs...)
		for _, bomb := range fleet.Spec.RetroBombs {
			retro += bomb.UnterraformRate * bomb.Quantity
		}
	}
	if result.NumBombers == 0 {
		return result
	}
	// 2. compute rates. Normal bombs add their kill rates (and minimum kills), while
	// smart bombs stack multiplicatively, each only killing what the previous
	// missed. Defenses reduce the kill rates, but only cover half of structure damage
	survival := 1 - planet.Spec.DefenseCoverage
	normalRate, floor, structureRate := 0.0, 0, 0.0
	for _, bomb := range normal {
		normalRate += bomb.KillRate * 10 * float64(bomb.Quantity)
		floor += bomb.MinKillRate * bomb.Quantity / 100
		structureRate += bomb.StructureDestroyRate * float64(bomb.Quantity)
	}
	for _, bomb := range smart {
		smartSurvival *= math.Pow(1-bomb.KillRate/100, float64(bomb.Quantity))
	}
	smartRate := int(math.Round((1 - smartSurvival) * 1000))
	smartRate = int(math.Round(float64(smartRate) * (1 - planet.Spec.DefenseCoverageSmart)))
	rate := int(math.Round(normalRate * survival))
	floor = int(math.Round(float64(floor) * survival))
	structures := int(math.Round(structureRate * (1 - planet.Spec.DefenseCoverage/2)))
	// 3. destroy installations. Factories and defenses each lose their share of the
	// destroyed structures, with the fractional part rounded up at random. Mines
	// take the rest, so the total destroyed matches the bombs' strength
	total := planet.Mines + planet.Factories + planet.Defenses
	proportional := func(count int) int {
		product := count * structures
		killed := product / total
		if remainder := product % total; remainder > 0 && b.rules.random.Intn(total) < remainder {
			killed++
		}
		return min(count, killed)
	}
	if total > 0 && structures > 0 {
		result.FactoriesDestroyed = proportional(planet.Factories)
		result.DefensesDestroyed = proportional(planet.Defenses)
		result.MinesDestroyed = min(planet.Mines, max(0, structures-result.FactoriesDestroyed-result.DefensesDestroyed))
	}
	// 4. kill colonists. Smart bombs go first and can't kill the last 100 colonists.
	// Normal bombs kill a share of the survivors (rounding up at random), always
	// killing at least 100. The combined total is never less than the minimum kill
	population := planet.Cargo.Colonists
	if population > 0 {
		smartKilled := min(population-1, population*smartRate/1000)
		product := (population - smartKilled) * rate
		killed := product / 1000
		if remainder := product % 1000; remainder > 0 && b.rules.random.Intn(1000) <= remainder {
			killed++
		}
		killed += smartKilled
		if rate > 0 {
			killed = max(1, killed)
		}
		result.ColonistsKilled = min(population, max(killed, floor)) * 100
	}
	planet.addPopulation(-result.ColonistsKilled)
	planet.Mines -= result.MinesDestroyed
	planet.Factories -= result.FactoriesDestroyed
	planet.Defenses -= result.DefensesDestroyed
	// 5. retro bombs undo terraforming, moving each habitat axis back toward the
	// planet's base habitat. Defenses cover half of their strength. This
	// subtracts a truncated defensive reduction, rather than rounding the final
	// retro strength
	retro -= int(float64(retro) * planet.Spec.DefenseCoverage / 2)
	result.UnterraformAmount = b.getUnterraformAmount(min(retro, 500), planet.BaseHab, planet.Hab)
	planet.Hab = planet.Hab.Add(result.UnterraformAmount)
	planet.TerraformedAmount = planet.Hab.Subtract(planet.BaseHab)

	// 6. update the planet so the next player bombs against the reduced defenses
	result.PlanetEmptied = planet.Cargo.Colonists == 0
	planet.Spec = ComputePlanetSpec(b.rules, defender, planet)
	planet.MarkDirty()

	b.log.Debug("fleet bombed planet",
		slog.Int("Player", result.fleet.PlayerNum),
		slog.String("Planet", planet.Name),
		slog.String("Fleet", result.BomberName),
		slog.Int("NumFleets", result.NumBombers),
		slog.Int("PlanetPlayer", planet.PlayerNum),
		slog.Int("Killed", result.ColonistsKilled),
		slog.Int("MinesDestroyed", result.MinesDestroyed),
		slog.Int("FactoriesDestroyed", result.FactoriesDestroyed),
		slog.Int("DefensesDestroyed", result.DefensesDestroyed),
		slog.String("UnterraformAmount", result.UnterraformAmount.String()))
	return result
}

// getUnterraformAmount gets the amount we should unterraform with retro bombs. Each
// axis moves back toward its base habitat by up to retroBombAmount.
func (b *bomber) getUnterraformAmount(retroBombAmount int, baseHab, hab Hab) Hab {
	amount := Hab{}
	for _, axis := range HabTypes {
		diff := baseHab.Get(axis) - hab.Get(axis)
		amount.Set(axis, max(-retroBombAmount, min(retroBombAmount, diff)))
	}
	return amount
}
