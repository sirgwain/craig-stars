package ai

import (
	"cmp"
	"slices"
	"strconv"

	"github.com/sirgwain/craig-stars/cs"
)

// the tag recording the year a fleet was first found with nothing to do
const tagIdleSince = "idleSince"

const (
	retireIdleYears       = 5   // years a ship sits idle before it is considered surplus
	partialFleetShare     = .6  // share of the war fleet recipe that is worth sending on its own
	obsoletePowerFactor   = 1.5 // how much stronger a newer design must be to retire an old one
	transportFleetsToKeep = 1   // idle transport fleets kept for each world with colonists to spare
)

// The AI is rebuilt every turn, so remember on the fleet how long it has been idle.
func (ai *aiPlayer) trackIdleFleets() {
	for _, fleet := range ai.Fleets {
		if !fleet.Idle() {
			delete(fleet.Tags, tagIdleSince)
		} else if fleet.GetTag(tagIdleSince) == "" {
			fleet.SetTag(tagIdleSince, strconv.Itoa(ai.game.Year))
		}
	}
}

func (ai *aiPlayer) idleYears(fleet *cs.Fleet) int {
	since, err := strconv.Atoi(fleet.GetTag(tagIdleSince))
	if err != nil || !fleet.Idle() {
		return 0
	}
	return ai.game.Year - since
}

// Scrap a fleet at the nearest starbase, which recovers most of its minerals.
func (ai *aiPlayer) scrap(fleet *cs.Fleet) {
	if planet := ai.getPlanet(fleet.OrbitingPlanetNum); planet != nil && planet.Spec.HasStarbase {
		fleet.Waypoints[0].Task = cs.WaypointTaskScrapFleet
		ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
		return
	}
	starbase := ai.getClosestStarbasePlanet(fleet)
	if starbase == nil {
		return
	}
	if index := fleet.AddWaypoint(ai.Player, cs.WaypointDest{MO: starbase.MapObject}, 0, false); index > 0 {
		fleet.Waypoints[index].Task = cs.WaypointTaskScrapFleet
		ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
	}
}

// Keep a transport fleet for each world that can supply colonists and scrap
// the rest once they have had nothing to carry for years.
func (ai *aiPlayer) retireTransports() {
	feeders := 0
	for _, planet := range ai.Planets {
		if ai.exportableColonists(planet, ai.config.colonistTransportDensity) >= minColonistLoad {
			feeders++
		}
	}
	idle := []*cs.Fleet{}
	for _, fleet := range ai.Fleets {
		if fleet.GetTag(cs.TagPurpose) == string(cs.FleetPurposeColonistFreighter) && !fleet.Spec.Colonizer && fleet.Idle() && fleet.Cargo.Total() == 0 {
			idle = append(idle, fleet)
		}
	}
	// the longest idle go first
	slices.SortStableFunc(idle, func(a, b *cs.Fleet) int {
		return cmp.Or(cmp.Compare(ai.idleYears(b), ai.idleYears(a)), cmp.Compare(a.Num, b.Num))
	})
	surplus := len(idle) - max(1, transportFleetsToKeep*feeders)
	for _, fleet := range idle {
		if surplus <= 0 || ai.idleYears(fleet) < retireIdleYears {
			break
		}
		ai.scrap(fleet)
		surplus--
	}
}

// Warships built for a fleet that never filled up are no use sitting at the
// shipyard. Send a mostly complete fleet as it is, and scrap ships whose
// design has been clearly surpassed.
func (ai *aiPlayer) retireWarships() error {
	recipe := 0
	for _, quantity := range ai.fleetsByPurpose[cs.FleetPurposeBomber].requiredShips() {
		recipe += quantity
	}

	planetNums := make([]int, 0, len(ai.fleetsByPlanetNum))
	for planetNum := range ai.fleetsByPlanetNum {
		if planetNum != cs.None {
			planetNums = append(planetNums, planetNum)
		}
	}
	slices.Sort(planetNums)
	for _, planetNum := range planetNums {
		waiting := []*cs.Fleet{}
		ships, bombers := 0, false
		for _, fleet := range ai.fleetsByPlanetNum[planetNum] {
			if fleet.GetTag(cs.TagPurpose) != string(cs.FleetPurposeBomber) || len(fleet.Tokens) != 1 || ai.idleYears(fleet) < retireIdleYears {
				continue
			}
			design := ai.GetDesign(fleet.Tokens[0].DesignNum)
			if design == nil {
				continue
			}
			if latest := ai.designsByPurpose[design.Purpose]; latest != nil && latest.Num != design.Num && design.Spec.PowerRating > 0 &&
				float64(latest.Spec.PowerRating) >= obsoletePowerFactor*float64(design.Spec.PowerRating) {
				ai.scrap(fleet)
				continue
			}
			waiting = append(waiting, fleet)
			ships += fleet.Spec.TotalShips
			bombers = bombers || fleet.Spec.Bomber
		}
		if len(waiting) > 1 && bombers && recipe > 0 && float64(ships) >= partialFleetShare*float64(recipe) {
			fleet, err := ai.merge(waiting)
			if err != nil {
				return err
			}
			fleet.Rename(ai.fleetName(cs.FleetPurposeBomber))
			fleet.SetTag(cs.TagPurpose, string(cs.FleetPurposeBomber))
			fleet.Purpose = cs.FleetPurposeBomber
		}
	}
	return nil
}

// Scouts with nothing left to explore watch worlds our planetary scanners
// can't see. A scout with nowhere useful to watch, or sitting inside
// the range of a better penetrating scanner, is scrapped.
func (ai *aiPlayer) picketScouts() {
	// worlds worth watching: explored, not ours, and outside our planetary scanners
	posts := map[int]*cs.Planet{}
	for _, intel := range ai.PlanetIntels {
		if intel.Explored() && !intel.OwnedBy(ai.Num) && !ai.inPlanetaryScanRange(intel.Position) {
			posts[intel.Num] = intel
		}
	}
	scouts := []*cs.Fleet{}
	for _, fleet := range ai.Fleets {
		if fleet.GetTag(cs.TagPurpose) != string(cs.FleetPurposeScout) || !fleet.Spec.Scanner {
			continue
		}
		if !fleet.Idle() {
			// on its way to a post
			if len(fleet.Waypoints) > 1 {
				delete(posts, fleet.Waypoints[len(fleet.Waypoints)-1].TargetNum)
			}
			continue
		}
		if _, onPost := posts[fleet.OrbitingPlanetNum]; onPost && !ai.outscanned(fleet) {
			delete(posts, fleet.OrbitingPlanetNum)
			continue
		}
		scouts = append(scouts, fleet)
	}

	for _, fleet := range scouts {
		if ai.outscanned(fleet) {
			ai.scrap(fleet)
			continue
		}
		// prefer watching an enemy, then whatever is closest
		var post *cs.Planet
		for _, candidate := range posts {
			if post == nil || ai.picketRank(fleet, candidate) < ai.picketRank(fleet, post) ||
				(ai.picketRank(fleet, candidate) == ai.picketRank(fleet, post) && candidate.Num < post.Num) {
				post = candidate
			}
		}
		if post == nil {
			if ai.idleYears(fleet) >= retireIdleYears {
				ai.scrap(fleet)
			}
			continue
		}
		if index := fleet.AddWaypoint(ai.Player, cs.WaypointDest{MO: post.MapObject}, 0, false); index > 0 {
			ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
			delete(posts, post.Num)
		}
	}
}

// lower is a better post for this scout
func (ai *aiPlayer) picketRank(fleet *cs.Fleet, post *cs.Planet) int {
	rank := fleet.Position.DistanceSquaredTo(post.Position)
	if !ai.IsEnemy(post.PlayerNum) {
		// any enemy world ranks ahead of every other world
		rank += 1 << 40
	}
	return rank
}

func (ai *aiPlayer) inPlanetaryScanRange(position cs.Vector) bool {
	for _, planet := range ai.Planets {
		if planet.Spec.ScanRange > 0 && position.DistanceSquaredTo(planet.Position) <= planet.Spec.ScanRange*planet.Spec.ScanRange {
			return true
		}
	}
	return false
}

// outscanned returns true if this scout sits inside the penetrating scan range
// of another of our stationary fleets with a better penetrating scanner
func (ai *aiPlayer) outscanned(scout *cs.Fleet) bool {
	for _, fleet := range ai.Fleets {
		if fleet == scout || !fleet.Idle() || fleet.Spec.ScanRangePen <= 0 || fleet.Spec.ScanRangePen <= scout.Spec.ScanRangePen {
			continue
		}
		if scout.Position.DistanceSquaredTo(fleet.Position) <= fleet.Spec.ScanRangePen*fleet.Spec.ScanRangePen {
			return true
		}
	}
	return false
}

// years a finished warship waits at its shipyard before going to the muster point
const musterIdleYears = 1

// The ships of a war fleet are often finished at different shipyards, years
// apart, and none of those piles is ever a whole fleet. Bring them together at
// the shipyard that already has the most, and form a fleet there once enough
// have gathered or nothing more is on its way.
func (ai *aiPlayer) musterWarships() error {
	if !ai.enabled("muster") {
		return nil
	}
	loose := map[int][]*cs.Fleet{}
	ships := map[int]int{}
	inbound := false
	for _, fleet := range ai.Fleets {
		if fleet.GetTag(cs.TagPurpose) != string(cs.FleetPurposeBomber) || len(fleet.Tokens) != 1 {
			continue
		}
		if !fleet.Idle() {
			inbound = true
			continue
		}
		if ai.getPlanet(fleet.OrbitingPlanetNum) == nil || ai.idleYears(fleet) < musterIdleYears {
			continue
		}
		loose[fleet.OrbitingPlanetNum] = append(loose[fleet.OrbitingPlanetNum], fleet)
		ships[fleet.OrbitingPlanetNum] += fleet.Spec.TotalShips
	}
	if len(loose) == 0 {
		return nil
	}

	// gather where the most ships already are
	muster := cs.None
	for num := range loose {
		if muster == cs.None || ships[num] > ships[muster] || (ships[num] == ships[muster] && num < muster) {
			muster = num
		}
	}
	planet := ai.getPlanet(muster)
	for num, fleets := range loose {
		if num == muster {
			continue
		}
		for _, fleet := range fleets {
			if index := fleet.AddWaypoint(ai.Player, cs.WaypointDest{MO: planet.MapObject}, 0, false); index > 0 {
				ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
				inbound = true
			}
		}
	}

	recipe := 0
	for _, quantity := range ai.fleetsByPurpose[cs.FleetPurposeBomber].requiredShips() {
		recipe += quantity
	}
	gathered := loose[muster]
	bombers := slices.ContainsFunc(gathered, func(fleet *cs.Fleet) bool { return fleet.Spec.Bomber })
	if len(gathered) < 2 || !bombers || (inbound && float64(ships[muster]) < partialFleetShare*float64(recipe)) {
		return nil
	}
	fleet, err := ai.merge(gathered)
	if err != nil {
		return err
	}
	fleet.Rename(ai.fleetName(cs.FleetPurposeBomber))
	fleet.SetTag(cs.TagPurpose, string(cs.FleetPurposeBomber))
	fleet.Purpose = cs.FleetPurposeBomber
	return nil
}
