//go:build !wasi && !wasm

package testgames

import "github.com/sirgwain/craig-stars/cs"

var longRangeScoutSlots []cs.ShipDesignSlot = []cs.ShipDesignSlot{
	{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: cs.RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
	{HullComponent: cs.FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
}

var santaMariaSlots []cs.ShipDesignSlot = []cs.ShipDesignSlot{
	{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: cs.ColonizationModule.Name, HullSlotIndex: 2, Quantity: 1},
}

var santaMariaARSlots []cs.ShipDesignSlot = []cs.ShipDesignSlot{
	{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: cs.OrbitalConstructionModule.Name, HullSlotIndex: 2, Quantity: 1},
}

var teamsterSlots []cs.ShipDesignSlot = []cs.ShipDesignSlot{
	{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: cs.Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
	{HullComponent: cs.RhinoScanner.Name, HullSlotIndex: 3, Quantity: 1},
}

var destroyerDeltaSlots []cs.ShipDesignSlot = []cs.ShipDesignSlot{
	{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: cs.DeltaTorpedo.Name, HullSlotIndex: 2, Quantity: 1},
	{HullComponent: cs.DeltaTorpedo.Name, HullSlotIndex: 3, Quantity: 1},
	{HullComponent: cs.DeltaTorpedo.Name, HullSlotIndex: 4, Quantity: 1},
	{HullComponent: cs.Crobmnium.Name, HullSlotIndex: 5, Quantity: 1},
	{HullComponent: cs.ManeuveringJet.Name, HullSlotIndex: 6, Quantity: 1},
	{HullComponent: cs.BattleComputer.Name, HullSlotIndex: 7, Quantity: 1},
}
