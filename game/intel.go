package game

import (
	"fmt"
	"time"
)

type MapObjectIntel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Dirty     bool      `json:"-" gorm:"-"`
	GameID    uint      `json:"gameId"`
	Position  Vector    `json:"position" gorm:"embedded"`
	Name      string    `json:"name"`
	Num       int       `json:"num"`
	PlayerNum *int      `json:"playerNum"`
	PlayerID  uint      `json:"playerId"`
	ReportAge int       `json:"reportAge"`
}

func (mo *MapObjectIntel) String() string {
	return fmt.Sprintf("GameID: %5d, ID: %5d, Num: %3d %s", mo.GameID, mo.ID, mo.Num, mo.Name)
}

type FleetIntel struct {
	MapObjectIntel
	PlanetIntelID uint `json:"-"` // for starbase fleets that are owned by a planet
}

type PlanetIntel struct {
	MapObjectIntel
	Hab                  Hab         `json:"hab,omitempty" gorm:"embedded;embeddedPrefix:hab_"`
	MineralConcentration Mineral     `json:"mineralConcentration,omitempty" gorm:"embedded;embeddedPrefix:mineral_conc_"`
	Population           uint        `json:"population,omitempty"`
	Starbase             *FleetIntel `json:"starbase,omitempty"`
}

func (p *PlanetIntel) String() string {
	return fmt.Sprintf("Planet %s", &p.MapObjectIntel)
}

func discoverPlanet(rules *Rules, player *Player, planet *Planet, penScanned bool) error {

	var intel *PlanetIntel
	planetIndex := planet.Num - 1

	if planetIndex < 0 || planetIndex >= len(player.PlanetIntels) {
		return fmt.Errorf("player %s cannot discover planet %s, planetIndex %d out of range", player, planet, planetIndex)
	}

	intel = &player.PlanetIntels[planetIndex]

	// if this intel is new, make sure it saves to the DB
	// once we create the object in the DB, it only gets saved to the DB
	// again if pen scanned
	if intel.PlayerID == 0 {
		intel.PlayerID = player.ID // this player owns this intel
		intel.Dirty = true
		intel.ReportAge = -1
	}

	// everyone knows these about planets
	intel.GameID = planet.GameID
	intel.Position = planet.Position
	intel.Name = planet.Name
	intel.Num = planet.Num

	ownedByPlayer := planet.PlayerNum != nil && player.Num == *planet.PlayerNum

	if penScanned || ownedByPlayer {
		intel.Dirty = true // flag for update
		intel.PlayerNum = planet.PlayerNum

		// if we pen scanned the planet, we learn some things
		intel.Hab = planet.Hab
		intel.MineralConcentration = planet.MineralConcentration
		intel.ReportAge = 0

		// players know their planet pops, but other planets are slightly off
		if ownedByPlayer {
			intel.Population = uint(planet.Population())
		} else {
			var randomPopulationError = rules.Random.Float64()*(rules.PopulationScannerError-(-rules.PopulationScannerError)) - rules.PopulationScannerError
			intel.Population = uint(float64(planet.Population()) * (1 - randomPopulationError))
		}
	}
	return nil
}
