package cs

import (
	"cmp"
	"fmt"
	"math"
	"slices"
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
	GameDBObject
	MapObject
	MinefieldOrders
	MinefieldType MinefieldType `json:"minefieldType"`
	NumMines      int           `json:"numMines"`
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
	Damage         int    `json:"damage,omitempty"`
	ShipsDestroyed int    `json:"shipsDestroyed,omitempty"`
	FleetDestroyed bool   `json:"fleetDestroyed,omitempty"`
	Position       Vector `json:"position,omitempty"` // where the fleet was hit
}

func (dmg MinefieldDamage) noDamage() bool {
	return dmg.Damage == 0 && dmg.ShipsDestroyed == 0 && !dmg.FleetDestroyed
}

// The radius of a minefield is the sqrt of its mines
func (mf *Minefield) Radius() float64 {
	return math.Sqrt(float64(mf.NumMines))
}

func ComputeMinefieldSpec(rules *Rules, player *Player, minefield *Minefield, numPlanets int) MinefieldSpec {
	spec := MinefieldSpec{}
	spec.Radius = minefield.Radius()
	spec.DecayRate = minefield.getDecayRate(rules, player, numPlanets)
	spec.CanDetonate = player.Race.Spec.CanDetonateMinefields && minefield.MinefieldType.CanDetonate()

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

// damageFleet damages a fleet that hit a minefield. Each ship takes
// damage per engine, more for ships with ramscoops, and small fleets take at least the field's minimum damage.
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
	if damagePerEngine == 0 {
		// speed bumps stop fleets without damaging them
		return MinefieldDamage{}
	}

	// fleets of fewer than 5 ships take extra damage to make up the minimum, on their first stack
	extraDamage := 0
	if fleet.Spec.TotalShips < 5 {
		extraDamage = max(0, minDamage-damagePerEngine*fleet.Spec.TotalShips)
	}

	for i := range fleet.Tokens {
		token := &fleet.Tokens[i]
		if minefield.Detonate && token.design.Spec.ImmuneToOwnDetonation && minefield.OwnedBy(fleetPlayer.Num) {
			continue
		}

		tokenDamage := (damagePerEngine*token.Quantity + extraDamage) * token.design.Spec.NumEngines
		extraDamage = 0
		totalDamage += tokenDamage
		result := token.applyMineDamage(tokenDamage)
		shipsDestroyed += result.shipsDestroyed
	}

	if totalDamage > 0 {
		fleet.noHeal = true
	}
	return MinefieldDamage{
		Damage:         totalDamage,
		ShipsDestroyed: shipsDestroyed,
		FleetDestroyed: fleet.Spec.TotalShips <= shipsDestroyed,
	}
}

// reduceMinefieldOnImpact removes mines from a field a fleet hit: 5% of the field, at least 10
// mines, and for fields over 1000 mines, 50 mines or 1% of the field, whichever is more
func (minefield *Minefield) reduceMinefieldOnImpact() {
	removed := minefield.NumMines / 20
	if removed > 50 {
		removed = max(50, minefield.NumMines/100)
	} else if removed < 10 {
		removed = 10
	}
	minefield.NumMines = max(0, minefield.NumMines-removed)
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

// minefieldStretch is a stretch of a fleet's path through enemy minefields of one type
type minefieldStretch struct {
	minefieldType MinefieldType
	start, end    float64 // ly along the path
	minefields    []*Minefield
}

// checkForMinefieldCollision checks a fleet moving distance ly toward dest for hitting enemy minefields.
// It returns the minefield hit and how far the fleet got before it hit, or nil and distance if it made it through.
//
// The fleet's speed for the check comes from how far it actually moves, not its warp, so a fleet arriving or
// making a short pursuit step is checked at a lower speed. Overlapping fields of the same type are one stretch
// of the path, checked once per ly in path order, and the last ly before arriving at dest isn't checked.
func checkForMinefieldCollision(rules *Rules, playerGetter playerGetter, mapObjectGetter mapObjectGetter, fleet *Fleet, dest Waypoint, distance float64) (minefield *Minefield, distanceTravelled float64) {
	totalDist := fleet.Position.DistanceTo(dest.Position)
	if distance <= 0 || totalDist == 0 {
		return nil, distance
	}

	// don't check the last ly before we arrive
	checkDist := min(distance, math.Ceil(totalDist)-1)

	// the speed we're going is the lowest warp from 3 to 10 that covers the distance
	warp := 3
	for warp < 10 && float64(warp*warp) < checkDist-1 {
		warp++
	}

	fleetPlayer := playerGetter.getPlayer(fleet.PlayerNum)
	safeWarpBonus := fleetPlayer.Race.Spec.MinefieldSafeWarpBonus
	if warp <= 3+safeWarpBonus {
		return nil, distance
	}

	from := fleet.Position.ToFloat64()
	heading := dest.Position.ToFloat64().Subtract(from).Normalized()
	stretches := map[MinefieldType][]minefieldStretch{}
	for _, mf := range mapObjectGetter.getAllMinefields() {
		// we don't hit our own minefields, or our allies'
		if mf.PlayerNum == fleet.PlayerNum || playerGetter.getPlayer(mf.PlayerNum).IsFriend(fleetPlayer.Num) {
			continue
		}

		start, end, ok := pathThroughCircle(from, heading, checkDist, mf.Position.ToFloat64(), mf.Radius())
		if !ok {
			continue
		}
		stretches[mf.MinefieldType] = addMinefieldStretch(stretches[mf.MinefieldType], minefieldStretch{mf.MinefieldType, start, end, []*Minefield{mf}})
	}

	// check each stretch in order along the path, for all field types
	ordered := []minefieldStretch{}
	for _, typeStretches := range stretches {
		ordered = append(ordered, typeStretches...)
	}
	slices.SortStableFunc(ordered, func(a, b minefieldStretch) int {
		return cmp.Or(cmp.Compare(a.start, b.start), cmp.Compare(a.minefieldType, b.minefieldType))
	})

	for _, stretch := range ordered {
		// fields only hit fleets going faster than they allow
		stats := rules.MinefieldStatsByType[stretch.minefieldType]
		unsafeWarp := warp - (stats.MaxSpeed + safeWarpBonus)
		if unsafeWarp <= 0 {
			continue
		}

		// Each type of minefield has a chance to hit based on how fast the fleet is travelling through
		// the field. A normal mine has a .3% chance of hitting a ship per warp over warp 4 per ly, so a
		// warp 9 ship has a 1.5% chance of hitting a mine per ly travelled
		chanceToHit := stats.ChanceOfHit * float64(unsafeWarp)
		lightYears := int(math.Ceil(stretch.end - stretch.start))
		for ly := range lightYears {
			if chanceToHit >= rules.random.Float64() {
				// ouch, we hit a minefield! we stop where we hit it
				fleet.struckMinefield = true
				distanceTravelled = stretch.start + float64(ly)
				return closestMinefield(stretch.minefields, from.Add(heading.Scale(distanceTravelled))), distanceTravelled
			}
		}
	}

	return nil, distance
}

// pathThroughCircle returns where a path from a point along heading for length ly enters and leaves a circle,
// in ly along the path, or false if it misses
func pathThroughCircle(from, heading VectorFloat64, length float64, center VectorFloat64, radius float64) (start, end float64, ok bool) {
	toCenter := center.Subtract(from)
	alongPath := toCenter.Dot(heading)
	distToPathSquared := toCenter.Dot(toCenter) - alongPath*alongPath
	if distToPathSquared >= radius*radius {
		return 0, 0, false
	}

	halfChord := math.Sqrt(radius*radius - distToPathSquared)
	start = max(0, alongPath-halfChord)
	end = min(length, alongPath+halfChord)
	return start, end, start < end
}

// addMinefieldStretch adds a stretch to stretches of the same field type, merging it with any it overlaps
// so the overlap is only checked once
func addMinefieldStretch(stretches []minefieldStretch, stretch minefieldStretch) []minefieldStretch {
	merged := []minefieldStretch{}
	for _, other := range stretches {
		if other.end < stretch.start || stretch.end < other.start {
			merged = append(merged, other)
			continue
		}
		stretch.start = min(stretch.start, other.start)
		stretch.end = max(stretch.end, other.end)
		stretch.minefields = append(stretch.minefields, other.minefields...)
	}
	return append(merged, stretch)
}

// closestMinefield returns the minefield whose center is closest to where a fleet hit
func closestMinefield(minefields []*Minefield, position VectorFloat64) *Minefield {
	return slices.MinFunc(minefields, func(a, b *Minefield) int {
		return cmp.Compare(a.Position.ToFloat64().DistanceTo(position), b.Position.ToFloat64().DistanceTo(position))
	})
}

// Move this minefield closer to us (in case it's not in our location)
// This was taken from the FreeStars codebase (like many other things)
func (minefield *Minefield) moveTowardsMineLayer(position Vector, minesLaid int) {
	totalDist := position.DistanceTo(minefield.Position)

	moveTowardsFactor := min(1, float64(minesLaid)/float64(minefield.NumMines))
	heading := position.Subtract(minefield.Position).Normalized()

	// move the minefield towards the fleet
	minefield.Position = minefield.Position.ToFloat64().Add(heading.Normalized().Scale(totalDist * moveTowardsFactor)).ToInt(true)
}
