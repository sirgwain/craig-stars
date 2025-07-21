package cs

import (
	"fmt"
	"math"
)

type MinefieldType string

const (
	MinefieldTypeStandard  MinefieldType = "Standard"
	MinefieldTypeHeavy     MinefieldType = "Heavy"
	MinefieldTypeSpeedBump MinefieldType = "SpeedBump"
)

func (t MinefieldType) String() string {
	switch t {
	case MinefieldTypeSpeedBump:
		return "Speed Bump"
	default:
		return string(t)
	}
}

func (t MinefieldType) CanDetonate() bool {
	switch t {
	case MinefieldTypeStandard:
		return true
	default:
		return false
	}
}

type Minefield struct {
	GameDBObject    `tstype:",extends"`
	MapObject       `tstype:",extends"`
	MinefieldOrders `tstype:",extends"`
	MinefieldType   MinefieldType `json:"minefieldType"`
	NumMines        int           `json:"numMines"`
	Spec            MinefieldSpec `json:"spec"`
}

type MinefieldOrders struct {
	Detonate bool `json:"detonate,omitempty"`
}

type MinefieldSpec struct {
	Radius      float64 `json:"radius"`
	DecayRate   int     `json:"decayRate"`
	CanDetonate bool    `json:"canDetonate"`
}

type MinefieldStats struct {
	MinDamagePerFleetRS int     `json:"minDamagePerFleetRS"`
	DamagePerEngineRS   int     `json:"damagePerEngineRS"`
	MaxSpeed            int     `json:"maxSpeed"`
	ChanceOfHit         float64 `json:"chanceOfHit"`
	MinDamagePerFleet   int     `json:"minDamagePerFleet"`
	DamagePerEngine     int     `json:"damagePerEngine"`
	SweepFactor         float64 `json:"sweepFactor"`
	MinDecay            int     `json:"minDecay"`
	CanDetonate         bool    `json:"canDetonate"`
}

type MinefieldDamage struct {
	Damage         int  `json:"damage,omitempty"`
	ShipsDestroyed int  `json:"shipsDestroyed,omitempty"`
	FleetDestroyed bool `json:"fleetDestroyed,omitempty"`
}

// The radius of a minefield is the sqrt of its mines
func (mf *Minefield) Radius() float64 {
	return math.Sqrt(float64(mf.NumMines))
}

func computeMinefieldSpec(rules *Rules, player *Player, minefield *Minefield, numPlanets int) MinefieldSpec {
	spec := MinefieldSpec{}
	spec.Radius = minefield.Radius()
	spec.DecayRate = minefield.getDecayRate(rules, player, numPlanets)
	spec.CanDetonate = minefield.MinefieldType.CanDetonate()

	return spec
}

func newMinefield(player *Player, minefieldType MinefieldType, numMines int, num int, position Vector) *Minefield {
	return &Minefield{
		MapObject: MapObject{
			Type:      MapObjectTypeMinefield,
			PlayerNum: player.Num,
			Num:       num,
			Name:      fmt.Sprintf("%s %s Minefield #%d", player.Race.PluralName, minefieldType.String(), num),
			Position:  position,
		},
		MinefieldType: minefieldType,
		NumMines:      numMines,
	}
}

func (minefield *Minefield) withOrders(orders MinefieldOrders) *Minefield {
	minefield.MinefieldOrders = orders
	return minefield
}

// get the number of mines that will decay this year
// * The base rate for minefield decay is 2% per year.
// * Minefields will decay an additional 4% per planet that is within the field, or 1% per planet for SD races.
// * A detonating SD minefield has an additional 25% decay each year.
// * Normal and Heavy Minefields have a minimum total decay rate of 10 mines per year
// * Speed Bump Minefields have a minimum total decay rate of 2 mines per year
// * There is a maximum total decay rate of 50% per year.
func (minefield *Minefield) getDecayRate(rules *Rules, player *Player, numPlanets int) int {
	if !minefield.Owned() {
		// we can't determine decay rate for minefields we don't own
		return -1
	}

	decayRate := player.Race.Spec.MinefieldBaseDecayRate
	decayRate += player.Race.Spec.MinefieldPlanetDecayRate * float64(numPlanets)
	if minefield.Detonate {
		decayRate += player.Race.Spec.MinefieldDetonateDecayRate
	}

	// Space Demolition mines decay slower
	decayFactor := player.Race.Spec.MinefieldMinDecayFactor
	decayRate *= decayFactor
	decayRate = min(decayRate, player.Race.Spec.MinefieldMaxDecayRate)

	// we decay at least 10 mines a year for normal and standard mines
	decayedMines := max(rules.MinefieldStatsByType[minefield.MinefieldType].MinDecay, int(float64(minefield.NumMines)*decayRate+0.5))
	return decayedMines
}

// damage a fleet that hit this minefield
// https://wiki.starsautohost.org/wiki/Guts_of_Minefields
func (minefield *Minefield) damageFleet(fleet *Fleet, fleetPlayer *Player, stats MinefieldStats) MinefieldDamage {
	hasRamScoop := false
	for _, token := range fleet.Tokens {
		if token.design.Spec.Engine.FreeSpeed > 1 {
			hasRamScoop = true
			break
		}
	}

	minDamage := stats.MinDamagePerFleet
	damagePerEngine := stats.DamagePerEngine
	if hasRamScoop {
		minDamage = stats.MinDamagePerFleetRS
		damagePerEngine = stats.DamagePerEngineRS
	}

	totalDamage := 0
	shipsDestroyed := 0

	if minDamage > 0 {
		if fleet.Spec.TotalShips <= 5 {
			firstDesignNumEngines := 0
			for i := range fleet.Tokens {
				token := &fleet.Tokens[i]
				if minefield.Detonate && token.design.Spec.ImmuneToOwnDetonation && minefield.OwnedBy(fleetPlayer.Num) {
					continue
				}

				design := token.design
				if firstDesignNumEngines == 0 {
					firstDesignNumEngines = design.Spec.NumEngines
					tokenDamage := firstDesignNumEngines * minDamage
					totalDamage += tokenDamage
					result := token.applyMineDamage(tokenDamage)
					shipsDestroyed += result.shipsDestroyed
				} else if design.Spec.NumEngines > firstDesignNumEngines {
					tokenDamage := damagePerEngine * (design.Spec.NumEngines - firstDesignNumEngines) * token.Quantity
					totalDamage += tokenDamage
					result := token.applyMineDamage(tokenDamage)
					shipsDestroyed += result.shipsDestroyed
				}
			}
		} else {
			for i := range fleet.Tokens {
				token := &fleet.Tokens[i]
				if minefield.Detonate && token.design.Spec.ImmuneToOwnDetonation && minefield.OwnedBy(fleetPlayer.Num) {
					continue
				}

				design := token.design
				tokenDamage := damagePerEngine * design.Spec.NumEngines * token.Quantity
				totalDamage += tokenDamage
				result := token.applyMineDamage(tokenDamage)
				shipsDestroyed += result.shipsDestroyed
			}
		}
	}

	return MinefieldDamage{
		Damage:         totalDamage,
		ShipsDestroyed: shipsDestroyed,
		FleetDestroyed: fleet.Spec.TotalShips <= shipsDestroyed,
	}
}

// When a minefield is collided with, reduce its number of mines
func (minefield *Minefield) reduceMinefieldOnImpact() {
	numMines := minefield.NumMines
	if numMines <= 10 {
		numMines = 0
	} else if numMines <= 200 {
		numMines -= 10
	} else if numMines <= 1000 {
		numMines = int(float64(numMines) * 0.95)
	} else if numMines <= 5000 {
		numMines -= 50
	} else {
		numMines = int(float64(numMines) * 0.95)
	}
	minefield.NumMines = numMines
}

func (minefield *Minefield) sweep(rules *Rules, fleetPosition Vector, mineSweep int) int {

	// we can only sweep up to our position in the minefield, so figure out how far we are from the center
	// and subtract that from the radius to determine the edge amount
	//		***
	// 	   *****
	// 	  *******
	// 	 *F**C**** // fleet is 1 from the edge, 3 from the center
	// 	  *******
	// 	   *****
	// 	    ***
	//
	radius := minefield.Radius()
	distFromCenter := fleetPosition.DistanceTo(minefield.Position)
	distFromEdge := minefield.Radius() - distFromCenter

	// radius of a minefield is sqrt(numMines) so we can sweet our dist^2 in mines
	sweepableMines := minefield.NumMines - int(math.Ceil((radius-distFromEdge)*(radius-distFromEdge)))

	old := minefield.NumMines
	minefield.NumMines -= min(sweepableMines, int(float64(mineSweep)*rules.MinefieldStatsByType[minefield.MinefieldType].SweepFactor))
	minefield.NumMines = max(minefield.NumMines, 0)

	numSwept := old - minefield.NumMines
	return numSwept
}

// / Check for minefield collisions. If we collide with one, do damage and stop the fleet
func checkForMinefieldCollision(rules *Rules, playerGetter playerGetter, mapObjectGetter mapObjectGetter, fleet *Fleet, dest Waypoint, distance float64) (minefield *Minefield, distanceTravelled float64) {
	distanceTravelled = distance
	fleetPlayer := playerGetter.getPlayer(fleet.PlayerNum)
	safeWarpBonus := fleetPlayer.Race.Spec.MinefieldSafeWarpBonus

	// see if we are colliding with any of these minefields
	for _, minefield := range mapObjectGetter.getAllMinefields() {
		// we don't hit our own minefields
		if minefield.PlayerNum == fleet.PlayerNum {
			continue
		}

		// our allies don't hit our minefields
		minefieldPlayer := playerGetter.getPlayer(minefield.PlayerNum)
		if minefieldPlayer.IsFriend(fleetPlayer.Num) {
			continue
		}

		// we only check if we are going faster than allowed by the minefield.
		stats := rules.MinefieldStatsByType[minefield.MinefieldType]
		if dest.WarpSpeed > stats.MaxSpeed+safeWarpBonus {
			// this is not our minefield, and we are going fast, check if we intersect.
			from := fleet.Position
			to := (dest.Position.Subtract(fleet.Position).Normalized()).Scale(distance).Add(from)
			collision := segmentIntersectsCircle(from, to, minefield.Position, minefield.Spec.Radius)
			if collision == -1 {
				// miss! phew, that was close!
				continue
			} else {
				// we are travelling through this minefield, for each light year we go through, check for a hit
				// collision is 0 to 1, which is the percent of our travel segment that is NOT in the field.
				// figure out what that is in lightYears
				// if we are travelling 32 light years and 3/4 of it is through the minefield, we need to check
				// for collision 24 times
				lightYearsInField := int(min(float64(minefield.Spec.Radius), math.Ceil(float64((1-collision)*distance))))
				lightYearsBeforeField := collision * distance

				// Each type of minefield has a chance to hit based on how fast
				// the fleet is travelling through the field. A normal mine has a .3% chance
				// of hitting a ship per extra warp over warp 4, so a warp 9 ship
				// has a 1.5% chance of hitting a mine per lightyear travelled
				unsafeWarp := dest.WarpSpeed - (stats.MaxSpeed + safeWarpBonus)
				chanceToHit := stats.ChanceOfHit * float64(unsafeWarp)
				for checkNum := 0; checkNum < lightYearsInField; checkNum++ {
					if chanceToHit >= rules.random.Float64() {
						// ouch, we hit a minefield!
						// we stop moving at the hit, so if we made it 8 checks out of 24 for our above example
						// we only travel 8 lightyears through the field (plus whatever distance we travelled to get to the field)
						fleet.struckMinefield = true
						distanceTravelled = lightYearsBeforeField + float64(checkNum)
						return minefield, distanceTravelled
					}
				}
			}
		}

	}

	return minefield, distance
}

// Move this minefield closer to us (in case it's not in our location)
// This was taken from the FreeStars codebase (like many other things)
func (minefield *Minefield) moveTowardsMineLayer(position Vector, minesLaid int) {
	totalDist := position.DistanceTo(minefield.Position)

	moveTowardsFactor := min(1, float64(minesLaid)/float64(minefield.NumMines))
	heading := position.Subtract(minefield.Position).Normalized()

	// move the minefield towards the fleet
	minefield.Position = minefield.Position.Add(heading.Normalized().Scale(totalDist * moveTowardsFactor)).Round()
}
