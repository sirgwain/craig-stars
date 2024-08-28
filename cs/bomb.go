package cs

import (
	"math"

	"github.com/rs/zerolog/log"
)

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
	UnterraformAmount  Hab    `json:"unterraformAmount,omitempty"`
	PlanetEmptied      bool   `json:"planetEmptied,omitempty"`
	fleet              *Fleet
}

type bomb struct {
	rules *Rules
}

// Bombers orbiting enemy planets will Bomb planets
// ============================================================================
// Algorithms:
// Normalpopkills = sum[bomb_kill_perc(n)*#(n)] * (1-Def(pop))
// Minkills = sum[bomb_kill_min(n)*#(n)] * (1-Def(pop))
//
// 10 Cherry and 5 M-70 bombing vs 100 Neutron Defs (97.92%)
//
// The calculations are, population kill:
//
// a    0.025 * 10  0.25        10 Cherry bombs
// b    0.012 * 5   0.06        5 M-70 bombs
// c    a + b       0.31        Total kill factor
// d    1 - 0.97    0.0208      1 - defense factor for 100 neutron defences
// e    c * d           0.006448    Total kill factor
// f    pop * c         64.48       Total colonists killed
//
// Minimum kill:
//
// a 10*300 + 5*300  4500
// b 1 - 0.97        0.0208   1 - defense factor for 100 neutron defences
// c a *b            156      Total minimum kill
// ============================================================================
type bomber interface {
	// Attempt to bomb this planet
	bombPlanet(planet *Planet, planetOwner *Player, enemyBombers []*Fleet, pg playerGetter)
}

func NewBomber(rules *Rules) bomber {
	return &bomb{rules: rules}
}

// add two bombing results and return the total
func (result BombingResult) Add(r BombingResult) BombingResult {
	// keep the new one, if the source is empty
	if result.NumBombers == 0 {
		return r
	}

	return BombingResult{
		NumBombers:         MaxInt(result.NumBombers, r.NumBombers), // we only care about the highest number of bombers for the result
		ColonistsKilled:    result.ColonistsKilled + r.ColonistsKilled,
		MinesDestroyed:     result.MinesDestroyed + r.MinesDestroyed,
		FactoriesDestroyed: result.FactoriesDestroyed + r.FactoriesDestroyed,
		DefensesDestroyed:  result.DefensesDestroyed + r.DefensesDestroyed,
		UnterraformAmount:  result.UnterraformAmount.Add(r.UnterraformAmount),
		PlanetEmptied:      result.PlanetEmptied || r.PlanetEmptied,
		fleet:              result.fleet, // we only care about the first fleet
	}
}

// bomb this planet if there are any bombers orbiting it
func (b *bomb) bombPlanet(planet *Planet, planetOwner *Player, enemyBombers []*Fleet, pg playerGetter) {
	// get a list of all players orbiting the planet
	orbitingPlayerNums := map[int]bool{}
	resultsByPlayer := map[int]BombingResult{}
	for _, fleet := range enemyBombers {
		orbitingPlayerNums[fleet.PlayerNum] = true
		resultsByPlayer[fleet.PlayerNum] = BombingResult{}
	}

	// bomb the planet with regular bombs
	for playerNum := range orbitingPlayerNums {
		result := b.normalBombPlanet(planet, planetOwner, pg.getPlayer(playerNum), b.getBombersForPlayer(enemyBombers, playerNum))

		resultsByPlayer[playerNum] = resultsByPlayer[playerNum].Add(result)
		// stop bombing if everyone is dead
		if planet.population() == 0 {
			break
		}
	}

	// bomb the planet with smart bombs
	if planet.population() > 0 {
		for playerNum := range orbitingPlayerNums {
			result := b.smartBombPlanet(planet, planetOwner, pg.getPlayer(playerNum), b.getBombersForPlayer(enemyBombers, playerNum))
			resultsByPlayer[playerNum] = resultsByPlayer[playerNum].Add(result)

			// stop bombing if everyone is dead
			if planet.population() == 0 {
				break
			}
		}
	}

	// deterraform planets
	if planet.population() > 0 && planet.BaseHab != planet.Hab {
		for playerNum := range orbitingPlayerNums {
			result := b.retroBombPlanet(planet, planetOwner, pg.getPlayer(playerNum), b.getBombersForPlayer(enemyBombers, playerNum))
			resultsByPlayer[playerNum] = resultsByPlayer[playerNum].Add(result)
		}
	}

	for playerNum, result := range resultsByPlayer {
		if result.NumBombers == 0 {
			continue
		}

		attacker := pg.getPlayer(playerNum)

		// let each player know a bombing happened
		messager.fleetBombedPlanet(attacker, result.fleet, planet, result)
		messager.planetBombed(planetOwner, planet, result.fleet, result)
	}

	// if, after bombing, the planet is all out of pop, empty it
	if planet.population() == 0 {
		planet.emptyPlanet()
		messager.planetDiedOff(planetOwner, planet)
	}
}

// get a slice of all bombers for a player
func (b *bomb) getBombersForPlayer(fleets []*Fleet, playerNum int) []*Fleet {
	result := []*Fleet{}
	for _, fleet := range fleets {
		if fleet.PlayerNum == playerNum {
			result = append(result, fleet)
		}
	}
	return result
}

// bomb this planet with a slice of fleets
func (b *bomb) normalBombPlanet(planet *Planet, defender *Player, attacker *Player, bombers []*Fleet) BombingResult {

	// do all normal bombs
	bombs := []Bomb{}
	fleets := []*Fleet{}
	for _, fleet := range bombers {
		if len(fleet.Spec.Bombs) > 0 {
			bombs = append(bombs, fleet.Spec.Bombs...)
			fleets = append(fleets, fleet)
		}
	}

	if len(bombs) == 0 {
		return BombingResult{}
	}

	// figure out the killRate and minKill for this fleet's bombs
	defenseCoverage := planet.Spec.DefenseCoverage
	killRateColonistsKilled := roundToNearest100f(b.getColonistsKilledForBombs(planet.population(), defenseCoverage, bombs))
	minColonistsKilled := roundToNearest100(b.getMinColonistsKilledForBombs(planet.population(), defenseCoverage, bombs))

	killed := MaxInt(killRateColonistsKilled, minColonistsKilled)
	leftoverPopulation := MaxInt(0, planet.population()-killed)
	actualKilled := planet.population() - leftoverPopulation
	planet.setPopulation(leftoverPopulation)

	// apply this against mines/factories and defenses proportionally
	structuresDestroyed := b.getStructuresDestroyed(defenseCoverage, bombs)
	totalStructures := planet.Mines + planet.Factories + planet.Defenses
	leftoverMines := 0
	leftoverFactories := 0
	leftoverDefenses := 0
	if totalStructures > 0 {
		leftoverMines = MaxInt(0, int(float64(planet.Mines)-float64(structuresDestroyed)*float64(planet.Mines)/float64(totalStructures)))
		leftoverFactories = MaxInt(0, int(float64(planet.Factories)-float64(structuresDestroyed)*float64(planet.Factories)/float64(totalStructures)))
		leftoverDefenses = MaxInt(0, int(float64(planet.Defenses)-float64(structuresDestroyed)*float64(planet.Defenses)/float64(totalStructures)))
	}

	// make sure we only count stuctures that were actually destroyed
	minesDestroyed := planet.Mines - leftoverMines
	factoriesDestroyed := planet.Factories - leftoverFactories
	defensesDestroyed := planet.Defenses - leftoverDefenses

	planet.Mines = leftoverMines
	planet.Factories = leftoverFactories
	planet.Defenses = leftoverDefenses

	// update planet spec
	planet.Spec = computePlanetSpec(b.rules, defender, planet)

	log.Debug().
		Int64("GameID", planet.GameID).
		Int("Player", attacker.Num).
		Str("Planet", planet.Name).
		Str("Fleet", fleets[0].Name).
		Int("NumFleets", len(fleets)).
		Int("PlanetPlayer", planet.PlayerNum).
		Int("ActualKilled", actualKilled).
		Int("MinesDestroyed", minesDestroyed).
		Int("FactoriesDestroyed", factoriesDestroyed).
		Int("DefensesDestroyed", defensesDestroyed).
		Msgf("fleet bombed planet")

	return BombingResult{
		BomberName:         fleets[0].Name,
		NumBombers:         len(fleets),
		ColonistsKilled:    actualKilled,
		MinesDestroyed:     minesDestroyed,
		FactoriesDestroyed: factoriesDestroyed,
		DefensesDestroyed:  defensesDestroyed,
		PlanetEmptied:      leftoverPopulation == 0,
		fleet:              fleets[0],
	}
}

// smartbomb the planet for each fleet
func (b *bomb) smartBombPlanet(planet *Planet, defender *Player, attacker *Player, bombers []*Fleet) BombingResult {
	smartDefenseCoverage := planet.Spec.DefenseCoverageSmart

	// get all smart bombs from these fleets
	bombs := []Bomb{}
	fleets := []*Fleet{}
	for _, fleet := range bombers {
		if len(fleet.Spec.SmartBombs) > 0 {
			bombs = append(bombs, fleet.Spec.SmartBombs...)
			fleets = append(fleets, fleet)
		}
	}

	if len(bombs) == 0 {
		return BombingResult{}
	}

	// figure out the killRate and minKill for this fleet's bombs
	smartKilled := roundToNearest100f(b.getColonistsKilledWithSmartBombs(planet.population(), smartDefenseCoverage, bombs))

	leftoverPopulation := MaxInt(0, planet.population()-smartKilled)
	actualKilled := planet.population() - leftoverPopulation
	planet.setPopulation(leftoverPopulation)

	// update planet spec
	planet.Spec = computePlanetSpec(b.rules, defender, planet)

	log.Debug().
		Int64("GameID", planet.GameID).
		Int("Player", attacker.Num).
		Str("Planet", planet.Name).
		Str("Fleet", fleets[0].Name).
		Int("NumFleets", len(fleets)).
		Int("PlanetPlayer", planet.PlayerNum).
		Int("ActualKilled", actualKilled).
		Msgf("fleet smart bombed planet")

	return BombingResult{
		BomberName:      fleets[0].Name,
		NumBombers:      len(fleets),
		ColonistsKilled: actualKilled,
		PlanetEmptied:   leftoverPopulation == 0,
		fleet:           fleets[0],
	}
}

// retroBombPlanet a planet for each fleet
func (b *bomb) retroBombPlanet(planet *Planet, defender *Player, attacker *Player, bombers []*Fleet) BombingResult {
	// do all retro bombs
	bombs := []Bomb{}
	fleets := []*Fleet{}
	for _, fleet := range bombers {
		if len(fleet.Spec.RetroBombs) > 0 {
			bombs = append(bombs, fleet.Spec.RetroBombs...)
			fleets = append(fleets, fleet)
		}
	}

	if len(bombs) == 0 {
		return BombingResult{}
	}

	// sum up all the unterraforming
	var retroBombAmount int
	for _, bomb := range bombs {
		retroBombAmount += bomb.UnterraformRate * bomb.Quantity
	}
	unterraformAmount := b.getUnterraformAmount(retroBombAmount, planet.BaseHab, planet.Hab)

	if unterraformAmount.absSum() == 0 {
		return BombingResult{}
	}

	// apply the unterraform amount
	planet.Hab = planet.Hab.Add(unterraformAmount)
	planet.TerraformedAmount = planet.TerraformedAmount.Add(unterraformAmount)

	// update planet spec
	planet.Spec = computePlanetSpec(b.rules, defender, planet)

	log.Debug().
		Int64("GameID", planet.GameID).
		Int("Player", attacker.Num).
		Str("Planet", planet.Name).
		Int("PlanetPlayer", planet.PlayerNum).
		Str("Fleet", fleets[0].Name).
		Int("NumFleets", len(fleets)).
		Str("UnterraformAmount", unterraformAmount.String()).
		Msgf("fleet retro bombed planet")

	return BombingResult{
		BomberName:        fleets[0].Name,
		NumBombers:        len(fleets),
		UnterraformAmount: unterraformAmount,
		fleet:             fleets[0],
	}
}

// Get the amount we should unterraform with retro bombs
func (b *bomb) getUnterraformAmount(retroBombAmount int, baseHab, hab Hab) Hab {
	unterraformAmount := Hab{}
	for i := 0; i < retroBombAmount; i++ {
		// find the current diff based on the unterraforming we've done so far
		habDiff := hab.Subtract(baseHab).Add(unterraformAmount)
		if habDiff.absSum() > 0 {
			largestTerraformHab := Grav
			largestTerraformAmount := 0
			for _, habType := range HabTypes {
				if AbsInt(habDiff.Get(habType)) > AbsInt(largestTerraformAmount) {
					largestTerraformAmount = habDiff.Get(habType)
					largestTerraformHab = habType
				}
			}

			// apply an unterraform amount in whatever direction we are going, to the largest terraform hab
			direction := 1
			if largestTerraformAmount > 0 {
				direction = -1
			}
			unterraformAmount.Set(largestTerraformHab, unterraformAmount.Get(largestTerraformHab)+direction)
		}
	}

	return unterraformAmount
}

// Get colonists killed using the KillRate of a bomb
func (b *bomb) getColonistsKilledForBombs(population int, defenseCoverage float64, bombs []Bomb) float64 {
	// calculate the killRate for all these bombs
	var killRate float64 = 0
	for _, bomb := range bombs {
		killRate += bomb.KillRate * float64(bomb.Quantity)
	}

	return killRate / 100.0 * (1 - defenseCoverage) * float64(population)
}

// Get minimum colonists killed using the MinKillRate of a bomb
func (b *bomb) getMinColonistsKilledForBombs(population int, defenseCoverage float64, bombs []Bomb) int {
	// calculate the minKill for all these bombs
	minKill := 0
	for _, bomb := range bombs {
		minKill += bomb.MinKillRate * bomb.Quantity
	}

	return int(float64(minKill) * (1 - defenseCoverage))
}

// Normal bombs versus buildings.
//
//	Destroy_Build = sum[destroy_build_type(n)*#(n)] * (1-Def(build))
//
// e.g. 10 Cherry + 5 M70 vs 100 Neutron Defs
//
//	= sum[10*10; 5*6] * (1-(97.92%/2))
//	= sum[100; 30] * (1-(48.96%))
//	= 130 * (1- 0.4896)
//	= 130 * 0.5104
//	= ~66 Buildings will be destroyed.
//
// Building kills are allotted proportionately to each building type on
// the planet.  For example, a planet with 1000 installations (of all
// three types combined) taking 400 building kills will lose 40% of each
// of its factories, mines, and defenses.  If there had been 350 mines,
// 550 factories, and 100 defenses, the losses would be 140 mines, 220
// factories, and 40 defenses.
//
// Normal bombs versus buildings.
//
//	Destroy_Build = sum[destroy_build_type(n)*#(n)] * (1-Def(build))
//
// e.g. 10 Cherry + 5 M70 vs 100 Neutron Defs
//
//	= sum[10*10; 5*6] * (1-(97.92%/2))
//	= sum[100; 30] * (1-(48.96%))
//	= 130 * (1- 0.4896)
//	= 130 * 0.5104
//	= ~66 Buildings will be destroyed.
//
// Building kills are allotted proportionately to each building type on
// the planet.  For example, a planet with 1000 installations (of all
// three types combined) taking 400 building kills will lose 40% of each
// of its factories, mines, and defenses.  If there had been 350 mines,
// 550 factories, and 100 defenses, the losses would be 140 mines, 220
// factories, and 40 defenses.

// Calculates the structures destroyed using the StructureDestroyRate of bombs
func (b *bomb) getStructuresDestroyed(defenseCoverage float64, bombs []Bomb) int {
	// calculate the StructureDestroyRate for all these bombs
	var structuresDestroyed float64 = 0
	for _, bomb := range bombs {
		structuresDestroyed += bomb.StructureDestroyRate * float64(bomb.Quantity)
	}

	// this will destroy some number of structures that are allocated proportionally
	// among mines, factories and defenses
	// NOTE: defense coverage is halved for structures
	return int(structuresDestroyed * (1 - defenseCoverage*0.5))
}

// Get the number of colonists killed by smart bombs
// ============================================================================
// Each smart bomb type has a specific pop-kill percentage.  The values
// given by _ONE_ bomb are summarized here:
//
// Smart              1.3%
// Neutron            2.2%
// Enriched Neutron   3.5%
// Peerless           5.0%
// Annihilator        7.0%
//
// Smart bombs do *not* add linearly; instead, they use this formula:
//
//	Pop_kill(smart) = (1-Def(smart))(1 - multiply[ (1 - kill_perc(n)^#n) ])
//
// Where "multiply[x(n)]" is the math "big-pi" operator, which means
// multiply all the terms together, i.e.:
//
//	multiply[x(n)] = x(n1)*x(n2)*x(n3)... *x(ni)
//
// e.g. 10 Annihilators + 5 neutron vs. 100 Neutron-Defs(Def(smart)=85.24%)
//
//	= (1-85.24%) * (1 -  multiply[((1-7%)^10); ((1-2.2%)^5)])
//	= (1-0.8524) * (1 -  ((1-0.07)^10) * ((1-0.022)^5))
//	= 0.1476 * (1 - (0.93^10) * (0.978^5))
//	= 0.1476 * (1 - 0.484 * 0.895)
//	= 0.1476 * 0.56682
//	= 0.0837
//	= 8.37% of planetary pop will be killed.
//
// ============================================================================

// Get number of colonists killed via smart bombs
func (b *bomb) getColonistsKilledWithSmartBombs(population int, defenseCoverageSmart float64, bombs []Bomb) float64 {
	smartKillRate := 0.0
	for _, bomb := range bombs {
		if smartKillRate == 0 {
			smartKillRate = math.Pow(1-bomb.KillRate/100.0, float64(bomb.Quantity))
		} else {
			smartKillRate *= math.Pow(1-bomb.KillRate/100.0, float64(bomb.Quantity))
		}
	}

	if smartKillRate != 0 {
		percentKilled := (1 - defenseCoverageSmart) * (1 - smartKillRate)
		return float64(population) * percentKilled
	}
	return 0
}
