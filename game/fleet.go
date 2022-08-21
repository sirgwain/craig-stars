package game

import (
	"fmt"
	"math"
	"time"
)

// warpfactor for using a stargate vs moving with warp drive
const StargateWarpFactor = 11

// time period to perform a task, like patrol
const Indefinite = -1

// use automatic warp factor for patrols
const PatrolWarpFactorAutomatic = -1

// target fleets in any range when patrolling
const PatrolRangeInfinite = -1

type Fleet struct {
	MapObject
	PlanetID         uint        `json:"-"` // for starbase fleets that are owned by a planet
	BaseName         string      `json:"baseName"`
	Cargo            Cargo       `json:"cargo,omitempty" gorm:"embedded;embeddedPrefix:cargo_"`
	Fuel             int         `json:"fuel"`
	Damage           int         `json:"damage"`
	BattlePlanID     uint        `json:"battlePlan"`
	Tokens           []ShipToken `json:"tokens" gorm:"constraint:OnDelete:CASCADE;"`
	Waypoints        []Waypoint  `json:"waypoints" gorm:"serializer:json"`
	RepeatOrders     bool        `json:"repeatOrders,omitempty"`
	Heading          Vector      `json:"heading,omitempty" gorm:"embedded;embeddedPrefix:heading_"`
	WarpSpeed        int         `json:"warpSpeed,omitempty"`
	PreviousPosition *Vector     `json:"previousPosition,omitempty" gorm:"embedded;embeddedPrefix:previous_position_"`
	Orbiting         bool        `json:"orbiting,omitempty"`
	Starbase         bool        `json:"starbase,omitempty"`
	Spec             *FleetSpec  `json:"spec" gorm:"serializer:json"`
}

type FleetSpec struct {
	ShipDesignSpec
	Purposes         map[ShipDesignPurpose]bool `json:"purposes"`
	TotalShips       int                        `json:"totalShips"`
	MassEmpty        int                        `json:"massEmpty"`
	BasePacketSpeed  int                        `json:"basePacketSpeed"`
	SafePacketSpeed  int                        `json:"safePacketSpeed"`
	BaseCloakedCargo int                        `json:"baseCloakedCargo"`
	HasMassDriver    bool                       `json:"hasMassDriver,omitempty"`
	HasStargate      bool                       `json:"hasStargate,omitempty"`
	Stargate         string                     `json:"stargate,omitempty"`
}

type ShipToken struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	FleetID   uint        `json:"gameId"`
	Design    *ShipDesign `json:"-" gorm:"foreignKey:DesignID"`
	DesignID  uint        `json:"designId"`
	Quantity  int         `json:"quantity"`
}

type Waypoint struct {
	FleetID           uint          `json:"-"`
	TargetID          uint          `json:"targetId,omitempty"`
	Position          Vector        `json:"position,omitempty" gorm:"embedded"`
	WarpFactor        int           `json:"warpFactor,omitempty"`
	WaitAtWaypoint    bool          `json:"waitAtWaypoint,omitempty"`
	TargetType        MapObjectType `json:"targetType,omitempty"`
	TargetNum         *int          `json:"targetNum,omitempty"`
	TargetPlayerNum   *int          `json:"targetPlayerNum,omitempty"`
	TransferToPlayer  *int          `json:"transferToPlayer,omitempty"`
	TargetName        string        `json:"targetName,omitempty"`
	PartiallyComplete bool          `json:"partiallyComplete,omitempty"`
}

// create a new fleet
func NewFleet(player *Player, design *ShipDesign, num int, name string, waypoints []Waypoint) Fleet {
	return Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			GameID:    player.GameID,
			PlayerID:  player.ID,
			PlayerNum: player.Num,
			Dirty:     true,
			Num:       num,
			Name:      fmt.Sprintf("%s #%d", name, num),
			Position:  waypoints[0].Position,
		},
		BaseName: name,
		Tokens: []ShipToken{
			{Design: design, Quantity: 1},
		},
		Waypoints: waypoints,
	}
}

// create a new fleet that is a starbase
func NewStarbase(player *Player, planet *Planet, design *ShipDesign, name string) Fleet {
	fleet := NewFleet(player, design, 0, name, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, 1)})
	fleet.Starbase = true

	return fleet
}

func (f *Fleet) String() string {
	return fmt.Sprintf("Fleet %s #%d", f.BaseName, f.Num)
}

func (f *Fleet) WithCargo(cargo Cargo) *Fleet {
	f.Cargo = cargo
	return f
}

func (f *Fleet) WithPosition(position Vector) *Fleet {
	f.Position = position
	// todo: should we set waypoints in a builder?
	f.Waypoints = []Waypoint{{Position: position}}
	return f
}

func NewPlanetWaypoint(position Vector, num int, name string, warpFactor int) Waypoint {
	return Waypoint{
		Position:   position,
		TargetNum:  &num,
		TargetName: name,
		WarpFactor: warpFactor,
	}
}

func ComputeFleetSpec(rules *Rules, player *Player, fleet *Fleet) *FleetSpec {
	spec := FleetSpec{
		ShipDesignSpec: ShipDesignSpec{
			ScanRange:    NoScanner,
			ScanRangePen: NoScanner,
			SpaceDock:    UnlimitedSpaceDock,
		},
		Purposes: map[ShipDesignPurpose]bool{},
	}
	spec.Mass = fleet.Cargo.Total()

	for _, token := range fleet.Tokens {

		// update our total ship count
		spec.TotalShips += token.Quantity

		if token.Design.Purpose != ShipDesignPurposeNone {
			spec.Purposes[token.Design.Purpose] = true
		}

		// use the lowest ideal speed for this fleet
		// if we have multiple engines
		if token.Design.Spec.Engine != "" {
			engine := rules.Techs.GetEngine(token.Design.Spec.Engine)
			if spec.IdealSpeed == 0 {
				spec.IdealSpeed = engine.IdealSpeed
			} else {
				spec.IdealSpeed = MinInt(spec.IdealSpeed, engine.IdealSpeed)
			}
		}
		// cost
		spec.Cost = token.Design.Spec.Cost.MultiplyInt(token.Quantity)

		// mass
		spec.Mass += token.Design.Spec.Mass * token.Quantity
		spec.MassEmpty += token.Design.Spec.Mass * token.Quantity

		// armor
		spec.Armor += token.Design.Spec.Armor * token.Quantity

		// shield
		spec.Shield += token.Design.Spec.Shield * token.Quantity

		// cargo
		spec.CargoCapacity += token.Design.Spec.CargoCapacity * token.Quantity

		// fuel
		spec.FuelCapacity += token.Design.Spec.FuelCapacity * token.Quantity

		// minesweep
		spec.MineSweep += token.Design.Spec.MineSweep * token.Quantity

		// remote mining
		spec.MiningRate += token.Design.Spec.MiningRate * token.Quantity

		// remote terraforming
		spec.TerraformRate += token.Design.Spec.TerraformRate * token.Quantity

		// colonization
		spec.Colonizer = spec.Colonizer || token.Design.Spec.Colonizer
		spec.OrbitalConstructionModule = spec.OrbitalConstructionModule || token.Design.Spec.OrbitalConstructionModule

		// spec all mine layers in the fleet
		if token.Design.Spec.CanLayMines {
			for key := range token.Design.Spec.MineLayingRateByMineType {
				if _, ok := spec.MineLayingRateByMineType[key]; ok {
					spec.MineLayingRateByMineType[key] = 0
				}
				spec.MineLayingRateByMineType[key] += token.Design.Spec.MineLayingRateByMineType[key] * token.Quantity
			}
		}

		// We should only have one ship stack with spacdock capabilities, but for this logic just go with the max
		spec.SpaceDock = MaxInt(spec.SpaceDock, token.Design.Spec.SpaceDock)

		// sadly, the fleet only gets the best repair bonus from one design
		spec.RepairBonus = math.Max(spec.RepairBonus, token.Design.Spec.RepairBonus)

		spec.ScanRange = MaxInt(spec.ScanRange, token.Design.Spec.ScanRange)
		spec.ScanRangePen = MaxInt(spec.ScanRangePen, token.Design.Spec.ScanRangePen)
		if token.Design.Spec.Scanner {
			spec.Scanner = true
		}

		// add bombs
		if token.Design.Spec.Bomber {
			spec.Bomber = true
			spec.Bombs = append(spec.Bombs, token.Design.Spec.Bombs...)
			spec.SmartBombs = append(spec.SmartBombs, token.Design.Spec.SmartBombs...)
			spec.RetroBombs = append(spec.RetroBombs, token.Design.Spec.RetroBombs...)
		}

		// check if any tokens have weapons
		// we process weapon slots per stack, so we don't need to spec all
		// weapons in a fleet
		if token.Design.Spec.HasWeapons {
			spec.HasWeapons = true
		}

		if token.Design.Spec.CloakUnits > 0 {
			// calculate the cloak units for this token based on the design's cloak units (i.e. 70 cloak units / kT for a stealh cloak)
			spec.CloakUnits += token.Design.Spec.CloakUnits
		} else {
			// if this ship doesn't have cloaking, it counts as cargo (except for races with free cargo cloaking)
			if !player.Race.Spec.FreeCargoCloaking {
				spec.BaseCloakedCargo += token.Design.Spec.Mass * token.Quantity
			}
		}

		// choose the best tachyon detector ship
		spec.ReduceCloaking = math.Max(spec.ReduceCloaking, token.Design.Spec.ReduceCloaking)

		spec.CanStealFleetCargo = spec.CanStealFleetCargo || token.Design.Spec.CanStealFleetCargo
		spec.CanStealPlanetCargo = spec.CanStealPlanetCargo || token.Design.Spec.CanStealPlanetCargo
	}

	// compute the cloaking based on the cloak units and cargo
	spec.CloakPercent = fleet.computeFleetCloakPercent(&spec, player.Race.Spec.FreeCargoCloaking)

	return &spec
}

// compute a fleet's cloak percent based on its current cargo/mass/cloak units
func (f *Fleet) computeFleetCloakPercent(spec *FleetSpec, freeCargoCloaking bool) int {
	cloakUnits := spec.CloakUnits

	// starbases have no mass or cargo, but fleet cloaking is adjusted for it
	if spec.Mass > 0 {
		// figure out how much cargo we are cloaking
		cloakedCargo := 0
		if !freeCargoCloaking {
			cloakedCargo = f.Cargo.Total()
		}

		cloakUnits = int(math.Round(float64(cloakUnits) * float64(spec.MassEmpty) / float64(spec.MassEmpty+cloakedCargo)))
	}
	return getCloakPercentForCloakUnits(cloakUnits)
}

func (f *Fleet) AvailableCargoSpace() int {
	return clamp(f.Spec.CargoCapacity-f.Cargo.Total(), 0, f.Spec.CargoCapacity)
}

// transfer cargo from a planet to/from a fleet
func (f *Fleet) TransferPlanetCargo(planet *Planet, transferAmount Cargo) error {

	if f.AvailableCargoSpace() < transferAmount.Total() {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", f.Name, f.AvailableCargoSpace(), transferAmount.Total(), planet.Name)
	}

	if !planet.Cargo.CanTransfer(transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %v from %s, there is not enough to transfer", f.Name, transferAmount, planet.Name)
	}

	// transfer the cargo
	f.Cargo = f.Cargo.Add(transferAmount)
	planet.Cargo = planet.Cargo.Subtract(transferAmount)

	return nil
}

// transfer cargo from a fleet to/from a fleet
func (f *Fleet) TransferFleetCargo(fleet *Fleet, transferAmount Cargo) error {

	if f.AvailableCargoSpace() < transferAmount.Total() {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", f.Name, f.AvailableCargoSpace(), transferAmount.Total(), fleet.Name)
	}

	if !fleet.Cargo.CanTransfer(transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %v from %s, there is not enough to transfer", f.Name, transferAmount, fleet.Name)
	}

	// transfer the cargo
	f.Cargo = f.Cargo.Add(transferAmount)
	fleet.Cargo = fleet.Cargo.Subtract(transferAmount)

	return nil
}

func (fleet *Fleet) moveFleet(game *Game, rules *Rules, player *Player, wp0, wp1 Waypoint, totalDist float64) {
	fleet.PreviousPosition = &Vector{fleet.Position.X, fleet.Position.Y}
	dist := float64(wp1.WarpFactor * wp1.WarpFactor)
	// round up, if we are <1 away, i.e. the target is 81.9 ly away, warp 9 (81 ly travel) should be able to make it there
	if dist < totalDist && totalDist-dist < 1 {
		dist = math.Ceil(totalDist)
	}

	// make sure we end up at a whole number
	vectorTravelled := wp1.Position.Subtract(fleet.Position).Normalized().Scale(dist)
	dist = vectorTravelled.Length()
	// don't overshoot
	dist = math.Min(totalDist, dist)

	// check for CE engine failure
	if player.Race.Spec.EngineFailureRate > 0 && wp1.WarpFactor > player.Race.Spec.EngineReliableSpeed && player.Race.Spec.EngineFailureRate >= rules.Random.Float64() {
		messager.fleetEngineFailure(player, fleet)
		return
	}

	// get the cost for the fleet
	fuelCost := fleet.GetFuelCost(rules.Techs, player, wp1.WarpFactor, dist)
	var fuelGenerated int = 0
	if fuelCost > fleet.Fuel {
		// we will run out of fuel
		// if this distance would have cost us 10 fuel but we have 6 left, only travel 60% of the distance.
		distanceFactor := float64(fleet.Fuel) / float64(fuelCost)
		dist = dist * distanceFactor

		// collide with minefields on route, but don't hit a minefield if we run out of fuel beforehand
		// TODO: add back in minefield check
		// dist = CheckForMineFields(fleet, player, wp1, dist)

		fleet.Fuel = 0
		wp1.WarpFactor = fleet.getNoFuelWarpFactor(rules.Techs, player)
		messager.fleetOutOfFuel(player, fleet, wp1.WarpFactor)

		// if we ran out of fuel 60% of the way to our normal distance, the remaining 40% of our time
		// was spent travelling at fuel generation speeds:
		remainingDistanceTravelled := (1 - distanceFactor) * float64(wp1.WarpFactor*wp1.WarpFactor)
		dist += remainingDistanceTravelled
		fuelGenerated = fleet.getFuelGeneration(rules.Techs, player, wp1.WarpFactor, remainingDistanceTravelled)
	} else {
		// collide with minefields on route, but don't hit a minefield if we run out of fuel beforehand
		// TODO: add back in minefield check
		// actualDist := CheckForMineFields(fleet, player, wp1, dist)
		actualDist := dist
		if actualDist != dist {
			dist = actualDist
			fuelCost = fleet.GetFuelCost(rules.Techs, player, wp1.WarpFactor, dist)
			// we hit a minefield, update fuel usage
		}

		fleet.Fuel -= fuelCost
		fuelGenerated = fleet.getFuelGeneration(rules.Techs, player, wp1.WarpFactor, dist)
	}

	// message the player about fuel generation
	fuelGenerated = MinInt(fuelGenerated, fleet.Spec.FuelCapacity-fleet.Fuel)
	if fuelGenerated > 0 {
		fleet.Fuel += fuelGenerated
		messager.fleetGeneratedFuel(player, fleet, fuelGenerated)
	}

	// assuming we move at all, make sure we are no longer orbiting any planets
	if dist > 0 && fleet.Orbiting {
		fleet.Orbiting = false
	}

	// TODO: repeat orders, can we just append wp0 to waypoints when we repeat?
	// if wp0.OriginalTarget == nil || !wp0.OriginalPosition.HasValue {
	// 	wp0.OriginalTarget = wp0.Target
	// 	wp0.OriginalPosition = fleet.Position
	// }

	if totalDist == dist {
		fleet.completeMove(game, wp0, wp1)
	} else {
		// move this fleet closer to the next waypoint
		fleet.WarpSpeed = wp1.WarpFactor
		fleet.Heading = (wp1.Position.Subtract(fleet.Position)).Normalized()
		wp0.TargetNum = nil
		wp0.TargetName = ""
		wp0.PartiallyComplete = true

		fleet.Position = fleet.Position.Add(fleet.Heading.Scale(dist))
		fleet.Position = fleet.Position.Round()
		wp0.Position = fleet.Position
		fleet.Waypoints[0] = wp0
	}
}

func (fleet *Fleet) gateFleet(rules *Rules, player *Player, wp0, wp1 Waypoint, totalDist float64) {
	panic("unimplemented")
}

// Engine fuel usage calculation courtesy of m.a@stars
func (fleet *Fleet) getFuelCostForEngine(warpFactor int, mass int, dist float64, ifeFactor float64, engine *TechEngine) int {
	if warpFactor == 0 {
		return 0
	}
	// 1 mg of fuel will move 200kT of weight 1 LY at a Fuel Usage Number of 100.
	// Number of engines doesn't matter. Neither number of ships with the same engine.

	distanceCeiling := math.Ceil(dist) // rounding to next integer gives best graph fit
	// window.status = 'Actual distance used is ' + Distan + 'ly';

	// IFE is applied to drive specifications, just as the helpfile hints.
	// Stars! probably does it outside here once per turn per engine to save time.
	engineEfficiency := math.Ceil(ifeFactor * float64(engine.FuelUsage[warpFactor]))

	// 20000 = 200*100
	// Safe bet is Stars! does all this with integer math tricks.
	// Subtracting 2000 in a loop would be a way to also get the rounding.
	// Or even bitshift for the 2 and adjust "decimal point" for the 1000
	teorFuel := (math.Floor(float64(mass)*engineEfficiency*distanceCeiling/2000) / 10)
	// using only one decimal introduces another artifact: .0999 gets rounded down to .0

	// The heavier ships will benefit the most from the accuracy
	intFuel := int(math.Ceil(teorFuel))

	// That's all. Nothing really fancy, much less random. Subtle differences in
	// math lib workings might explain the rarer and smaller discrepancies observed
	return intFuel
	// Unrelated to this fuel math are some quirks inside the
	// "negative fuel" watchdog when the remainder of the
	// trip is < 1 ly. Aahh, the joys of rounding! ;o)
}

// Get the Fuel cost for this fleet to travel a certain distance at a certain speed
func (fleet *Fleet) GetFuelCost(techStore *TechStore, player *Player, warpFactor int, distance float64) int {

	// figure out how much fuel we're going to use
	efficiencyFactor := 1 - player.Race.Spec.FuelEfficiencyOffset

	var fuelCost int = 0

	// compute each ship stack separately
	for _, token := range fleet.Tokens {
		// figure out this ship stack's mass as well as it's proportion of the cargo
		engine := techStore.GetEngine(token.Design.Spec.Engine)
		mass := token.Design.Spec.Mass * token.Quantity
		fleetCargo := fleet.Cargo.Total()
		stackCapacity := token.Design.Spec.CargoCapacity * token.Quantity
		fleetCapacity := fleet.Spec.CargoCapacity

		if fleetCapacity > 0 {
			mass += int(float64(fleetCargo) * (float64(stackCapacity) / float64(fleetCapacity)))
		}

		fuelCost += fleet.getFuelCostForEngine(warpFactor, mass, distance, efficiencyFactor, engine)
	}

	return fuelCost
}

// Get the warp factor for when we run out of fuel.
func (fleet *Fleet) getNoFuelWarpFactor(techStore *TechStore, player *Player) int {
	// find the lowest freeSpeed from all the fleet's engines
	var freeSpeed = math.MaxInt
	for _, token := range fleet.Tokens {
		engine := techStore.GetEngine(token.Design.Spec.Engine)
		freeSpeed = MinInt(freeSpeed, engine.FreeSpeed)
	}
	return freeSpeed
}

// Get the warp factor for when we run out of fuel.
func (fleet *Fleet) getFuelGeneration(techStore *TechStore, player *Player, warpFactor int, distance float64) int {
	// find the lowest freeSpeed from all the fleet's engines
	var freeSpeed = math.MaxInt
	for _, token := range fleet.Tokens {
		engine := techStore.GetEngine(token.Design.Spec.Engine)
		freeSpeed = MinInt(freeSpeed, engine.FreeSpeed)
	}
	return freeSpeed
}

// Complete a move from one waypoint to another
func (fleet *Fleet) completeMove(game *Game, wp0 Waypoint, wp1 Waypoint) {
	fleet.Position = wp1.Position

	// find out if we arrived at a planet, either by reaching our target fleet
	// or reaching a planet
	if wp1.TargetType == MapObjectTypeFleet && wp1.TargetPlayerNum != nil && wp1.TargetNum != nil {
		target := game.getFleet(*wp1.TargetPlayerNum, *wp1.TargetNum)
		fleet.Orbiting = target.Orbiting
	}

	if wp1.TargetType == MapObjectTypePlanet && wp1.TargetNum != nil {
		target := game.getPlanet(*wp1.TargetNum)
		fleet.Orbiting = true
		if fleet.PlayerNum == target.PlayerNum && target.Spec.HasStarbase {
			// refuel at starbases
			fleet.Fuel = fleet.Spec.FuelCapacity
		}
	} else if wp1.TargetType == MapObjectTypeWormhole {
		target := game.getWormhole(*wp1.TargetNum)
		dest := game.getWormhole(target.DestinationNum)
		fleet.Position = dest.Position
	}

	// if we wait at a waypoint while unloading, we "complete" our move but don't actually move
	// TODO: this is weird, can we just not complete a move if we are waiting at a waypoint?
	if !wp0.WaitAtWaypoint {
		fleet.Waypoints = fleet.Waypoints[1:]
	}

	// we arrived, process the current task (the previous waypoint)
	if len(fleet.Waypoints) == 1 {
		fleet.WarpSpeed = 0
		fleet.Heading = Vector{}
	} else {
		wp1 = fleet.Waypoints[1]
		fleet.WarpSpeed = wp1.WarpFactor
		fleet.Heading = (wp1.Position.Subtract(fleet.Position)).Normalized()
	}
}
