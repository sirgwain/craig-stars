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

var teamsterSlots []cs.ShipDesignSlot = []cs.ShipDesignSlot{
	{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: cs.Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
	{HullComponent: cs.RhinoScanner.Name, HullSlotIndex: 3, Quantity: 1},
}
