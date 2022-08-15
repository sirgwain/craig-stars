package game

import (
	"fmt"
	"math"
	"time"
)

type Fleet struct {
	MapObject
	PlanetID     uint        `json:"-"` // for starbase fleets that are owned by a planet
	BaseName     string      `json:"baseName"`
	Cargo        Cargo       `json:"cargo,omitempty" gorm:"embedded;embeddedPrefix:cargo_"`
	Fuel         int         `json:"fuel"`
	Damage       int         `json:"damage"`
	BattlePlanID uint        `json:"battlePlan"`
	Tokens       []ShipToken `json:"tokens" gorm:"constraint:OnDelete:CASCADE;"`
	Waypoints    []Waypoint  `json:"waypoints" gorm:"serializer:json"`
	Spec         *FleetSpec  `json:"spec" gorm:"serializer:json"`
}

type FleetSpec struct {
	ShipDesignSpec
	Purposes         []ShipDesignPurpose `json:"purposes"`
	TotalShips       int                 `json:"totalShips"`
	MassEmpty        int                 `json:"massEmpty"`
	BasePacketSpeed  int                 `json:"basePacketSpeed"`
	SafePacketSpeed  int                 `json:"safePacketSpeed"`
	BaseCloakedCargo int                 `json:"baseCloakedCargo"`
	HasMassDriver    bool                `json:"hasMassDriver,omitempty"`
	HasStargate      bool                `json:"hasStargate,omitempty"`
	Stargate         string              `json:"stargate,omitempty"`
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
	FleetID          uint   `json:"-"`
	TargetID         uint   `json:"targetId,omitempty"`
	Position         Vector `json:"position,omitempty" gorm:"embedded"`
	WarpFactor       int    `json:"warpFactor,omitempty"`
	TargetPlanetNum  *int   `json:"targetPlanetNum,omitempty"`
	TransferToPlayer *int   `json:"transferToPlayer,omitempty"`
	TargetName       string `json:"targetName,omitempty"`
}

func NewFleet(player *Player, design *ShipDesign, num int, name string, waypoints []Waypoint) Fleet {
	return Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			GameID:    player.GameID,
			PlayerID:  player.ID,
			PlayerNum: &player.Num,
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

func (f *Fleet) WithCargo(cargo Cargo) *Fleet {
	f.Cargo = cargo
	return f
}

func NewPlanetWaypoint(planet *Planet, warpFactor int) Waypoint {
	return Waypoint{
		Position:        planet.Position,
		TargetPlanetNum: &planet.Num,
		TargetName:      planet.Name,
		WarpFactor:      warpFactor,
	}
}

func ComputeFleetSpec(rules *Rules, player *Player, fleet *Fleet) *FleetSpec {
	spec := FleetSpec{
		ShipDesignSpec: ShipDesignSpec{
			ScanRange:    NoScanner,
			ScanRangePen: NoScanner,
			SpaceDock:    UnlimitedSpaceDock,
		},
	}
	spec.Mass = fleet.Cargo.Total()

	for _, token := range fleet.Tokens {

		// update our total ship count
		spec.TotalShips += token.Quantity

		if token.Design.Purpose != ShipDesignPurposeNone {
			spec.Purposes = append(spec.Purposes, token.Design.Purpose)
		}

		// TODO: which default engine do we use for multiple fleets?
		spec.Engine = token.Design.Spec.Engine
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

	return &spec
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
