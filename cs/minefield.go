package cs

import (
	"fmt"
	"math"
)

// Minefields are zones in space that damage and destroy enemy ships traveling through them.
type MineField struct {
	GameDBObject    `tstype:",extends"`
	MapObject       `tstype:",extends"`
	MineFieldOrders `tstype:",extends"`
	MineFieldType   MineFieldType `json:"mineFieldType"`
	NumMines        int           `json:"numMines"`
	Spec            MineFieldSpec `json:"spec"`
}

type MineFieldOrders struct {
	Detonate bool `json:"detonate,omitempty"`
}

// The type of a mine field.
type MineFieldType string

const (
	MineFieldTypeNone      MineFieldType = "" // TODO: Make standard the zero value
	MineFieldTypeStandard  MineFieldType = "Standard"
	MineFieldTypeHeavy     MineFieldType = "Heavy"
	MineFieldTypeSpeedBump MineFieldType = "SpeedBump"
)

func (t MineFieldType) String() string {
	switch t {
	case MineFieldTypeSpeedBump:
		return "Speed Bump"
	default:
		return string(t)
	}
}

func (t MineFieldType) CanDetonate() bool {
	switch t {
	case MineFieldTypeStandard:
		return true
	default:
		return false
	}
}

type MineFieldSpec struct {
	Radius      float64 `json:"radius"`
	DecayRate   int     `json:"decayRate"`
	CanDetonate bool    `json:"canDetonate"`
}

type MineFieldStats struct {
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

type MineFieldDamage struct {
	Damage         int  `json:"damage,omitempty"`
	ShipsDestroyed int  `json:"shipsDestroyed,omitempty"`
	FleetDestroyed bool `json:"fleetDestroyed,omitempty"`
}

// The radius of a minefield is the sqrt of its mines
func (mf *MineField) Radius() float64 {
	return math.Sqrt(float64(mf.NumMines))
}

func computeMinefieldSpec(rules *Rules, player *Player, mineField *MineField, numPlanets int) MineFieldSpec {
	spec := MineFieldSpec{}
	spec.Radius = mineField.Radius()
	spec.DecayRate = mineField.getDecayRate(rules, player, numPlanets)
	spec.CanDetonate = mineField.MineFieldType.CanDetonate()

	return spec
}

func newMineField(player *Player, mineFieldType MineFieldType, numMines int, num int, position Vector) *MineField {
	return &MineField{
		MapObject: MapObject{
			Type:      MapObjectTypeMineField,
			PlayerNum: player.Num,
			Num:       num,
			Name:      fmt.Sprintf("%s %s Mine Field #%d", player.Race.PluralName, mineFieldType.String(), num),
			Position:  position,
		},
		MineFieldType: mineFieldType,
		NumMines:      numMines,
	}
}

func (mineField *MineField) withOrders(orders MineFieldOrders) *MineField {
	mineField.MineFieldOrders = orders
	return mineField
}

// get the number of mines that will decay this year
// * The base rate for minefield decay is 2% per year.
// * Minefields will decay an additional 4% per planet that is within the field, or 1% per planet for SD races.
// * A detonating SD minefield has an additional 25% decay each year.
// * Normal and Heavy Minefields have a minimum total decay rate of 10 mines per year
// * Speed Bump Minefields have a minimum total decay rate of 2 mines per year
// * There is a maximum total decay rate of 50% per year.
func (mineField *MineField) getDecayRate(rules *Rules, player *Player, numPlanets int) int {
	if !mineField.Owned() {
		// we can't determine decay rate for minefields we don't own
		return -1
	}

	decayRate := player.Race.Spec.MineFieldBaseDecayRate
	decayRate += player.Race.Spec.MineFieldPlanetDecayRate * float64(numPlanets)
	if mineField.Detonate {
		decayRate += player.Race.Spec.MineFieldDetonateDecayRate
	}

	// Space Demolition mines decay slower
	decayFactor := player.Race.Spec.MineFieldMinDecayFactor
	decayRate *= decayFactor
	decayRate = min(decayRate, player.Race.Spec.MineFieldMaxDecayRate)

	// we decay at least 10 mines a year for normal and standard mines
	decayedMines := max(rules.MineFieldStatsByType[mineField.MineFieldType].MinDecay, int(float64(mineField.NumMines)*decayRate+0.5))
	return decayedMines
}

// damage a fleet that hit this minefield
// https://wiki.starsautohost.org/wiki/Guts_of_Minefields
func (mineField *MineField) damageFleet(fleet *Fleet, fleetPlayer *Player, stats MineFieldStats) MineFieldDamage {
	minDamage := stats.MinDamagePerFleet
	damagePerEngine := stats.DamagePerEngine

	// figure out if we have ramscoops, increasing damage if applicable
	for _, token := range fleet.Tokens {
		if token.design.Spec.Engine.FreeSpeed > 1 {
			minDamage = stats.MinDamagePerFleetRS
			damagePerEngine = stats.DamagePerEngineRS
			break
		}
	}

	if minDamage <= 0 && damagePerEngine <= 0 {
		// no minefield damage makes our job very easy
		return MineFieldDamage{}
	}

	totalDamage := 0
	shipsDestroyed := 0
	if fleet.Spec.TotalShips <= 5 {
		// for the first 5 ships, damage is allocated proportionally to engine count
		firstDesignNumEngines := 0
		for i := range fleet.Tokens {
			token := &fleet.Tokens[i]
			if mineField.Detonate && token.design.Spec.ImmuneToOwnDetonation && mineField.OwnedBy(fleetPlayer.Num) {
				// don't damage this fleet if we're remote detonating ourselves
				continue
			}

			design := token.design
			if firstDesignNumEngines != 0 && design.Spec.NumEngines <= firstDesignNumEngines {
				// TODO @sirgwain: please put a good comment here IDK what to say that makes sense for this
				continue
			}

			var tokenDamage int
			if firstDesignNumEngines == 0 {
				firstDesignNumEngines = design.Spec.NumEngines
				tokenDamage = firstDesignNumEngines * minDamage
			} else {
				tokenDamage = damagePerEngine * (design.Spec.NumEngines - firstDesignNumEngines) * token.Quantity
			}

			totalDamage += tokenDamage
			result := token.applyMineDamage(tokenDamage)
			shipsDestroyed += result.shipsDestroyed
		}
	} else {
		for i := range fleet.Tokens {
			token := &fleet.Tokens[i]
			if mineField.Detonate && token.design.Spec.ImmuneToOwnDetonation && mineField.OwnedBy(fleetPlayer.Num) {
				// don't damage this fleet if we're remote detonating ourselves
				continue
			}

			design := token.design
			tokenDamage := damagePerEngine * design.Spec.NumEngines * token.Quantity
			totalDamage += tokenDamage
			result := token.applyMineDamage(tokenDamage)
			shipsDestroyed += result.shipsDestroyed
		}
	}

	return MineFieldDamage{
		Damage:         totalDamage,
		ShipsDestroyed: shipsDestroyed,
		FleetDestroyed: shipsDestroyed == fleet.Spec.TotalShips,
	}
}

// When a minefield is collided with, reduce its number of mines
// based on the number of tokens that hit it.
func (mineField *MineField) reduceMineFieldOnImpact(numTokens int) {
	numMines := mineField.NumMines

	// Successively reduce the field with each collision
	for numReductions := 0; numTokens > 0 && numMines > 0; numTokens -= numReductions {
		switch {
		case numMines <= 200:
			// Reduce by 10 mines per hit, clearing the field if mines drops below 0
			numReductions = min(numTokens, divideRoundAway0(numMines, 10))
			numMines -= 10 * numReductions
		case numMines <= 1000:
			// Apply a 5% reduction per hit until reaching 200
			numReductions = min(numTokens,
				int(math.Ceil(LogBase(0.95, 200/float64(numMines)))))
			numMines = int(float64(numMines) * math.Pow(0.95, float64(numReductions)))
		case numMines <= 5000:
			// Reduce by 50 mines per hit until reaching 1K
			numReductions = min(numTokens, divideRoundAway0(numMines-1000, 10))
			numMines -= 50 * numReductions
		default:
			// Apply a 5% reduction per hit until reaching 5000
			numReductions = min(numTokens,
				int(math.Ceil(LogBase(0.95, 5000/float64(numMines)))))
			numMines = int(float64(numMines) * math.Pow(0.95, float64(numReductions)))

		}
	}

	mineField.NumMines = max(0, numMines)
}

// sweep reduces this MineField's NumMines based on a fleet's inside the minefield
// and minesweep rate.
//
// Returns the number of mines swept, up to the field's current mine count.
func (mineField *MineField) sweep(fleetPosition Vector, mineSweep int, sweepFactor float64) (numSwept int) {
	// we can only sweep up to our position in the minefield, so figure out how far
	// we are from the center and subtract that from the radius to get distance from edge.
	// (Radius is just distance from center to edge)

	radius := mineField.Radius()
	distFromCenter := fleetPosition.DistanceTo(mineField.Position)
	distFromEdge := radius - distFromCenter

	// radius of a minefield is sqrt(numMines) so we can sweep up to dist^2 in mines
	sweepableMines := mineField.NumMines - int(math.Ceil(math.Pow(radius-distFromEdge, 2)))

	old := mineField.NumMines
	mineField.NumMines -= min(sweepableMines, int(float64(mineSweep)*sweepFactor))
	if mineField.NumMines <= 10 {
		// delete fields with under 10 mines in them
		mineField.NumMines = 0
	}

	return old - mineField.NumMines
}

// Check for mine field collisions, damaging and stopping the fleet as applicable.
func checkForMineFieldCollision(rules *Rules, playerGetter playerGetter, mapObjectGetter mapObjectGetter, fleet *Fleet, dest Waypoint, distTraveled float64) (mineField *MineField, actualDist float64) {
	fleetPlayer := playerGetter.getPlayer(fleet.PlayerNum)
	safeWarpBonus := fleetPlayer.Race.Spec.MineFieldSafeWarpBonus

	// see if we are colliding with any of these minefields
	for _, mineField := range mapObjectGetter.getAllMineFields() {
		// we don't hit our own minefields
		if mineField.PlayerNum == fleet.PlayerNum {
			continue
		}

		// our allies don't hit our minefields
		mineFieldPlayer := playerGetter.getPlayer(mineField.PlayerNum)
		if mineFieldPlayer.IsFriend(fleetPlayer.Num) {
			continue
		}

		// we only check if we are going faster than allowed by the minefield
		stats := rules.MineFieldStatsByType[mineField.MineFieldType]
		if dest.WarpSpeed <= stats.MaxSpeed+safeWarpBonus {
			// we are going slow enough to not trigger an explosion
			continue
		}

		// check if we intersect with this minefield
		to := dest.Position.Subtract(fleet.Position).Normalized().Multiply(distTraveled).Add(fleet.Position)
		percentNotInField := segmentIntersectsCircle(fleet.Position, to, mineField.Position, mineField.Spec.Radius)
		if percentNotInField == doesNotIntersect {
			// we aren't colliding with the minefield
			continue
		}

		// check collisions for each light year of travel through the minefield
		lightYearsInField := min(int(mineField.Spec.Radius),
			int(math.Ceil((1-percentNotInField)*distTraveled)))
		lightYearsBeforeField := percentNotInField * distTraveled

		// Each type of minefield has their hit rate multiplied by how many warp speeds
		// the fleet is travelling over the safe limit.
		// Warp 9 in a standard minefield has a 1.5% chance/LY (0.3%*(9-4)

		unsafeWarp := dest.WarpSpeed - (stats.MaxSpeed + safeWarpBonus)
		chanceToHit := stats.ChanceOfHit * float64(unsafeWarp)
		for i := range lightYearsInField {
			if chanceToHit >= rules.random.Float64() {
				// ouch, we hit a mine!
				// We stop moving immediately, so our distance traveled
				fleet.struckMineField = true
				actualDist = lightYearsBeforeField + float64(i)
				return mineField, actualDist
			}
		}

	}

	return mineField, distTraveled
}

// Move this minefield closer to us (in case it's not in our location)
// This was taken from the FreeStars codebase (like many other things)
func (mineField *MineField) moveTowardsMineLayer(position Vector, minesLaid int) {
	totalDist := position.DistanceTo(mineField.Position)

	moveTowardsFactor := min(1, float64(minesLaid)/float64(mineField.NumMines))
	heading := position.Subtract(mineField.Position).Normalized()

	// move the minefield towards the fleet
	mineField.Position = mineField.Position.Add(heading.Normalized().Multiply(totalDist * moveTowardsFactor)).Round()
}
