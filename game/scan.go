package game

import "math"

type scanner struct {
	Position            Vector
	RangeSquared        int
	RangePenSquared     int
	DiscoverFleetCargo  bool
	DiscoverPlanetCargo bool
	CloakReduction      float64
}

type playerScan struct {
	universe   *Universe
	rules      *Rules
	player     *Player
	discoverer discoverer
}

type playerScanner interface {
	scan() error
}

func newPlayerScanner(universe *Universe, rules *Rules, player *Player) playerScanner {
	discoverer := newDiscoverer(player)
	return &playerScan{universe, rules, player, discoverer}
}

// scan planets, fleets, etc for a player
func (scan *playerScan) scan() error {
	// clear out any reports that we recreate each year
	player := scan.player
	scan.discoverer.clearTransientReports()

	for i := range player.PlanetIntels {
		planet := &player.PlanetIntels[i]
		if planet.ReportAge != Unexplored {
			planet.ReportAge++
		}
	}

	// TODO: add in player mineral packets, minefields, etc
	scanners := scan.getScanners()
	cargoScanners := scan.getCargoScanners()

	// scan planets
	if err := scan.scanPlanets(scanners, cargoScanners); err != nil {
		return err
	}

	// scan fleets
	scan.scanFleets(scanners, cargoScanners)

	return nil
}

// scan all planets with this player's scanners
func (scan *playerScan) scanPlanets(scanners []scanner, cargoScanners []scanner) error {
	for _, planet := range scan.universe.Planets {
		if planet.ownedBy(scan.player.Num) {
			if err := scan.discoverer.discoverPlanet(scan.rules, scan.player, planet, false); err != nil {
				return err
			}
			continue
		}

		// try and scan the planet with this scanner
		for _, scanner := range scanners {
			scanned, err := scan.scanPlanet(planet, scanner)
			if err != nil {
				return err
			}
			if scanned {
				break
			}
		}

		// try and scan the planet with a cargo scanner
		for _, scanner := range cargoScanners {
			scanned, err := scan.scanPlanet(planet, scanner)
			if err != nil {
				return err
			}
			if scanned {
				break
			}
		}
	}

	return nil
}

// scan this planet
func (scan *playerScan) scanPlanet(planet *Planet, scanner scanner) (bool, error) {
	dist := scanner.Position.DistanceSquaredTo(planet.Position)
	_ = dist
	if float64(scanner.RangePenSquared) >= scanner.Position.DistanceSquaredTo(planet.Position) {
		if err := scan.discoverer.discoverPlanet(scan.rules, scan.player, planet, true); err != nil {
			return false, err
		}
		if scanner.DiscoverPlanetCargo {
			if err := scan.discoverer.discoverPlanetCargo(scan.player, planet); err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return false, nil
}

// scan all fleets and discover their designs if we should
func (scan *playerScan) scanFleets(scanners []scanner, cargoScanners []scanner) {
	// scan fleets
	fleetsToScan := []*Fleet{}
	fleetsToCargoScan := []*Fleet{}
	for _, fleet := range scan.universe.Fleets {
		if fleet.ownedBy(scan.player.Num) {
			// The player already gets a copy of all their own fleets
			continue
		}

		// try and scan the planet with this scanner
		for _, scanner := range scanners {
			if scan.fleetInScannerRange(fleet, scanner) {
				fleetsToScan = append(fleetsToScan, fleet)
				break
			}
		}

		// try and scan the planet with a cargo scanner
		for _, scanner := range cargoScanners {
			if scan.fleetInScannerRange(fleet, scanner) {
				fleetsToCargoScan = append(fleetsToScan, fleet)
				break
			}
		}
	}

	for _, fleet := range fleetsToScan {
		scan.discoverer.discoverFleet(scan.player, fleet)
		if scan.player.Race.Spec.DiscoverDesignOnScan {
			for _, token := range fleet.Tokens {
				scan.discoverer.discoverDesign(scan.player, token.design, true)
			}
		}
	}

	for _, fleet := range fleetsToCargoScan {
		scan.discoverer.discoverFleetCargo(scan.player, fleet)
	}
}

// return true if this scanner successfully scans this fleet, taking into account cloaking
// and the fleet's cloak penetration
func (scan *playerScan) fleetInScannerRange(fleet *Fleet, scanner scanner) bool {
	var cloakFactor = 1 - (float64(fleet.Spec.CloakPercent) * (1 - scanner.CloakReduction) / 100.0)
	var distance = scanner.Position.DistanceSquaredTo(fleet.Position)

	// if we pen scanned this, update the report
	if float64(scanner.RangePenSquared)*cloakFactor >= distance {
		// update the fleet report with pen scanners
		return true
	}

	// if we aren't orbiting a planet, we can be seen with regular scanners
	if !fleet.Orbiting() && float64(scanner.RangeSquared)*cloakFactor >= distance {
		return true
	}
	return false
}

// get a list of unique scanners per player.
// This is a minimal list only containing the best scanner values for each position
func (scan *playerScan) getScanners() []scanner {
	planetaryScanner := scan.player.Spec.PlanetaryScanner
	scanningFleetsByPosition := map[Vector]scanner{}
	for _, fleet := range scan.universe.Fleets {
		if fleet.PlayerNum == scan.player.Num && fleet.Spec.Scanner {
			scanner, found := scanningFleetsByPosition[fleet.Position]
			if !found {
				// start with NoScanner (-1)
				scanner.Position = fleet.Position
				scanner.RangeSquared = NoScanner
				scanner.RangePenSquared = NoScanner
			}
			scanner.RangeSquared = maxInt(scanner.RangeSquared, fleet.Spec.ScanRange*fleet.Spec.ScanRange)
			scanner.RangePenSquared = maxInt(scanner.RangePenSquared, fleet.Spec.ScanRangePen*fleet.Spec.ScanRangePen)
			scanner.CloakReduction = math.Max(scanner.CloakReduction, fleet.Spec.ReduceCloaking)
			scanningFleetsByPosition[fleet.Position] = scanner
		}
	}

	// build a list of scanners for this player
	scanners := []scanner{}
	for _, planet := range scan.universe.Planets {
		if planet.PlayerNum == scan.player.Num && planet.Scanner {
			scanner := scanner{
				Position:        planet.Position,
				RangeSquared:    planetaryScanner.ScanRange * planetaryScanner.ScanRange,
				RangePenSquared: planetaryScanner.ScanRangePen * planetaryScanner.ScanRangePen,
			}

			// use the fleet scanner if it's better
			if fleetScanner, ok := scanningFleetsByPosition[planet.Position]; ok {
				scanner.RangeSquared = maxInt(scanner.RangeSquared, fleetScanner.RangeSquared)
				scanner.RangePenSquared = maxInt(scanner.RangePenSquared, fleetScanner.RangePenSquared)
				scanner.CloakReduction = math.Max(scanner.CloakReduction, fleetScanner.CloakReduction)
			}
			scanners = append(scanners, scanner)
		}
	}

	// add in any fleet scanners that weren't on a planet
	if len(scanners) == 0 {
		// we have no planetary scanners (weird, but possible if all planets with scanners are destroyed)
		// so just add the fleet scanners to the list
		for _, fleetScanner := range scanningFleetsByPosition {
			scanners = append(scanners, fleetScanner)
		}
	} else {
		scannersByPosition := map[Vector]scanner{}
		for _, scanner := range scanners {
			scannersByPosition[scanner.Position] = scanner
		}
		for position, fleetScanner := range scanningFleetsByPosition {
			// if we don't find a scanner at this position, add the fleetScanner
			// to our existing scanners list
			if _, found := scannersByPosition[position]; !found {
				scanners = append(scanners, fleetScanner)
			}
		}
	}

	return scanners
}

// get a list of scanners that can scan cargo from fleets or planets
func (scan *playerScan) getCargoScanners() []scanner {
	scanners := []scanner{}
	scanningFleetsByPosition := map[Vector]scanner{}

	for i := range scan.universe.Fleets {
		fleet := scan.universe.Fleets[i]
		if fleet.PlayerNum == scan.player.Num && fleet.Spec.Scanner && (fleet.Spec.CanStealFleetCargo || fleet.Spec.CanStealPlanetCargo) {
			scanner, found := scanningFleetsByPosition[fleet.Position]
			if !found {
				// start with NoScanner (-1)
				scanner.Position = fleet.Position
				scanner.RangeSquared = NoScanner
				scanner.RangePenSquared = NoScanner
			}
			scanner.RangeSquared = maxInt(scanner.RangeSquared, fleet.Spec.ScanRange*fleet.Spec.ScanRange)
			scanner.RangePenSquared = maxInt(scanner.RangePenSquared, fleet.Spec.ScanRangePen*fleet.Spec.ScanRangePen)
			scanner.CloakReduction = math.Max(scanner.CloakReduction, fleet.Spec.ReduceCloaking)
			scanner.DiscoverFleetCargo = fleet.Spec.CanStealFleetCargo
			scanner.DiscoverPlanetCargo = fleet.Spec.CanStealPlanetCargo
			scanningFleetsByPosition[fleet.Position] = scanner
		}
	}

	for _, fleetScanner := range scanningFleetsByPosition {
		scanners = append(scanners, fleetScanner)
	}

	return scanners
}
