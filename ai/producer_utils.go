package ai

import (
	"github.com/sirgwain/craig-stars/cs"
)

// check if an item type is in the queue
func (ai *aiPlayer) isItemInQueue(planet *cs.Planet, t cs.QueueItemType) bool {
	for _, item := range planet.ProductionQueue {
		if item.Type == t {
			return true
		}
	}
	return false
}

// check if a starbase is already in the queue
func (ai *aiPlayer) isStarbaseInQueue(planet *cs.Planet) bool {
	for _, item := range planet.ProductionQueue {
		if item.Type == cs.QueueItemTypeStarbase {
			return true
		}
	}
	return false
}

// check if a ShipDesign with the given purpose is in the queue
func (ai *aiPlayer) isShipInQueue(planet *cs.Planet, fleetPurpose cs.FleetPurpose, purpose cs.ShipDesignPurpose, quantity int) bool {
	for _, item := range planet.ProductionQueue {
		if item.Type == cs.QueueItemTypeShipToken && item.GetTag(cs.TagPurpose) == string(fleetPurpose) {
			design := ai.GetDesign(item.DesignNum)
			if design != nil && design.Purpose == purpose && item.Quantity >= quantity {
				return true
			}
		}
	}
	return false
}

// prepend item to this planet's production queue.
func (ai *aiPlayer) prependToQueue(planet *cs.Planet, item cs.ProductionQueueItem) {
	planet.ProductionQueue = append([]cs.ProductionQueueItem{item}, planet.ProductionQueue...)
}
