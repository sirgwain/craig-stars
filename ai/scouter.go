package ai

import (
	"math"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"golang.org/x/exp/slices"
)

// dispatch scouts to unknown planets
func (ai *aiPlayer) scout() error {
	unknownPlanetsByNum := map[int]cs.PlanetIntel{}

	// find all the unexplored planets
	for _, planet := range ai.Player.PlanetIntels {
		if planet.Unexplored() {
			unknownPlanetsByNum[planet.Num] = planet
		}
	}

	// find all idle fleets that have scanners
	scannerFleets := []*cs.Fleet{}
	for _, fleet := range ai.Fleets {
		if _, contains := fleet.Spec.Purposes[cs.ShipDesignPurposeScout]; contains && fleet.Spec.Scanner {
			if len(fleet.Waypoints) <= 1 {
				// this fleet can be sent to scan a planet
				scannerFleets = append(scannerFleets, fleet)
			} else {
				// this fleet is already scanning a planet, remove the target from the unknown planets list
				for _, wp := range fleet.Waypoints[1:] {
					if wp.TargetNum != cs.None {
						delete(unknownPlanetsByNum, wp.TargetNum)

						target := ai.getPlanetIntel(wp.TargetNum)
						if target.ReportAge != cs.ReportAgeUnexplored {
							log.Debug().
								Int64("GameID", ai.GameID).
								Int("PlayerNum", ai.Num).
								Msgf("Scout %s no longer targeting %s, it's already explored", fleet.Name, target.Name)

							// we discovered this target in some other way, change targets
							fleet.Waypoints = fleet.Waypoints[:1]
							scannerFleets = append(scannerFleets, fleet)
							continue
						}
					}
				}
			}
		}
	}

	for _, fleet := range scannerFleets {
		closestPlanet := ai.getClosestPlanetIntel(fleet.Position, unknownPlanetsByNum)
		if closestPlanet != nil {
			warpSpeed := ai.getScoutWarpSpeed(fleet, closestPlanet.Position)

			fleet.Waypoints = append(fleet.Waypoints, cs.NewPlanetWaypoint(closestPlanet.Position, closestPlanet.Num, closestPlanet.Name, warpSpeed))
			ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
			delete(unknownPlanetsByNum, closestPlanet.Num)

			log.Debug().
				Int64("GameID", ai.GameID).
				Int("PlayerNum", ai.Num).
				Int("WarpSpeed", warpSpeed).
				Msgf("Scout %s targeting %s", fleet.Name, closestPlanet.Name)

		}
	}

	// build scout ships where necessary
	if len(unknownPlanetsByNum) > 0 {
		ai.addFleetBuildRequest(cs.FleetPurposeScout, len(unknownPlanetsByNum))
	}

	return nil
}

func (ai *aiPlayer) getScoutWarpSpeed(fleet *cs.Fleet, position cs.Vector) int {
	dist := fleet.Position.DistanceTo(position)
	if float64(fleet.Fuel)/float64(fleet.Spec.FuelCapacity) > .5 {
		return ai.getMaxWarp(dist, fleet)
	}
	// slow down when we run low on fuel
	return ai.getMinimalWarp(dist, fleet.Spec.Engine.IdealSpeed, fleet)
}

// fling packets at planets to scout
func (ai *aiPlayer) scoutPackets() error {
	if !ai.Player.Race.Spec.PacketBuiltInScanner {
		return nil
	}
	unknownPlanetsByNum := map[int]cs.PlanetIntel{}

	// find all the unexplored planets
	for _, planet := range ai.Player.PlanetIntels {
		if planet.Unexplored() {
			unknownPlanetsByNum[planet.Num] = planet
		}
	}

	// remove any planets we're already targetting
	for _, packet := range ai.MineralPackets {
		// this packet travels along a path to the target
		// build a rectangle in the shape of the scan path
		//
		angle := math.Atan2(packet.Heading.Y, packet.Heading.X)
		target := ai.getPlanetIntel(packet.TargetPlanetNum)
		dist := packet.Position.DistanceTo(target.Position)

		// the rectangle is the height of the scan range and the width of the distance
		// from the packet to the target
		// i.e
		// rect origin is upper left corner
		//          *-------------
		// packet    |           |  target
		// (x,y)    *|>>>>>>>>>>>|* (x, y)
		//           |           |
		//           -------------
		height := float64(packet.ScanRangePen * 2)
		width := dist
		rect := cs.Rect{X: packet.Position.X, Y: packet.Position.Y - height/2, Width: width, Height: height}
		for num, planet := range unknownPlanetsByNum {
			// we will catch this in our scanners, remove it
			if rect.PointInRotatedRectangle(planet.Position, angle) {
				delete(unknownPlanetsByNum, num)
			}
		}
		delete(unknownPlanetsByNum, packet.TargetPlanetNum)
	}

	for _, planet := range ai.Planets {
		if planet.Spec.HasMassDriver {
			existingQueueItemIndex := slices.IndexFunc(planet.ProductionQueue, func(item cs.ProductionQueueItem) bool { return item.Type.IsPacket() })
			if existingQueueItemIndex == -1 {

				// find the farthest planet
				farthest := ai.getFarthestPlanetIntel(planet.Position, unknownPlanetsByNum)
				if farthest == nil {
					continue
				}

				// fling a packet with the mineral we have the most of
				cargoType := planet.Cargo.GreatestMineralType()
				queueItemType := cs.QueueItemTypeMixedMineralPacket
				switch cargoType {
				case cs.Ironium:
					queueItemType = cs.QueueItemTypeIroniumMineralPacket
				case cs.Boranium:
					queueItemType = cs.QueueItemTypeBoraniumMineralPacket
				case cs.Germanium:
					queueItemType = cs.QueueItemTypeGermaniumMineralPacket
				}

				// Build a new packet targetted towards this planet
				planet.PacketTargetNum = farthest.Num
				planet.ProductionQueue = append([]cs.ProductionQueueItem{{Type: queueItemType, Quantity: 1}}, planet.ProductionQueue...)
				delete(unknownPlanetsByNum, farthest.Num)

				if err := ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders); err != nil {
					return err
				}

				log.Debug().
					Int64("GameID", ai.GameID).
					Int("PlayerNum", ai.Num).
					Int("WarpSpeed", planet.PacketSpeed).
					Msgf("Planet %s is sending a scout packet to %s", planet.Name, farthest.Name)

			}
		}
	}
	return nil
}
